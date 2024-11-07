package tables

import (
	"context"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

// WafTrafficLogTable - table for Waf traffic logs
type WafTrafficLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[rows.WafTrafficLog, *WafTrafficLogTableConfig, *config.AwsConnection]
}

func NewWafTrafficLogTable() table.Table {
	return &WafTrafficLogTable{}
}

func (c *WafTrafficLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
	// call base init
	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
		return err
	}

	c.initMapper()
	return nil
}

func (c *WafTrafficLogTable) initMapper() {
	// TODO switch on source

	// if the source is an artifact source, we need a mapper
	c.Mapper = mappers.NewWafMapper()
}

// Identifier implements table.Table
func (c *WafTrafficLogTable) Identifier() string {
	return "aws_waf_traffic_log"
}

// GetSourceOptions returns any options which should be passed to the given source type
func (c *WafTrafficLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	var opts []row_source.RowSourceOption

	switch sourceType {
	case artifact_source.AwsS3BucketSourceIdentifier:
		// the default file layout for Waf traffic logs in S3
		defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
			// TODO #config finalise default Waf traffic file layout
			FileLayout: utils.ToStringPointer("waf-access-log-sample_(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})(?P<hour>\\d{2})(?P<minute>\\d{2})(?P<second>\\d{2}).gz"),
		}
		opts = append(opts, artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig), artifact_source.WithRowPerLine())
	}

	return opts
}

// GetRowSchema implements table.Table
func (c *WafTrafficLogTable) GetRowSchema() any {
	return rows.WafTrafficLog{}
}

func (c *WafTrafficLogTable) GetConfigSchema() parse.Config {
	return &WafTrafficLogTableConfig{}
}

// EnrichRow implements table.Table
func (c *WafTrafficLogTable) EnrichRow(row rows.WafTrafficLog, sourceEnrichmentFields *enrichment.CommonFields) (rows.WafTrafficLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpSourceType = "aws_waf_traffic_log"
	row.TpTimestamp = helpers.UnixMillis(*row.Timestamp)
	row.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))

	// Hive fields
	// Check if row.HttpSourceId is not nil before dereferencing
	if row.HttpSourceId != nil {
		row.TpIndex = strings.ReplaceAll(*row.HttpSourceId, "/", `_`)
	} else {
		// Handle the case where HttpSourceId is nil, if needed
		row.TpIndex = "" // or assign a default value if desired
	}
	// convert to date in format yy-mm-dd
	row.TpDate = time.UnixMilli(int64(*row.Timestamp)).Format("2006-01-02")

	return row, nil
}
