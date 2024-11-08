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
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

// CostAndUsageLogTable - table for CostAndUsageLogs
type CostAndUsageLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[rows.CostAndUsageLog, *CostAndUsageLogTableConfig, *config.AwsConnection]
}

func NewCCostAndUsageLogTable() table.Table {
	return &CostAndUsageLogTable{}
}

func (c *CostAndUsageLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
	// call base init
	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
		return err
	}

	c.initMapper()
	return nil
}

func (c *CostAndUsageLogTable) initMapper() {
	// TODO switch on source

	// if the source is an artifact source, we need a mapper
	c.Mapper = mappers.NewCostAndUsageMapper()
}

// Identifier implements table.Table
func (c *CostAndUsageLogTable) Identifier() string {
	return "aws_cost_and_usage_log"
}

// GetSourceOptions returns any options which should be passed to the given source type
func (c *CostAndUsageLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	var opts []row_source.RowSourceOption

	switch sourceType {

	// TODO - update to use AwsS3BucketSourceIdentifier (using FileSystemSourceIdentifier for now)
	case artifact_source.FileSystemSourceIdentifier:
		defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
			FileLayout: utils.ToStringPointer("/Users/pskrbasu/s3-logs/(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})/ecsv_10_2024.csv"),
		}
		opts = append(opts, artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig), artifact_source.WithRowPerLine(), artifact_source.WithSkipHeaderRow())

	}

	return opts
}

// GetRowSchema implements table.Table
func (c *CostAndUsageLogTable) GetRowSchema() any {
	return rows.CostAndUsageLog{}
}

func (c *CostAndUsageLogTable) GetConfigSchema() parse.Config {
	return &CostAndUsageLogTableConfig{}
}

// EnrichRow implements table.Table
func (c *CostAndUsageLogTable) EnrichRow(row rows.CostAndUsageLog, sourceEnrichmentFields *enrichment.CommonFields) (rows.CostAndUsageLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpSourceType = "aws_cost_and_usage_log"
	row.TpTimestamp = *row.InvoiceDate
	row.TpIngestTimestamp = time.Now()
	// if row.PayerAccountName != nil {
	// 	row.TpSourceIP = row.PayerAccountName
	// 	row.TpIps = append(row.TpIps, *row.PayerAccountName)
	// }

	// Hive fields
	row.TpPartition = "default" // TODO - should be based on the definition in HCL
	row.TpIndex = *row.PayerAccountId
	// convert to date in format yy-mm-dd
	row.TpDate = row.InvoiceDate.Format("2006-01-02")

	return row, nil
}
