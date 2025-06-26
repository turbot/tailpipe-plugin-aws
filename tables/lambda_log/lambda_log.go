package lambda_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type LambdaLog struct {
	schema.CommonFields

	Timestamp      *time.Time             `json:"timestamp,omitempty"`
	RequestID      *string                `json:"request_id,omitempty"`
	LogType        *string                `json:"log_type,omitempty"`
	LogLevel       *string                `json:"log_level,omitempty"`
	Message        *string                `json:"message,omitempty"`
	MessageJson    map[string]interface{} `json:"message_json,omitempty"`
	RawMessage     *string                `json:"raw_message,omitempty"`
	RawMessageJson map[string]interface{} `json:"raw_message_json,omitempty"`
	LogGroupName   *string                `json:"log_group_name,omitempty"`

	// Report Specific Fields
	Duration       *float64 `json:"duration,omitempty"`
	BilledDuration *float64 `json:"billed_duration,omitempty"`
	MemorySize     *int     `json:"memory_size,omitempty"`
	MaxMemoryUsed  *int     `json:"max_memory_used,omitempty"`
}
