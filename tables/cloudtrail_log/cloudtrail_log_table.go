package cloudtrail_log

import (
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-aws/sources/cloudwatch"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-aws/tables"

	"github.com/rs/xid"

	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const CloudTrailLogTableIdentifier = "aws_cloudtrail_log"

// CloudTrailLogTable - table for CloudTrailLog logs
type CloudTrailLogTable struct{}

func (t *CloudTrailLogTable) GetSourceMetadata() []*table.SourceMetadata[*CloudTrailLog] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"),
	}

	return []*table.SourceMetadata[*CloudTrailLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewCloudTrailLogExtractor()),
			},
		},
		{
			// any other artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewCloudTrailLogExtractor()),
			},
		},
		{
			SourceName: cloudwatch.AwsCloudwatchSourceIdentifier,
			Mapper:     &CloudTrailMapper{},
		},
	}
}

// Identifier implements table.Table
func (t *CloudTrailLogTable) Identifier() string {
	return CloudTrailLogTableIdentifier
}

// EnrichRow implements table.Table
func (t *CloudTrailLogTable) EnrichRow(row *CloudTrailLog, sourceEnrichmentFields schema.SourceEnrichment) (*CloudTrailLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpTimestamp = *row.EventTime
	row.TpIngestTimestamp = time.Now()

	if row.SourceIPAddress != nil {
		row.TpSourceIP = row.SourceIPAddress
		row.TpIps = append(row.TpIps, *row.SourceIPAddress)
	}
	for _, resource := range row.Resources {
		if resource.ARN != nil {
			newAkas := tables.AwsAkasFromArn(*resource.ARN)
			row.TpAkas = append(row.TpAkas, newAkas...)
		}
	}
	// If it's an AKIA, then record that as an identity. Do not record ASIA*
	// keys etc.
	if row.UserIdentity.AccessKeyId != nil {
		if strings.HasPrefix(*row.UserIdentity.AccessKeyId, "AKIA") {
			row.TpUsernames = append(row.TpUsernames, *row.UserIdentity.AccessKeyId)
		}
	}
	if row.UserIdentity.UserName != nil {
		row.TpUsernames = append(row.TpUsernames, *row.UserIdentity.UserName)
	}

	// Hive fields
	row.TpIndex = row.RecipientAccountId
	// convert to date in format yy-mm-dd
	row.TpDate = row.EventTime.Truncate(24 * time.Hour)

	return row, nil
}

func (c *CloudTrailLogTable) GetDescription() string {
	return "AWS CloudTrail logs capture API activity and user actions within your AWS account."
}
