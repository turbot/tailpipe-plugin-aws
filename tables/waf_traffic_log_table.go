package tables

import (
	"strings"
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

const WaftTrafficLogTableIdentifier = "aws_waf_traffic_log"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.WafTrafficLog, *WafTrafficLogTable]()
}

// WafTrafficLogTable - table for Waf traffic logs
type WafTrafficLogTable struct{}

func (c *WafTrafficLogTable) GetSourceMetadata() []*table.SourceMetadata[*rows.WafTrafficLog] {
	// the default file layout for Waf traffic logs in S3
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/WAFLogs/%{DATA:log_group_name}/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{DATA}.gz"),
	}

	return []*table.SourceMetadata[*rows.WafTrafficLog]{
		{
			// S3 artifact source
			SourceName: sources.AwsS3BucketSourceIdentifier,
			Mapper:     &mappers.WafMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &mappers.WafMapper{},
			Options:    []row_source.RowSourceOption{artifact_source.WithRowPerLine()},
		},
	}
}

func (c *WafTrafficLogTable) Identifier() string {
	return WaftTrafficLogTableIdentifier
}

// EnrichRow implements table.Table
func (c *WafTrafficLogTable) EnrichRow(row *rows.WafTrafficLog, sourceEnrichmentFields schema.SourceEnrichment) (*rows.WafTrafficLog, error) { // initialize the enrichment fields to any fields provided by the source
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpTimestamp = *row.Timestamp
	row.TpIngestTimestamp = time.Now()
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)

	// TODO: #enrichment figure out correct field for TpIndex
	// Hive fields
	// Check if row.HttpSourceId is not nil before dereferencing
	if row.HttpSourceId != nil {
		row.TpIndex = strings.ReplaceAll(*row.HttpSourceId, "/", `_`)
	} else {
		// Handle the case where HttpSourceId is nil, if needed
		row.TpIndex = "" // or assign a default value if desired
	}

	return row, nil
}
