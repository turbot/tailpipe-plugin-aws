package lambda_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type LambdaLog struct {
	schema.CommonFields

	Timestamp    *time.Time `json:"timestamp,omitempty"`
	RequestID    *string    `json:"request_id,omitempty"`
	LogType      *string    `json:"log_type,omitempty"`
	LogLevel     *string    `json:"log_level,omitempty"`
	Message      *string    `json:"message,omitempty"`
	LogGroupName *string    `json:"log_group_name,omitempty"`

	// Report Specific Fields
	Duration       *float64 `json:"duration,omitempty"`
	BilledDuration *float64 `json:"billed_duration,omitempty"`
	MemorySize     *int     `json:"memory_size,omitempty"`
	MaxMemoryUsed  *int     `json:"max_memory_used,omitempty"`
}
