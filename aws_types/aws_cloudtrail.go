package aws_types

import (
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
)

type AWSCloudTrailBatch struct {
	Records []AWSCloudTrail `json:"Records"`
}

// TODO validate the struct does not have omitempty fields
type AWSCloudTrail struct {
	// embed required enrichment fields (be sure to skip in parquet)
	enrichment.CommonFields `parquet:"-"`

	AdditionalEventData *helpers.JSONString `json:"additionalEventData"`
	APIVersion          *string             `json:"apiVersion" parquet:"name=apiversion"`
	AwsRegion           string              `json:"awsRegion"`
	ErrorCode           *string             `json:"errorCode"`
	ErrorMessage        *string             `json:"errorMessage"`
	EventID             string              `json:"eventID"`
	EventName           string              `json:"eventName"`
	EventSource         string              `json:"eventSource"`
	EventTime           helpers.UnixMillis  `json:"eventTime"`
	EventType           string              `json:"eventType"`
	EventVersion        string              `json:"eventVersion"`
	ManagementEvent     *bool               `json:"managementEvent"`
	ReadOnly            *bool               `json:"readOnly"`
	RecipientAccountId  string              `json:"recipientAccountId"`
	RequestID           *string             `json:"requestID"`
	RequestParameters   *helpers.JSONString `json:"requestParameters"`
	Resources           []*Resource         `json:"resources"`
	ResponseElements    *helpers.JSONString `json:"responseElements"`
	ServiceEventDetails *helpers.JSONString `json:"serviceEventDetails"`
	SharedEventID       *string             `json:"sharedEventID"`
	SourceIPAddress     *string             `json:"sourceIPAddress"`
	UserAgent           *string             `json:"userAgent"`
	UserIdentity        UserIdentity        `json:"userIdentity"`
	VpcEndpointId       string              `json:"vpcEndpointId"`
	EventCategory       string              `json:"eventCategory"`
	// TODO - this looks like a bool, but is in the JSON as a string ... should we convert it?
	SessionCredentialFromConsole *string             `json:"sessionCredentialFromConsole"`
	EdgeDeviceDetails            *helpers.JSONString `json:"edgeDeviceDetails"`
	TLSDetails                   *TLSDetails         `json:"tlsDetails"`
}

type UserIdentity struct {
	Type             string          `json:"type"`
	PrincipalId      *string         `json:"principalId"`
	ARN              *string         `json:"arn"`
	AccountId        *string         `json:"accountId"`
	AccessKeyId      *string         `json:"accessKeyId"`
	UserName         *string         `json:"userName"`
	SessionContext   *SessionContext `json:"sessionContext"`
	InvokedBy        *string         `json:"invokedBy"`
	IdentityProvider *string         `json:"identityProvider"`
}

type Resource struct {
	ARN       *string `json:"ARN"`
	AccountId *string `json:"accountId"`
	Type      *string `json:"type"`
}

type SessionContext struct {
	Attributes          *SessionAttributes   `json:"attributes"`
	SessionIssuer       *SessionIssuer       `json:"sessionIssuer"`
	WebIdFederationData *WebIdFederationData `json:"webIdFederationData"`
	EC2RoleDelivery     *string              `json:"ec2RoleDelivery"`
}

type SessionAttributes struct {
	MfaAuthenticated *string             `json:"mfaAuthenticated"`
	CreationDate     *helpers.UnixMillis `json:"creationDate"`
}

type SessionIssuer struct {
	Type        *string `json:"type"`
	PrincipalId *string `json:"principalId"`
	ARN         *string `json:"arn"`
	AccountId   *string `json:"accountId"`
	UserName    *string `json:"userName"`
}

type WebIdFederationData struct {
	FederatedProvider *string             `json:"federatedProvider"`
	Attributes        *helpers.JSONString `json:"attributes"`
}

type TLSDetails struct {
	TLSVersion               *string `json:"tlsVersion"`
	CipherSuite              *string `json:"cipherSuite"`
	ClientProvidedHostHeader *string `json:"clientProvidedHostHeader"`
}
