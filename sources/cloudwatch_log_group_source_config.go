package sources

// AwsCloudWatchSourceConfig is the configuration for an [AwsCloudWatchSource]
type AwsCloudWatchSourceConfig struct {
	LogGroupName    string  `hcl:"log_group_name"`
	LogStreamPrefix *string `hcl:"log_stream_prefix"`
	Region          *string `hcl:"region"`
}

func (c *AwsCloudWatchSourceConfig) Validate() error {
	return nil
}

func (c *AwsCloudWatchSourceConfig) Identifier() string {
	return AwsCloudwatchSourceIdentifier
}
