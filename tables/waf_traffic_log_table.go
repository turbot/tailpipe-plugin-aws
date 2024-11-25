package tables

import (
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const WaftTrafficLogTableIdentifier = "aws_waf_traffic_log"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.WafTrafficLog, *WafTrafficLogTableConfig, *WafTrafficLogTable]()
}

// WafTrafficLogTable - table for Waf traffic logs
type WafTrafficLogTable struct{}

func (c *WafTrafficLogTable) SupportedSources(*WafTrafficLogTableConfig) []*table.SourceMetadata[*rows.WafTrafficLog] {
	// the default file layout for Waf traffic logs in S3
	defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
		FileLayout: utils.ToStringPointer("AWSLogs(?:/o-[a-z0-9]{8,12})?/(?P<account_id>\\d+)/WAFLogs/(?P<log_group_name>[a-zA-Z0-9-_]+)/(?P<region>[a-z0-9-]+)/(?P<year>\\d{4})/(?P<month>\\d{2})/(?P<day>\\d{2})/(?P<hour>\\d{2})/(?P<filename>.+\\.gz)"),
	}

	return []*table.SourceMetadata[*rows.WafTrafficLog]{
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			MapperFunc: mappers.NewWafMapper,
			Options:    []row_source.RowSourceOption{artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig), artifact_source.WithRowPerLine()},
		},
	}
}

func (c *WafTrafficLogTable) Identifier() string {
	return WaftTrafficLogTableIdentifier
}

// EnrichRow implements table.Table
func (c *WafTrafficLogTable) EnrichRow(row *rows.WafTrafficLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.WafTrafficLog, error) { // initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

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
