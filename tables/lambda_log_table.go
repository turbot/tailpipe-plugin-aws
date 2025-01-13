package tables

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-aws/sources"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const LambdaLogTableIdentifier = "aws_lambda_log"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.LambdaLog, *LambdaLogTable]()
}

type LambdaLogTable struct{}

func (c *LambdaLogTable) Identifier() string {
	return LambdaLogTableIdentifier
}

func (c *LambdaLogTable) GetSourceMetadata() []*table.SourceMetadata[*rows.LambdaLog] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/%{NUMBER:account_id}/lambda/%{DATA:function_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA:function_name}-%{TIMESTAMP_ISO8601:timestamp}-%{DATA:suffix}.log"),
	}

	return []*table.SourceMetadata[*rows.LambdaLog]{
		{
			// S3 artifact source
			SourceName: sources.AwsS3BucketSourceIdentifier,
			Mapper:     &mappers.LambdaLogMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &mappers.LambdaLogMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *LambdaLogTable) EnrichRow(row *rows.LambdaLog, sourceEnrichmentFields schema.SourceEnrichment) (*rows.LambdaLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	if !row.TpTimestamp.IsZero() {
		row.TpTimestamp = *row.Timestamp
		row.TpDate = row.Timestamp.Truncate(24 * time.Hour)
	}

	return row, nil
}
