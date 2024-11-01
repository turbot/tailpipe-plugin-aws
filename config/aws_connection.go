package config

import "github.com/turbot/tailpipe-plugin-sdk/parse"

type AwsConnection struct {
	Regions               []string `hcl:"regions,optional"`
	DefaultRegion         *string  `hcl:"default_region"`
	Profile               *string  `hcl:"profile"`
	AccessKey             *string  `hcl:"access_key"`
	SecretKey             *string  `hcl:"secret_key"`
	SessionToken          *string  `hcl:"session_token"`
	MaxErrorRetryAttempts *int     `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int     `hcl:"min_error_retry_delay"`
	IgnoreErrorCodes      []string `hcl:"ignore_error_codes,optional"`
	EndpointUrl           *string  `hcl:"endpoint_url"`
	S3ForcePathStyle      *bool    `hcl:"s3_force_path_style"`
}

func NewAwsConnection() parse.Config {
	return &AwsConnection{}
}

func (c *AwsConnection) Validate() error {
	return nil
}
