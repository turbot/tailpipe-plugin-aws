// Package cloudwatch provides functionality to collect logs from AWS CloudWatch
package cloudwatch

import "fmt"

// AwsCloudWatchSourceConfig defines the configuration parameters for collecting logs from AWS CloudWatch.
// It specifies which log group to collect from, optionally filters log streams by prefix,
// and allows specifying the AWS region to connect to.
type AwsCloudWatchSourceConfig struct {
	// LogGroupName is the name of the CloudWatch log group to collect logs from (required)
	LogGroupName string `hcl:"log_group_name"`
	// LogStreamPrefix optionally filters log streams by their name prefix
	LogStreamPrefix *string `hcl:"log_stream_prefix"`
	// Region specifies the AWS region where the log group exists
	// If not provided, defaults to us-east-1
	Region *string `hcl:"region"`
}

// Validate checks if the configuration is valid.
// It ensures that the required LogGroupName field is provided and not empty.
func (c *AwsCloudWatchSourceConfig) Validate() error {
	if c.LogGroupName == "" {
		return fmt.Errorf("log_group_name is required and cannot be empty")
	}
	return nil
}

// Identifier returns the unique identifier for this source type.
// This is used to identify the source type in the plugin system.
func (c *AwsCloudWatchSourceConfig) Identifier() string {
	return AwsCloudwatchSourceIdentifier
}
