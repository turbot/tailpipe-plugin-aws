package cloudtrail

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type CloudTrailLogBatch struct {
	Records []CloudTrailLog `json:"Records"`
}

type CloudTrailLog struct {
	// embed required enrichment fields
	schema.CommonFields

	// json tags for marshalling to/from the source & parquet tags handle the parquet column names for the table
	AdditionalEventData          *map[string]interface{} `json:"additionalEventData,omitempty" parquet:"name=additional_event_data, type=JSON"`
	APIVersion                   *string                 `json:"apiVersion,omitempty" parquet:"name=api_version"`
	AwsRegion                    string                  `json:"awsRegion" parquet:"name=aws_region"`
	EdgeDeviceDetails            *map[string]interface{} `json:"edgeDeviceDetails,omitempty" parquet:"name=edge_device_details, type=JSON"`
	ErrorCode                    *string                 `json:"errorCode,omitempty" parquet:"name=error_code"`
	ErrorMessage                 *string                 `json:"errorMessage,omitempty" parquet:"name=error_message"`
	EventCategory                string                  `json:"eventCategory,omitempty" parquet:"name=event_category"`
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
	SessionCredentialFromConsole *string                 `json:"sessionCredentialFromConsole,omitempty" parquet:"name=session_credential_from_console"`
	SharedEventID                *string                 `json:"sharedEventID,omitempty" parquet:"name=shared_event_id"`
	SourceIPAddress              *string                 `json:"sourceIPAddress,omitempty" parquet:"name=source_ip_address"`
	TLSDetails                   *TLSDetails             `json:"tlsDetails,omitempty" parquet:"name=tls_details"`
	UserAgent                    *string                 `json:"userAgent,omitempty" parquet:"name=user_agent"`
	UserIdentity                 UserIdentity            `json:"userIdentity" parquet:"name=user_identity"`
	VpcEndpointId                string                  `json:"vpcEndpointId,omitempty" parquet:"name=vpc_endpoint_id"`
}

type UserIdentity struct {
	AccessKeyId      *string         `json:"accessKeyId,omitempty" parquet:"name=access_key_id"`
	AccountId        *string         `json:"accountId,omitempty" parquet:"name=account_id"`
	ARN              *string         `json:"arn,omitempty" parquet:"name=arn"`
	IdentityProvider *string         `json:"identityProvider,omitempty" parquet:"name=identity_provider"`
	InvokedBy        *string         `json:"invokedBy,omitempty" parquet:"name=invoked_by"`
	PrincipalId      *string         `json:"principalId,omitempty" parquet:"name=principal_id"`
	SessionContext   *SessionContext `json:"sessionContext,omitempty" parquet:"name=session_context"`
	Type             string          `json:"type" parquet:"name=type"`
	UserName         *string         `json:"userName,omitempty" parquet:"name=user_name"`
}

type Resource struct {
	AccountId *string `json:"accountId,omitempty"`
	ARN       *string `json:"ARN,omitempty"`
	Type      *string `json:"type,omitempty"`
}

type SessionContext struct {
	Attributes          *SessionAttributes   `json:"attributes,omitempty" parquet:"name=attributes"`
	EC2RoleDelivery     *string              `json:"ec2RoleDelivery,omitempty" parquet:"name=ec2_role_delivery"`
	SessionIssuer       *SessionIssuer       `json:"sessionIssuer,omitempty" parquet:"name=session_issuer"`
	WebIdFederationData *WebIdFederationData `json:"webIdFederationData,omitempty" parquet:"name=web_id_federation_data"`
}

type SessionAttributes struct {
	CreationDate     *types.UnixMillis `json:"creationDate,omitempty" parquet:"name=creation_date"`
	MfaAuthenticated *string           `json:"mfaAuthenticated,omitempty" parquet:"name=mfa_authenticated"`
}

type SessionIssuer struct {
	AccountId   *string `json:"accountId,omitempty" parquet:"name=account_id"`
	ARN         *string `json:"arn,omitempty" parquet:"name=arn"`
	PrincipalId *string `json:"principalId,omitempty" parquet:"name=principal_id"`
	Type        *string `json:"type,omitempty" parquet:"name=type"`
	UserName    *string `json:"userName,omitempty" parquet:"name=user_name"`
}

type WebIdFederationData struct {
	Attributes        *types.JSONString `json:"attributes,omitempty" parquet:"name=attributes, type=JSON"`
	FederatedProvider *string           `json:"federatedProvider,omitempty" parquet:"name=federated_provider"`
}

type TLSDetails struct {
	CipherSuite              *string `json:"cipherSuite,omitempty" parquet:"name=cipher_suite"`
	ClientProvidedHostHeader *string `json:"clientProvidedHostHeader,omitempty" parquet:"name=client_provided_host_header"`
	TLSVersion               *string `json:"tlsVersion,omitempty" parquet:"name=tls_version"`
}

func (c *CloudTrailLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"additional_event_data":           "Additional information about the event that is specific to the service being called.",
		"api_version":                     "The API version associated with the event.",
		"aws_region":                      "The AWS region where the event originated.",
		"edge_device_details":             "Details about an edge device involved in the event, in JSON format.",
		"error_code":                      "The error code returned, if the event resulted in an error.",
		"error_message":                   "The error message returned, if the event resulted in an error.",
		"event_category":                  "The category of the event, such as 'Management' or 'Data'.",
		"event_id":                        "A unique identifier for the event.",
		"event_name":                      "The name of the API operation that was invoked.",
		"event_source":                    "The AWS service that the request was made to, such as 'ec2.amazonaws.com'.",
		"event_time":                      "The date and time the event occurred, in ISO 8601 format.",
		"event_type":                      "The type of event (e.g., 'AwsApiCall', 'AwsServiceEvent').",
		"event_version":                   "The version of the event log schema.",
		"management_event":                "Indicates whether the event is a management event (true or false).",
		"read_only":                       "Indicates whether the request was a read-only operation (true or false).",
		"recipient_account_id":            "The AWS account ID that received the request.",
		"request_id":                      "The ID of the request associated with the event.",
		"request_parameters":              "The request parameters sent with the request, in JSON format.",
		"resources":                       "A list of resources that were affected by the event, including ARNs and resource types.",
		"response_elements":               "The response elements returned by the service in response to the request, in JSON format.",
		"service_event_details":           "Details about a service event, in JSON format.",
		"session_credential_from_console": "Indicates whether the session credential originated from the AWS Management Console.",
		"shared_event_id":                 "An identifier for shared events, when multiple entries represent the same event.",
		"source_ip_address":               "The IP address from which the request was made.",
		"tls_details":                     "Details about the TLS connection, if applicable.",
		"user_agent":                      "The user agent string of the client that made the request.",
		"user_identity":                   "Details about the IAM identity that made the request, including user, role, or service.",
		"vpc_endpoint_id":                 "The ID of the VPC endpoint through which the request was made, if applicable.",

		// Override table specific tp_* column descriptions
		"tp_akas":      "Resource ARNs associated with the event.",
		"tp_index":     "The AWS account ID that received the request.",
		"tp_ips":       "IP addresses associated with the event, including the source IP address.",
		"tp_timestamp": "The date and time the event occurred, in ISO 8601 format.",
		"tp_usernames": "Usernames or access key IDs associated with the event.",
	}
}
