package tables

import (
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"time"

	"github.com/rs/xid"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// register the table from the package init function
func init() {
	table.RegisterTable[*rows.CostAndUsageLog, *CostAndUsageLogTable]()
}

// CostAndUsageLogTable - table for CostAndUsageLogs
type CostAndUsageLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[*rows.CostAndUsageLog, *CostAndUsageLogTableConfig, *config.AwsConnection]
}

// Identifier implements table.Table
func (t *CostAndUsageLogTable) Identifier() string {
	return "aws_cost_usage_log"
}

func (t *CostAndUsageLogTable) SupportedSources() []*table.SourceMetadata[*rows.CostAndUsageLog] {
	// TODO fix FileLayout
	defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
		FileLayout: utils.ToStringPointer("/Users/vedmisra/billing-info/(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})"),
	}

	return []*table.SourceMetadata[*rows.CostAndUsageLog]{
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			MapperFunc: mappers.NewCostAndUsageMapper,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig),
				artifact_source.WithRowPerLine(), artifact_source.WithSkipHeaderRow()},
		},
	}
}

// EnrichRow implements table.Table
func (t *CostAndUsageLogTable) EnrichRow(row *rows.CostAndUsageLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.CostAndUsageLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpSourceType = "aws_cost_and_usage_log"

	// we are using the invoice date as the tp_timestamp, but we dont always get the invoice date
	// so we use the BillingPeriodStartDate as a fallback
	// TODO - should we use the BillingPeriodEndDate instead?
	if row.InvoiceDate != nil {
		row.TpTimestamp = *row.InvoiceDate
	} else if row.BillingPeriodStartDate != nil {
		row.TpTimestamp = *row.BillingPeriodStartDate
	} else if row.BillingPeriodEndDate != nil {
		row.TpTimestamp = *row.BillingPeriodEndDate
	}

	row.TpIngestTimestamp = time.Now()
	// if row.PayerAccountName != nil {
	// 	row.TpSourceIP = row.PayerAccountName
	// 	row.TpIps = append(row.TpIps, *row.PayerAccountName)
	// }

	// Hive fields
	// for some rows we dont get the linked account id, so we use the payer account id as a fallback
	// TODO - should we use the payer account id instead?
	if accountId := typehelpers.SafeString(row.LinkedAccountId); accountId != "" {
		row.TpIndex = accountId
	} else {
		row.TpIndex = typehelpers.SafeString(row.PayerAccountId)
	}
	// convert to date in format yy-mm-dd
	if row.InvoiceDate != nil {
		row.TpDate = row.InvoiceDate.Truncate(24 * time.Hour)
	} else if row.BillingPeriodStartDate != nil {
		row.TpDate = row.BillingPeriodStartDate.Truncate(24 * time.Hour)
	} else if row.BillingPeriodEndDate != nil {
		row.TpDate = row.BillingPeriodEndDate.Truncate(24 * time.Hour)
	}

	return row, nil
}
