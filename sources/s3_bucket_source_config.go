package sources

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
)

// AwsS3BucketSourceConfig is the configuration for an [AwsS3BucketSource]
type AwsS3BucketSourceConfig struct {
	// required to allow partial decoding
	Remain hcl.Body `hcl:",remain" json:"-"`
	artifact_source_config.ArtifactSourceConfigImpl

	Bucket               string  `hcl:"bucket"`
	Prefix               *string `hcl:"prefix"`
	Region               *string `hcl:"region"`
	StartAfterKey        *string `hcl:"start_after_key"`
	LexicographicalOrder bool    `hcl:"lexicographical_order,optional"`
}

func (c *AwsS3BucketSourceConfig) Validate() error {
	if c.Bucket == "" {
		return fmt.Errorf("bucket is required and cannot be empty")
	}

	if c.Region != nil && !isValidAWSRegion(*c.Region) {
		return fmt.Errorf("invalid AWS region '%s'", *c.Region)
	}

	return nil
}

func (c *AwsS3BucketSourceConfig) Identifier() string {
	return AwsS3BucketSourceIdentifier
}

// IsValidAWSRegion checks if the given region is in the FORMAT of an AWS region.
func isValidAWSRegion(region string) bool {
	// TODO: #refactor should we ensure it is an actual AWS region via CLI/API?
	pattern := `^(us|eu|ap|af|me|sa|ca)-(west|east|central|south|north|northeast|southeast|northwest|southwest)-\d+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(region)
}
