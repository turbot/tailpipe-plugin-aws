package cost_and_usage_report

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

const CostUsageReportTableIdentifier = "aws_cost_and_usage_report"

// CostUsageReportTable - table for CostUsageReports
type CostUsageReportTable struct{}

// Identifier implements table.Table
func (t *CostUsageReportTable) Identifier() string {
	return CostUsageReportTableIdentifier
}

func (t *CostUsageReportTable) GetSourceMetadata() []*table.SourceMetadata[*CostUsageReport] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("%{DATA:prefix}/%{DATA:exportName}/%{DATA:folderName}/%{DATA:timestamp}/%{DATA}.csv.(?:gz|zip)"),
	}

	return []*table.SourceMetadata[*CostUsageReport]{
		{
			// any artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewCostUsageReportExtractor())},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewCostUsageReportExtractor())},
		},
	}
}

// EnrichRow implements table.Table
func (t *CostUsageReportTable) EnrichRow(row *CostUsageReport, sourceEnrichmentFields schema.SourceEnrichment) (*CostUsageReport, error) {
	// initialize the enrichment fields to any fields provided by the source
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()

	if row.LineItemUsageStartDate != nil {
		row.TpTimestamp = *row.LineItemUsageStartDate

		// convert to date in format yyyy-mm-dd
		row.TpDate = row.LineItemUsageStartDate.Truncate(24 * time.Hour)
	} else if row.LineItemUsageEndDate != nil {
		row.TpTimestamp = *row.LineItemUsageEndDate

		// convert to date in format yyyy-mm-dd
		row.TpDate = row.LineItemUsageEndDate.Truncate(24 * time.Hour)
	} else if row.BillBillingPeriodStartDate != nil {
		row.TpTimestamp = *row.BillBillingPeriodStartDate

		// convert to date in format yyyy-mm-dd
		row.TpDate = row.BillBillingPeriodStartDate.Truncate(24 * time.Hour)
	} else if row.BillBillingPeriodEndDate != nil {
		row.TpTimestamp = *row.BillBillingPeriodEndDate

		// convert to date in format yyyy-mm-dd
		row.TpDate = row.BillBillingPeriodEndDate.Truncate(24 * time.Hour)
	}

	// TpIndex
	switch {
	case typehelpers.SafeString(row.LineItemUsageAccountId) != "":
		row.TpIndex = typehelpers.SafeString(row.LineItemUsageAccountId)
	case typehelpers.SafeString(row.BillPayerAccountId) != "":
		row.TpIndex = typehelpers.SafeString(row.BillPayerAccountId)
	default:
		row.TpIndex = schema.DefaultIndex
	}

	return row, nil
}

func (c *CostUsageReportTable) GetDescription() string {
	return "AWS Cost and Usage Reports (CUR) provide a comprehensive breakdown of AWS service costs and usage. This table offers a structured view of billing data, including service charges, account-level spending, resource consumption, discounts, and pricing details. It enables cost analysis, budget tracking, and optimization insights across AWS accounts."
}