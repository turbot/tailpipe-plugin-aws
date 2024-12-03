package tables

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const LambdaLogTableIdentifier = "aws_lambda_log"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.LambdaLog, *LambdaLogTableConfig, *LambdaLogTable]()
}

type LambdaLogTable struct{}

func (c *LambdaLogTable) Identifier() string {
	return LambdaLogTableIdentifier
}

func (c *LambdaLogTable) GetSourceMetadata(_ *LambdaLogTableConfig) []*table.SourceMetadata[*rows.LambdaLog] {
	return []*table.SourceMetadata[*rows.LambdaLog]{
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

func (c *LambdaLogTable) EnrichRow(row *rows.LambdaLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.LambdaLog, error) {
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	if !row.TpTimestamp.IsZero() {
		row.TpTimestamp = *row.Timestamp
		row.TpDate = row.Timestamp.Truncate(24 * time.Hour)
	}

	return row, nil
}
