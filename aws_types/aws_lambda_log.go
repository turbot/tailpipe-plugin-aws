package aws_types

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type AwsLambdaLog struct {
	enrichment.CommonFields

	Timestamp *time.Time `json:"timestamp,omitempty"`
	RequestID *string    `json:"request_id,omitempty"`
	LogType   *string    `json:"log_type,omitempty"`
	LogLevel  *string    `json:"log_level,omitempty"`
	Message   *string    `json:"message,omitempty"`

	// Report Specific Fields
	Duration       *float64 `json:"duration,omitempty"`
	BilledDuration *float64 `json:"billed_duration,omitempty"`
	MemorySize     *int     `json:"memory_size,omitempty"`
	MaxMemoryUsed  *int     `json:"max_memory_used,omitempty"`
}
