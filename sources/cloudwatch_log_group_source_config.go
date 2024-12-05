package sources

import (
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
)

// AwsCloudWatchSourceConfig is the configuration for an [AwsCloudWatchSource]
type AwsCloudWatchSourceConfig struct {
	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`

	LogGroupName    string  `hcl:"log_group_name"`
	LogStreamPrefix *string `hcl:"log_stream_prefix"`
	StartTimeString string  `hcl:"start_time"`
	EndTimeString   *string `hcl:"end_time"`
	StartTime       time.Time
	EndTime         *time.Time
	Region          *string `hcl:"region"`
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
	if c.EndTimeString == nil {
		endTime, err := time.Parse(time.RFC3339, *c.EndTimeString)
		if err != nil {
			return fmt.Errorf("invalid end_time: %v", err)
		}
		c.EndTime = &endTime
		if c.StartTime.After(endTime) {
			return fmt.Errorf("start_time must be before end_time")
		}
	}

	return nil
}

func (c *AwsCloudWatchSourceConfig) Identifier() string {
	return AwsCloudwatchSourceIdentifier
}
