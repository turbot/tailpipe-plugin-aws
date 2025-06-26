package cloudtrail_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type CloudTrailLog struct {
	schema.CommonFields

	// CloudTrail specific fields
	EventVersion      string                 `json:"eventVersion"`
	UserIdentity      map[string]interface{} `json:"userIdentity"`
	EventTime         time.Time              `json:"eventTime"`
	EventSource       string                 `json:"eventSource"`
	EventName         string                 `json:"eventName"`
	AWSRegion         string                 `json:"awsRegion"`
	SourceIPAddress   string                 `json:"sourceIPAddress"`
	UserAgent         string                 `json:"userAgent"`
	RequestParameters map[string]interface{} `json:"requestParameters"`
	ResponseElements  map[string]interface{} `json:"responseElements"`
	RequestID         string                 `json:"requestID"`
	EventID           string                 `json:"eventID"`
	EventType         string                 `json:"eventType"`
	APIVersion        string                 `json:"apiVersion,omitempty"`
	ReadOnly          bool                   `json:"readOnly"`
	Resources         []struct {
		ResourceType string                 `json:"resourceType"`
		ResourceName string                 `json:"resourceName"`
		ARN          string                 `json:"ARN"`
		Tags         map[string]interface{} `json:"tags,omitempty"`
	} `json:"resources"`
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	AccountID    string `json:"recipientAccountId"`
}

// Validate validates the CloudTrailLog struct
func (c CloudTrailLog) Validate() error {
	return nil
}
