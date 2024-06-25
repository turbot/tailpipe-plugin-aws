package aws_types

import (
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
)

type AWSCloudTrailBatch struct {
	Records []AWSCloudTrail `json:"Records"`
}

type AWSCloudTrail struct {
	// embed required enrichment fields
	enrichment.EnrichmentFields

	AdditionalEventData          *helpers.JSONString `json:"additionalEventData,omitempty"`
	APIVersion                   *string             `json:"apiVersion,omitempty"`
	AwsRegion                    string              `json:"awsRegion"`
	ErrorCode                    *string             `json:"errorCode,omitempty"`
	ErrorMessage                 *string             `json:"errorMessage,omitempty"`
	EventID                      string              `json:"eventID"`
	EventName                    string              `json:"eventName"`
	EventSource                  string              `json:"eventSource"`
	EventTime                    helpers.UnixMillis  `json:"eventTime"`
	EventType                    string              `json:"eventType"`
	EventVersion                 string              `json:"eventVersion"`
	ManagementEvent              *bool               `json:"managementEvent,omitempty"`
	ReadOnly                     *bool               `json:"readOnly,omitempty"`
	RecipientAccountId           string              `json:"recipientAccountId,omitempty"`
	RequestID                    *string             `json:"requestID,omitempty"`
	RequestParameters            *helpers.JSONString `json:"requestParameters,omitempty"`
	Resources                    []*Resource         `json:"resources,omitempty"`
	ResponseElements             *helpers.JSONString `json:"responseElements,omitempty"`
	ServiceEventDetails          *helpers.JSONString `json:"serviceEventDetails,omitempty"`
	SharedEventID                *string             `json:"sharedEventID,omitempty"`
	SourceIPAddress              *string             `json:"sourceIPAddress,omitempty"`
	UserAgent                    *string             `json:"userAgent,omitempty"`
	UserIdentity                 UserIdentity        `json:"userIdentity"`
	VpcEndpointId                string              `json:"vpcEndpointId,omitempty"`
	EventCategory                string              `json:"eventCategory,omitempty"`
	SessionCredentialFromConsole *string             `json:"sessionCredentialFromConsole,omitempty"`
	EdgeDeviceDetails            *helpers.JSONString `json:"edgeDeviceDetails,omitempty"`
	TLSDetails                   *TLSDetails         `json:"tlsDetails,omitempty"`
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
	MfaAuthenticated *string             `json:"mfaAuthenticated,omitempty" parquet:"name=mfa_authenticated, type=BYTE_ARRAY, convertedtype=UTF8"`
	CreationDate     *helpers.UnixMillis `json:"creationDate,omitempty" parquet:"name=creation_date, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
}

type SessionIssuer struct {
	Type        *string `json:"type,omitempty" parquet:"name=type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	PrincipalId *string `json:"principalId,omitempty" parquet:"name=principal_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	ARN         *string `json:"arn,omitempty" parquet:"name=arn, type=BYTE_ARRAY, convertedtype=UTF8"`
	AccountId   *string `json:"accountId,omitempty" parquet:"name=account_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	UserName    *string `json:"userName,omitempty" parquet:"name=user_name, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type WebIdFederationData struct {
	FederatedProvider *string             `json:"federatedProvider,omitempty" parquet:"name=federated_provider, type=BYTE_ARRAY, convertedtype=UTF8"`
	Attributes        *helpers.JSONString `json:"attributes,omitempty" parquet:"name=attributes, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type TLSDetails struct {
	TLSVersion  *string `json:"tlsVersion,omitempty" parquet:"name=tls_version, type=BYTE_ARRAY, convertedtype=UTF8"`
	CipherSuite *string `json:"cipherSuite,omitempty" parquet:"name=cipher_suite, type=BYTE_ARRAY, convertedtype=UTF8"`
	//ClientProvidedHostHeader *string `json:"clientProvidedHostHeader,omitempty" parquet:"name=client_provided_host_header, type=BYTE_ARRAY, convertedtype=UTF8"`
	ClientProvidedHostHeader *string `json:"clientProvidedHostHeader,omitempty" parquet:"name=client_provided_host_header, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN"`
}
