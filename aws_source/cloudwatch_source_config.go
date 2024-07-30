package aws_source

import (
	"github.com/turbot/tailpipe-plugin-sdk/artifact"
	"time"
)

// AwsCloudWatchSourceConfig is the configuration for an [CloudWatchArtifactSource]
type AwsCloudWatchSourceConfig struct {
	artifact.SourceConfigBase
	// TODO #confif connection based credentiuals mechanism
	AccessKey    string
	SecretKey    string
	SessionToken string

	// the log group to collect
	LogGroupName string

	// collect log streams with this prefixthe log stream prefix
	LogStreamPrefix *string

	// the time range to collect for
	StartTime time.Time
	EndTime   time.Time
}
