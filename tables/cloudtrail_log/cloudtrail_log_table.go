package cloudtrail_log

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/cloudwatch_log_group"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const CloudTrailLogTableIdentifier = "aws_cloudtrail_log"

type CloudTrailLogTable struct{}

func (c CloudTrailLogTable) Identifier() string {
	return CloudTrailLogTableIdentifier
}

func (c CloudTrailLogTable) GetSourceMetadata() ([]*table.SourceMetadata[CloudTrailLog], error) {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"),
	}

	return []*table.SourceMetadata[CloudTrailLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     &CloudTrailMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// CloudWatch source
			SourceName: cloudwatch_log_group.AwsCloudwatchLogGroupSourceIdentifier,
			Mapper:     &CloudTrailMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &CloudTrailMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}, nil
}

func (c CloudTrailLogTable) EnrichRow(row CloudTrailLog, sourceEnrichmentFields schema.SourceEnrichment) (CloudTrailLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	row.TpTimestamp = row.EventTime
	row.TpDate = row.EventTime.Truncate(24 * time.Hour)

	// Set tp_index to account ID
	row.TpIndex = row.AccountID

	// Add source IP to tp_ips
	if row.SourceIPAddress != "" {
		row.TpSourceIP = &row.SourceIPAddress
		row.TpIps = append(row.TpIps, row.SourceIPAddress)
	}

	// Add resource ARNs to tp_akas
	for _, resource := range row.Resources {
		if resource.ARN != "" {
			row.TpAkas = append(row.TpAkas, resource.ARN)
		}
	}

	return row, nil
}

func (c CloudTrailLogTable) GetDescription() string {
	return "AWS CloudTrail logs record detailed information about API calls and resource changes in your AWS account, helping track user activity, security analysis, and compliance auditing."
}
