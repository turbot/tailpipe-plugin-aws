package tables

import (
	"context"
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

// WafTrafficLogTable - table for CloudTrailLog logs
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
		// the default file layout for Cloudtrail logs in S3
		defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
			// TODO #config finalise default cloudtrail file layout
			FileLayout: utils.ToStringPointer("AWSLogs(?:/o-[a-z0-9]{8,12})?/\\d+/CloudTrail/[a-z-0-9]+/\\d{4}/\\d{2}/\\d{2}/(?P<index>\\d+)_CloudTrail_(?P<region>[a-z-0-9]+)_(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})T(?P<hour>\\d{2})(?P<minute>\\d{2})Z_.+.json.gz"),
		}
		opts = append(opts, artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig))
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
	row.TpTimestamp = helpers.UnixMillis(row.Timestamp.UnixNano() / int64(time.Millisecond))
	row.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))

	// Hive fields
	row.TpPartition = "aws_waf_traffic_log"
	// row.TpIndex = row.RecipientAccountId
	// convert to date in format yy-mm-dd
	row.TpDate = time.UnixMilli(int64(row.Timestamp.UnixNano())).Format("2006-01-02")

	return row, nil
}
