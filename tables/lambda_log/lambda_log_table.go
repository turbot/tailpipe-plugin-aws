package lambda_log

import (
	"regexp"
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

const LambdaLogTableIdentifier = "aws_lambda_log"

type LambdaLogTable struct{}

func (c *LambdaLogTable) Identifier() string {
	return LambdaLogTableIdentifier
}

func (c *LambdaLogTable) GetSourceMetadata() ([]*table.SourceMetadata[*LambdaLog], error) {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/lambda/%{DATA:function_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log"),
	}

	return []*table.SourceMetadata[*LambdaLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     &LambdaLogMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// S3 artifact source
			SourceName: cloudwatch_log_group.AwsCloudwatchLogGroupSourceIdentifier,
			Mapper:     &LambdaLogMapper{},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &LambdaLogMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}, nil
}

func (c *LambdaLogTable) EnrichRow(row *LambdaLog, sourceEnrichmentFields schema.SourceEnrichment) (*LambdaLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	if row.Timestamp != nil {
		row.TpTimestamp = *row.Timestamp
		row.TpDate = row.Timestamp.Truncate(24 * time.Hour)
	} else if !row.TpTimestamp.IsZero() {
		row.TpDate = row.TpTimestamp.Truncate(24 * time.Hour)
	}

	row.TpIndex = schema.DefaultIndex

	var arnRegex = regexp.MustCompile(`arn:aws:[^,\s'"\\]+`)

	seen := map[string]struct{}{}
	for _, match := range arnRegex.FindAllString(*row.Message, -1) {
		if _, exists := seen[match]; !exists {
			seen[match] = struct{}{}
			row.TpAkas = append(row.TpAkas, match)
		}
	}

	// TODO: Add enrichment fields

	return row, nil
}
