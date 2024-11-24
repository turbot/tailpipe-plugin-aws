package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type CloudTrailLogBatch struct {
	Records []CloudTrailLog `json:"Records"`
}

type CloudTrailLog struct {
	// embed required enrichment fields
	enrichment.CommonFields

	// json tags for marshalling to/from the source & parquet tags handle the parquet column names for the table
	AdditionalEventData          *map[string]interface{} `json:"additionalEventData,omitempty" parquet:"name=additional_event_data, type=JSON"`
	APIVersion                   *string                 `json:"apiVersion,omitempty" parquet:"name=api_version"`
	AwsRegion                    string                  `json:"awsRegion" parquet:"name=aws_region"`
	ErrorCode                    *string                 `json:"errorCode,omitempty" parquet:"name=error_code"`
	ErrorMessage                 *string                 `json:"errorMessage,omitempty" parquet:"name=error_message"`
	EventID                      string                  `json:"eventID" parquet:"name=event_id"`
	EventName                    string                  `json:"eventName" parquet:"name=event_name"`
	EventSource                  string                  `json:"eventSource" parquet:"name=event_source"`
	EventTime                    *time.Time              `json:"eventTime" parquet:"name=event_time"`
	EventType                    string                  `json:"eventType" parquet:"name=event_type"`
	EventVersion                 string                  `json:"eventVersion" parquet:"name=event_version"`
	ManagementEvent              *bool                   `json:"managementEvent,omitempty" parquet:"name=management_event"`
	ReadOnly                     *bool                   `json:"readOnly,omitempty" parquet:"name=read_only"`
	RecipientAccountId           string                  `json:"recipientAccountId,omitempty" parquet:"name=recipient_account_id"`
	RequestID                    *string                 `json:"requestID,omitempty" parquet:"name=request_id"`
	RequestParameters            *map[string]interface{} `json:"requestParameters,omitempty" parquet:"name=request_parameters, type=JSON"`
	Resources                    []*Resource             `json:"resources,omitempty" parquet:"name=resources, type=JSON"`
	ResponseElements             *map[string]interface{} `json:"responseElements,omitempty" parquet:"name=response_elements, type=JSON"`
	ServiceEventDetails          *map[string]interface{} `json:"serviceEventDetails,omitempty" parquet:"name=service_event_details, type=JSON"`
	SharedEventID                *string                 `json:"sharedEventID,omitempty" parquet:"name=shared_event_id"`
	SourceIPAddress              *string                 `json:"sourceIPAddress,omitempty" parquet:"name=source_ip_address"`
	UserAgent                    *string                 `json:"userAgent,omitempty" parquet:"name=user_agent"`
	UserIdentity                 UserIdentity            `json:"userIdentity" parquet:"name=user_identity, type=JSON"`
	VpcEndpointId                string                  `json:"vpcEndpointId,omitempty" parquet:"name=vpc_endpoint_id"`
	EventCategory                string                  `json:"eventCategory,omitempty" parquet:"name=event_category"`
	SessionCredentialFromConsole *string                 `json:"sessionCredentialFromConsole,omitempty" parquet:"name=session_credential_from_console"`
	EdgeDeviceDetails            *map[string]interface{} `json:"edgeDeviceDetails,omitempty" parquet:"name=edge_device_details, type=JSON"`
	TLSDetails                   *TLSDetails             `json:"tlsDetails,omitempty" parquet:"name=tls_details, type=JSON"`
}

type UserIdentity struct {
	Type             string          `json:"type"`
	PrincipalId      *string         `json:"principalId,omitempty"`
	ARN              *string         `json:"arn,omitempty"`
	AccountId        *string         `json:"accountId,omitempty"`
	AccessKeyId      *string         `json:"accessKeyId,omitempty"`
	UserName         *string         `json:"userName,omitempty"`
	SessionContext   *SessionContext `json:"sessionContext,omitempty"`
	InvokedBy        *string         `json:"invokedBy,omitempty"`
	IdentityProvider *string         `json:"identityProvider,omitempty"`
}

type Resource struct {
	ARN       *string `json:"ARN,omitempty"`
	AccountId *string `json:"accountId,omitempty"`
	Type      *string `json:"type,omitempty"`
}

type SessionContext struct {
	Attributes          *SessionAttributes   `json:"attributes,omitempty"`
	SessionIssuer       *SessionIssuer       `json:"sessionIssuer,omitempty"`
	WebIdFederationData *WebIdFederationData `json:"webIdFederationData,omitempty"`
	EC2RoleDelivery     *string              `json:"ec2RoleDelivery,omitempty"`
}

type SessionAttributes struct {
	MfaAuthenticated *string           `json:"mfaAuthenticated,omitempty"`
	CreationDate     *types.UnixMillis `json:"creationDate,omitempty"`
}

type SessionIssuer struct {
	Type        *string `json:"type,omitempty"`
	PrincipalId *string `json:"principalId,omitempty"`
	ARN         *string `json:"arn,omitempty"`
	AccountId   *string `json:"accountId,omitempty"`
	UserName    *string `json:"userName,omitempty"`
}

type WebIdFederationData struct {
	FederatedProvider *string           `json:"federatedProvider,omitempty"`
	Attributes        *types.JSONString `json:"attributes,omitempty"`
}

type TLSDetails struct {
	TLSVersion               *string `json:"tlsVersion,omitempty"`
	CipherSuite              *string `json:"cipherSuite,omitempty"`
	ClientProvidedHostHeader *string `json:"clientProvidedHostHeader,omitempty"`
}
