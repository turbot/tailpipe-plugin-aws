package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type CloudTrailBatch struct {
	Records []CloudTrailLog `json:"Records"`
}

type CloudTrailLog struct {
	// embed required enrichment fields
	enrichment.CommonFields

	// TODO json fields should be snake case - these are probably camel case because
	// they are coming from the AWS SDK? in which case we should just use a mapper
	AdditionalEventData *types.JSONString `json:"additionalEventData,omitempty"`
	APIVersion          *string           `json:"apiVersion,omitempty"`
	AwsRegion           string            `json:"awsRegion"`
	ErrorCode           *string           `json:"errorCode,omitempty"`
	ErrorMessage        *string           `json:"errorMessage,omitempty"`
	EventID             string            `json:"eventID"`
	EventName           string            `json:"eventName"`
	EventSource         string            `json:"eventSource"`
	EventTime           *time.Time        `json:"eventTime"`
	EventType           string            `json:"eventType"`
	EventVersion        string            `json:"eventVersion"`
	ManagementEvent     *bool             `json:"managementEvent,omitempty" `
	ReadOnly            *bool             `json:"readOnly,omitempty"`
	RecipientAccountId  string            `json:"recipientAccountId,omitempty" `
	RequestID           *string           `json:"requestID,omitempty" `
	RequestParameters   *types.JSONString `json:"requestParameters,omitempty" `
	Resources           []*Resource       `json:"resources,omitempty" `
	ResponseElements    *types.JSONString `json:"responseElements,omitempty" `
	ServiceEventDetails *types.JSONString `json:"serviceEventDetails,omitempty"`
	SharedEventID       *string           `json:"sharedEventID,omitempty"`
	SourceIPAddress     *string           `json:"sourceIPAddress,omitempty"`
	UserAgent           *string           `json:"userAgent,omitempty"`
	UserIdentity        UserIdentity      `json:"userIdentity"`
	VpcEndpointId       string            `json:"vpcEndpointId,omitempty"`
	EventCategory       string            `json:"eventCategory,omitempty"`
	// TODO - this looks like a bool, but is in the JSON as a string ... should we convert it?
	SessionCredentialFromConsole *string           `json:"sessionCredentialFromConsole,omitempty"`
	EdgeDeviceDetails            *types.JSONString `json:"edgeDeviceDetails,omitempty"`
	TLSDetails                   *TLSDetails       `json:"tlsDetails,omitempty"`
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
