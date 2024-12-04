package sources

import (
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
)

// AwsCloudWatchSourceConfig is the configuration for an [AwsCloudWatchSource]
type AwsCloudWatchSourceConfig struct {
	artifact_source_config.ArtifactSourceConfigBase
	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	// the log group to collect
	LogGroupName string `hcl:"log_group_name"`
	// collect log streams with this prefix the log stream prefix
	LogStreamPrefix *string `hcl:"log_stream_prefix"`
	// the time range to collect for
	StartTimeString string `hcl:"start_time"`
	EndTimeString   string `hcl:"end_time"`
	StartTime       time.Time
	EndTime         time.Time

	Region *string `hcl:"region"`
}

func (c *AwsCloudWatchSourceConfig) Validate() error {
	// parse  start  and end time
	if c.StartTimeString == "" {
		return fmt.Errorf("start_time is required")
	}
	startTime, err := time.Parse(time.RFC3339, c.StartTimeString)
	if err != nil {
		return fmt.Errorf("invalid start_time: %v", err)
	}
	c.StartTime = startTime
	if c.EndTimeString == "" {
		return fmt.Errorf("end_time is required")

	}
	endTime, err := time.Parse(time.RFC3339, c.EndTimeString)
	if err != nil {
		return fmt.Errorf("invalid end_time: %v", err)
	}
	c.EndTime = endTime
	if c.StartTime.After(c.EndTime) {
		return fmt.Errorf("start_time must be before end_time")
	}

	return c.ArtifactSourceConfigBase.Validate()
}

func (c *AwsCloudWatchSourceConfig) Identifier() string {
	return AwsCloudwatchSourceIdentifier
}
