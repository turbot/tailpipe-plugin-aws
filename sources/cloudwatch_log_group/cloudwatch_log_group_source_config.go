// Package cloudwatch provides functionality to collect logs from AWS CloudWatch
package cloudwatch_log_group

import "fmt"

// AwsCloudWatchLogGroupSourceConfig defines the configuration parameters for collecting logs from AWS CloudWatch log groups.
// It specifies which log group to collect from, optionally filters log streams by prefix,
// and allows specifying the AWS region to connect to.
type AwsCloudWatchLogGroupSourceConfig struct {
	// LogGroupName is the name of the CloudWatch log group to collect logs from (required)
	LogGroupName string `hcl:"log_group_name"`
	// LogStreamNames optionally filters log streams by their names. Supports wildcards (*).
	// If not specified, logs from all available streams will be collected.
	// Example: ["456789012345_CloudTrail_*", "123456789012_CloudTrail_us-east-1"]
	LogStreamNames []string `hcl:"log_stream_names,optional"`
	// Region specifies the AWS region where the log group exists
	// If not provided, defaults to us-east-1
	Region *string `hcl:"region"`
}

// Validate checks if the configuration is valid.
// It ensures that the required LogGroupName field is provided and not empty.
func (c *AwsCloudWatchLogGroupSourceConfig) Validate() error {
	if c.LogGroupName == "" {
		return fmt.Errorf("log_group_name is required and cannot be empty")
	}
	if c.Region == nil {
		return fmt.Errorf("region is required and cannot be empty")
	}
	return nil
}

// Identifier returns the unique identifier for this source type.
// This is used to identify the source type in the plugin system.
func (c *AwsCloudWatchLogGroupSourceConfig) Identifier() string {
	return AwsCloudwatchLogGroupSourceIdentifier
}
