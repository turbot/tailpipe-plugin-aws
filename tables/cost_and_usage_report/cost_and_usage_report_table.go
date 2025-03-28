package cost_and_usage_report

import (
	"strings"
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
		// Grok pattern to match AWS CUR legacy, and CUR 2.0 report file paths in Amazon S3.
		//
		// Pattern:
		// %{DATA:export_name}/(?:data/%{DATA:partition}/)?(?:%{INT:from_date}-%{INT:to_date}/)?(?:%{DATA:assembly_id}/)?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.(?:zip|gz)
		//
		// This pattern supports both "Create new" and "Overwrite" file naming conventions used by AWS CUR reports:
		// - "Create new" layout: <export-name>/data/<partition>/<timestamp>-<execution-id>/<export-name>-<chunk-number>.csv.gz
		// - "Overwrite" layout: <export-name>/data/<partition>/<export-name>-<chunk-number>.csv.gz
		//
		// Additionally, it supports the legacy CUR layout:
		// - Legacy layout: <export-name>/<from_date>-<to_date>/<timestamp>/<export-name>-<chunk-number>.csv.zip
		//
		// Notes:
		// - `partition` captures CUR 2.0 and FOCUS 1.0 partition values (e.g., BILLING_PERIOD=YYYY-MM) and cost optimization format (e.g., date=YYYY-MM-DD).
		// - `from_date` and `to_date` are used in legacy CUR exports (e.g., 20250301-20250401).
		// - `assembly_id` and `execution_id` are optional identifiers that vary by report version.
		// - `timestamp` is a string like 20250307T053621Z.
		// - The pattern deliberately omits the S3 prefix, as it is handled by the `prefix` argument in `aws_s3_bucket` sources.
		//
		// Example S3 keys matched:
		// - cur-2-0-daily-csv/data/BILLING_PERIOD=2025-03/cur-2-0-daily-csv-00003.csv.gz
		// - report-name/20250101-20250201/assembly123/report-name-00001.csv.zip
		// - report-name/20250101-20250201/report-name-00001.csv.gz
		// - export/data/PARTITION1/20250307T053621Z-exec123/export-00001.csv.zip
		// - export/data/PARTITION1/export-00002.csv.gz
		// - cost-usage-legacy-export/20250301-20250401/20250307T053621Z/cost-usage-legacy-export-00001.csv.zip
		//
		// References:
		// - AWS Export Delivery Formats CUR 2.0: https://docs.aws.amazon.com/cur/latest/userguide/dataexports-export-delivery.html
		// - AWS Report Versioning CUR legacy: https://docs.aws.amazon.com/cur/latest/userguide/understanding-report-versions.html
		FileLayout: utils.ToStringPointer("%{DATA:export_name}/(?:data/%{DATA:partition}/)?(?:%{INT:from_date}-%{INT:to_date}/)?(?:%{DATA:assembly_id}/)?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.(?:zip|gz)"),
	}

	return []*table.SourceMetadata[*CostUsageReport]{
		{
			// any artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewCostUsageReportExtractor()),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewCostUsageReportExtractor()),
			},
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
	// Set TpIndex for the row to help uniquely identify the resource owner or origin account.
	// Priority:
	//   1. Use LineItemUsageAccountId if available (typically the owning AWS account).
	//   2. If not, attempt to extract the 5th element from LineItemResourceId ARN (e.g., the account ID).
	//   3. If both are missing, fall back to a default value.
	//
	// Example:
	//   For LineItemResourceId: "arn:aws:ec2:us-east-1:123456789012:volume/vol-0abcd1234efgh5678"
	//   → TpIndex will be set to "123456789012".
	if row.LineItemUsageAccountId != nil {
		row.TpIndex = typehelpers.SafeString(row.LineItemUsageAccountId)
	} else if row.LineItemResourceId != nil {
		splitResourceId := strings.Split(*row.LineItemResourceId, ":")
		if len(splitResourceId) > 4 {
			row.TpIndex = splitResourceId[4]
		}
	} else {
		row.TpIndex = schema.DefaultIndex
	}

	return row, nil
}

func (c *CostUsageReportTable) GetDescription() string {
	return "AWS Cost and Usage Reports (CUR) provide a comprehensive breakdown of AWS service costs and usage. This table offers a structured view of billing data, including service charges, account-level spending, resource consumption, discounts, and pricing details. It enables cost analysis, budget tracking, and optimization insights across AWS accounts."
}
