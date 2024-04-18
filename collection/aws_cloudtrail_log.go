package collection

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	sdkconfig "github.com/turbot/tailpipe-plugin-sdk/config"
	"github.com/turbot/tailpipe-plugin-sdk/source"

	"github.com/turbot/tailpipe-plugin-aws/util"
)

type AwsCloudTrailLogCollectionConfig struct{}

type AwsCloudTrailLogCollection struct {
	Config AwsCloudTrailLogCollectionConfig

	ctx context.Context

	// observers is a list of observers that will be notified of events.
	observers      []collection.CollectionObserver
	observersMutex sync.RWMutex
}

func (c *AwsCloudTrailLogCollection) Identifier() string {
	return "aws_cloudtrail_log"
}

func (c *AwsCloudTrailLogCollection) Init(ctx context.Context) error {
	c.ctx = ctx
	return nil
}

func (c *AwsCloudTrailLogCollection) Context() context.Context {
	return c.ctx
}

func (c *AwsCloudTrailLogCollection) AddObserver(observer collection.CollectionObserver) {
	c.observersMutex.Lock()
	defer c.observersMutex.Unlock()
	c.observers = append(c.observers, observer)
}

func (c *AwsCloudTrailLogCollection) RemoveObserver(observer collection.CollectionObserver) {
	c.observersMutex.Lock()
	defer c.observersMutex.Unlock()
	for i, o := range c.observers {
		if o == observer {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			break
		}
	}
}

func (c *AwsCloudTrailLogCollection) LoadConfig(configRaw []byte) error {
	return sdkconfig.Load(configRaw, &c.Config)
}

func (c *AwsCloudTrailLogCollection) ValidateConfig() error {
	return nil
}

func (c *AwsCloudTrailLogCollection) Schema() collection.Row {
	return &AWSCloudTrail{}
}

func (c *AwsCloudTrailLogCollection) ExtractArtifactRows(ctx context.Context, a *source.Artifact) error {

	inputPath := a.Name

	gzFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	var log AWSCloudTrailBatch
	if err := json.NewDecoder(gzReader).Decode(&log); err != nil {
		return err
	}

	for _, record := range log.Records {

		// Record standardization
		record.TpID = xid.New().String()
		record.TpSourceType = "aws_cloudtrail_log"
		record.TpTimestamp = record.EventTime
		record.TpSourceLocation = &inputPath
		record.TpIngestTimestamp = UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))
		if record.SourceIPAddress != nil {
			record.TpSourceIP = record.SourceIPAddress
			record.TpIps = append(record.TpIps, *record.SourceIPAddress)
		}
		for _, resource := range record.Resources {
			if resource.ARN != nil {
				newAkas := util.AwsAkasFromArn(*resource.ARN)
				record.TpAkas = append(record.TpAkas, newAkas...)
			}
		}
		// If it's an AKIA, then record that as an identity. Do not record ASIA*
		// keys etc.
		if record.UserIdentity.AccessKeyId != nil {
			if strings.HasPrefix(*record.UserIdentity.AccessKeyId, "AKIA") {
				record.TpUsernames = append(record.TpUsernames, *record.UserIdentity.AccessKeyId)
			}
		}
		if record.UserIdentity.UserName != nil {
			record.TpUsernames = append(record.TpUsernames, *record.UserIdentity.UserName)
		}

		// Hive fields
		record.TpCollection = "default" // TODO - should be based on the definition in HCL
		record.TpConnection = record.RecipientAccountId
		record.TpYear = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Year())
		record.TpMonth = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Month())
		record.TpDay = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Day())

		//c.Collection.NotifyRow(a, &record)
		for _, o := range c.observers {
			o.NotifyRow(a, &record)
		}

	}

	return nil

}

type AWSCloudTrailBatch struct {
	Records []AWSCloudTrail `json:"Records"`
}

type AWSCloudTrail struct {

	// Metadata
	TpID              string     `json:"tp_id" parquet:"name=tp_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	TpSourceType      string     `json:"tp_source_type" parquet:"name=tp_source_type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpSourceName      string     `json:"tp_source_name" parquet:"name=tp_source_name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpSourceLocation  *string    `json:"tp_source_location" parquet:"name=tp_source_location, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpIngestTimestamp UnixMillis `json:"tp_ingest_timestamp" parquet:"name=tp_ingest_timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"`

	// Standardized
	TpTimestamp     UnixMillis `json:"tp_timestamp" parquet:"name=tp_timestamp, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	TpSourceIP      *string    `json:"tp_source_ip" parquet:"name=tp_source_ip, type=BYTE_ARRAY, convertedtype=UTF8"`
	TpDestinationIP *string    `json:"tp_destination_ip" parquet:"name=tp_destination_ip, type=BYTE_ARRAY, convertedtype=UTF8"`

	// Hive fields
	TpCollection string `json:"tp_collection" parquet:"name=tp_collection, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpConnection string `json:"tp_connection" parquet:"name=tp_connection, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpYear       int32  `json:"tp_year" parquet:"name=tp_year, type=INT32, convertedtype=INT32"`
	TpMonth      int32  `json:"tp_month" parquet:"name=tp_month, type=INT32, convertedtype=INT32"`
	TpDay        int32  `json:"tp_day" parquet:"name=tp_day, type=INT32, convertedtype=INT32"`

	// Searchable
	TpAkas      []string `json:"tp_akas,omitempty" parquet:"name=tp_akas, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpIps       []string `json:"tp_ips,omitempty" parquet:"name=tp_ips, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8"`
	TpTags      []string `json:"tp_tags,omitempty" parquet:"name=tp_tags, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8"`
	TpDomains   []string `json:"tp_domains,omitempty" parquet:"name=tp_domains, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpEmails    []string `json:"tp_emails,omitempty" parquet:"name=tp_emails, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TpUsernames []string `json:"tp_usernames,omitempty" parquet:"name=tp_usernames, type=LIST, convertedtype=LIST, valuetype=BYTE_ARRAY, valueconvertedtype=UTF8, encoding=PLAIN_DICTIONARY"`

	AdditionalEventData *JSONString  `json:"additionalEventData,omitempty" parquet:"name=additional_event_data, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	APIVersion          *string      `json:"apiVersion,omitempty" parquet:"name=api_version, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	AwsRegion           string       `json:"awsRegion" parquet:"name=aws_region, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	ErrorCode           *string      `json:"errorCode,omitempty" parquet:"name=error_code, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	ErrorMessage        *string      `json:"errorMessage,omitempty" parquet:"name=error_message, type=BYTE_ARRAY, convertedtype=UTF8"`
	EventID             string       `json:"eventID" parquet:"name=event_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	EventName           string       `json:"eventName" parquet:"name=event_name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	EventSource         string       `json:"eventSource" parquet:"name=event_source, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	EventTime           UnixMillis   `json:"eventTime" parquet:"name=event_time, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	EventType           string       `json:"eventType" parquet:"name=event_type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	EventVersion        string       `json:"eventVersion" parquet:"name=event_version, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	ManagementEvent     *bool        `json:"managementEvent,omitempty" parquet:"name=management_event, type=BOOLEAN"`
	ReadOnly            *bool        `json:"readOnly,omitempty" parquet:"name=read_only, type=BOOLEAN"`
	RecipientAccountId  string       `json:"recipientAccountId,omitempty" parquet:"name=recipient_account_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	RequestID           *string      `json:"requestID,omitempty" parquet:"name=request_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	RequestParameters   *JSONString  `json:"requestParameters,omitempty" parquet:"name=request_parameters, type=BYTE_ARRAY, convertedtype=UTF8"`
	Resources           []*Resource  `json:"resources,omitempty" parquet:"name=resources, type=LIST"`
	ResponseElements    *JSONString  `json:"responseElements,omitempty" parquet:"name=response_elements, type=BYTE_ARRAY, convertedtype=UTF8"`
	ServiceEventDetails *JSONString  `json:"serviceEventDetails,omitempty" parquet:"name=service_event_details, type=BYTE_ARRAY, convertedtype=UTF8"`
	SharedEventID       *string      `json:"sharedEventID,omitempty" parquet:"name=shared_event_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	SourceIPAddress     *string      `json:"sourceIPAddress,omitempty" parquet:"name=source_ip_address, type=BYTE_ARRAY, convertedtype=UTF8"`
	UserAgent           *string      `json:"userAgent,omitempty" parquet:"name=user_agent, type=BYTE_ARRAY, convertedtype=UTF8"`
	UserIdentity        UserIdentity `json:"userIdentity" parquet:"name=user_identity, type=STRUCT"`
	VpcEndpointId       string       `json:"vpcEndpointId,omitempty" parquet:"name=vpc_endpoint_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	EventCategory       string       `json:"eventCategory,omitempty" parquet:"name=event_category, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	// TODO - this looks like a bool, but is in the JSON as a string ... should we convert it?
	SessionCredentialFromConsole *string     `json:"sessionCredentialFromConsole,omitempty" parquet:"name=session_credential_from_console, type=BYTE_ARRAY, convertedtype=UTF8"`
	EdgeDeviceDetails            *JSONString `json:"edgeDeviceDetails,omitempty" parquet:"name=edge_device_details, type=BYTE_ARRAY, convertedtype=UTF8"`
	TLSDetails                   *TLSDetails `json:"tlsDetails,omitempty" parquet:"name=tls_details, type=STRUCT"`
}

func (a *AWSCloudTrail) GetTpID() string {
	return a.TpID
}

func (a *AWSCloudTrail) GetTpTimestamp() int64 {
	return int64(a.TpTimestamp)
}

func (a *AWSCloudTrail) GetConnection() string {
	return a.TpConnection
}

func (a *AWSCloudTrail) GetYear() int {
	return int(a.TpYear)
}

func (a *AWSCloudTrail) GetMonth() int {
	return int(a.TpMonth)
}

func (a *AWSCloudTrail) GetDay() int {
	return int(a.TpDay)
}

type UserIdentity struct {
	Type             string          `json:"type" parquet:"name=type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	PrincipalId      *string         `json:"principalId,omitempty" parquet:"name=principal_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	ARN              *string         `json:"arn,omitempty" parquet:"name=arn, type=BYTE_ARRAY, convertedtype=UTF8"`
	AccountId        *string         `json:"accountId,omitempty" parquet:"name=account_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	AccessKeyId      *string         `json:"accessKeyId,omitempty" parquet:"name=access_key_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	UserName         *string         `json:"userName,omitempty" parquet:"name=user_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	SessionContext   *SessionContext `json:"sessionContext,omitempty" parquet:"name=session_context, type=STRUCT"`
	InvokedBy        *string         `json:"invokedBy,omitempty" parquet:"name=invoked_by, type=BYTE_ARRAY, convertedtype=UTF8"`
	IdentityProvider *string         `json:"identityProvider,omitempty" parquet:"name=identity_provider, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type Resource struct {
	ARN       *string `json:"ARN,omitempty" parquet:"name=arn, type=BYTE_ARRAY, convertedtype=UTF8"`
	AccountId *string `json:"accountId,omitempty" parquet:"name=account_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Type      *string `json:"type,omitempty" parquet:"name=type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

type SessionContext struct {
	Attributes          *SessionAttributes   `json:"attributes,omitempty" parquet:"name=attributes, type=STRUCT"`
	SessionIssuer       *SessionIssuer       `json:"sessionIssuer,omitempty" parquet:"name=session_issuer, type=STRUCT"`
	WebIdFederationData *WebIdFederationData `json:"webIdFederationData,omitempty" parquet:"name=web_id_federation_data, type=STRUCT"`
	EC2RoleDelivery     *string              `json:"ec2RoleDelivery,omitempty" parquet:"name=ec2_role_delivery, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type SessionAttributes struct {
	MfaAuthenticated *string     `json:"mfaAuthenticated,omitempty" parquet:"name=mfa_authenticated, type=BYTE_ARRAY, convertedtype=UTF8"`
	CreationDate     *UnixMillis `json:"creationDate,omitempty" parquet:"name=creation_date, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
}

type SessionIssuer struct {
	Type        *string `json:"type,omitempty" parquet:"name=type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	PrincipalId *string `json:"principalId,omitempty" parquet:"name=principal_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	ARN         *string `json:"arn,omitempty" parquet:"name=arn, type=BYTE_ARRAY, convertedtype=UTF8"`
	AccountId   *string `json:"accountId,omitempty" parquet:"name=account_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	UserName    *string `json:"userName,omitempty" parquet:"name=user_name, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type WebIdFederationData struct {
	FederatedProvider *string     `json:"federatedProvider,omitempty" parquet:"name=federated_provider, type=BYTE_ARRAY, convertedtype=UTF8"`
	Attributes        *JSONString `json:"attributes,omitempty" parquet:"name=attributes, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type TLSDetails struct {
	TLSVersion  *string `json:"tlsVersion,omitempty" parquet:"name=tls_version, type=BYTE_ARRAY, convertedtype=UTF8"`
	CipherSuite *string `json:"cipherSuite,omitempty" parquet:"name=cipher_suite, type=BYTE_ARRAY, convertedtype=UTF8"`
	//ClientProvidedHostHeader *string `json:"clientProvidedHostHeader,omitempty" parquet:"name=client_provided_host_header, type=BYTE_ARRAY, convertedtype=UTF8"`
	ClientProvidedHostHeader *string `json:"clientProvidedHostHeader,omitempty" parquet:"name=client_provided_host_header, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
}
