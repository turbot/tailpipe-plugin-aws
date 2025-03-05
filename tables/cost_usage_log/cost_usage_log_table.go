package cost_usage_log

import (
	"time"

	"github.com/rs/xid"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const CostUsageLogTableIdentifier = "aws_cost_usage_log"

// CostAndUsageLogTable - table for CostAndUsageLogs
type CostAndUsageLogTable struct{}

// Identifier implements table.Table
func (t *CostAndUsageLogTable) Identifier() string {
	return CostUsageLogTableIdentifier
}

func (t *CostAndUsageLogTable) GetSourceMetadata() []*table.SourceMetadata[*CostAndUsageLog] {

	// TODO: Cross-check https://docs.aws.amazon.com/cur/latest/userguide/understanding-report-versions.html
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("%{DATA:prefix}/%{DATA:exportName}/%{DATA:data}/%{DATA:timestampz}/%{DATA}.csv.zip"),

		// 		s3://cost-usage-report-log/report/test52/20250201-20250301/20250228T214620Z/test52-00001.csv.zip
		//report/
		// test52/
		// 20250201-20250301/
		// 20250228T214620Z/
		// test52-00001.csv.zip
	}

	return []*table.SourceMetadata[*CostAndUsageLog]{
		{
			// any artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			// Mapper:     &CostAndUsageLogMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewCostUsageLogExtractor())},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			// Mapper:     &CostAndUsageLogMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewCostUsageLogExtractor())},
		},
	}
}

// EnrichRow implements table.Table
func (t *CostAndUsageLogTable) EnrichRow(row *CostAndUsageLog, sourceEnrichmentFields schema.SourceEnrichment) (*CostAndUsageLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()

	// we are using the invoice date as the tp_timestamp, but we dont always get the invoice date
	// so we use the BillingPeriodStartDate as a fallback
	// TODO - should we use the BillingPeriodEndDate instead?
	if row.BillBillingPeriodStartDate != nil {
		row.TpTimestamp = *row.BillBillingPeriodStartDate

		// convert to date in format yyyy-mm-dd
		row.TpDate = row.BillBillingPeriodStartDate.Truncate(24 * time.Hour)
	} else if row.BillBillingPeriodEndDate != nil {
		row.TpTimestamp = *row.BillBillingPeriodEndDate

		// convert to date in format yyyy-mm-dd
		row.TpDate = row.BillBillingPeriodEndDate.Truncate(24 * time.Hour)
	}
	

	// if row.PayerAccountName != nil {
	// 	row.TpSourceIP = row.PayerAccountName
	// 	row.TpIps = append(row.TpIps, *row.PayerAccountName)
	// }

	// Hive fields
	// for some rows we dont get the linked account id, so we use the payer account id as a fallback
	// TODO - should we use the payer account id instead?
	switch {
	case typehelpers.SafeString(row.BillPayerAccountId) != "":
		row.TpIndex = typehelpers.SafeString(row.BillPayerAccountId)
	case typehelpers.SafeString(row.LineItemUsageAccountId) != "":
		row.TpIndex = typehelpers.SafeString(row.LineItemUsageAccountId)
	default:
		row.TpIndex = schema.DefaultIndex
	}

	return row, nil
}
