package s3_bucket

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
)

// AwsS3BucketSourceConfig is the configuration for an [AwsS3BucketSource]
type AwsS3BucketSourceConfig struct {
	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`
	artifact_source_config.ArtifactSourceConfigImpl

	Bucket string  `hcl:"bucket"`
	Prefix *string `hcl:"prefix,optional"`
}

func (c *AwsS3BucketSourceConfig) Validate() error {
	if c.Bucket == "" {
		return fmt.Errorf("bucket is required and cannot be empty")
	}

	return nil
}

func (c *AwsS3BucketSourceConfig) Identifier() string {
	return AwsS3BucketSourceIdentifier
}
