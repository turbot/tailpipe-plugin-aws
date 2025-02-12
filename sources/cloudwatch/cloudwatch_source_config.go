package cloudwatch

import "fmt"

// AwsCloudWatchSourceConfig is the configuration for an [AwsCloudWatchSource]
type AwsCloudWatchSourceConfig struct {
	LogGroupName    string `hcl:"log_group_name"`
	LogStreamPrefix *string `hcl:"log_stream_prefix"`
	Region          *string `hcl:"region"`
}

func (c *AwsCloudWatchSourceConfig) Validate() error {
	if c.LogGroupName == "" {
		return fmt.Errorf("log group is required and cannot be empty")
	}
	return nil
}

func (c *AwsCloudWatchSourceConfig) Identifier() string {
	return AwsCloudwatchSourceIdentifier
}
