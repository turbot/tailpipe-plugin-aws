package detailed_billing_report

import (
	"log/slog"
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

const DetailedBillingReportTableIdentifier = "aws_detailed_billing_report"

type DetailedBillingReportTable struct{}

func (t *DetailedBillingReportTable) Identifier() string {
	return DetailedBillingReportTableIdentifier
}

func (t *DetailedBillingReportTable) GetSourceMetadata() []*table.SourceMetadata[*DetailedBillingReport] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("%{DATA:prefix}/%{DATA:exportName}/%{DATA:folderName}/%{DATA:billing_period}/%{DATA:assembly_id}/%{DATA}.csv.(?:gz|zip)"),
	}

	return []*table.SourceMetadata[*DetailedBillingReport]{
		{
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewDetailedBillingReportExtractor()),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewDetailedBillingReportExtractor()),
			},
		},
	}
}

func (t *DetailedBillingReportTable) EnrichRow(row *DetailedBillingReport, sourceEnrichmentFields schema.SourceEnrichment) (*DetailedBillingReport, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()

	// Fallbacks for TpIndex and TpTimestamp
	if typehelpers.SafeString(row.LinkedAccountId) != "" {
		row.TpIndex = *row.LinkedAccountId
	} else if typehelpers.SafeString(row.PayerAccountId) != "" {
		row.TpIndex = *row.PayerAccountId
	} else {
		row.TpIndex = schema.DefaultIndex
	}

	if row.UsageStartDate != nil {
		row.TpTimestamp = *row.UsageStartDate
		row.TpDate = row.UsageStartDate.Truncate(24 * time.Hour)
	} else if row.UsageEndDate != nil {
		row.TpTimestamp = *row.UsageEndDate
		row.TpDate = row.UsageEndDate.Truncate(24 * time.Hour)
	} else if row.BillingPeriodStart != nil {
		row.TpTimestamp = *row.BillingPeriodStart
		row.TpDate = row.BillingPeriodStart.Truncate(24 * time.Hour)
	} else if row.BillingPeriodEnd != nil {
		row.TpTimestamp = *row.BillingPeriodEnd
		row.TpDate = row.BillingPeriodEnd.Truncate(24 * time.Hour)
	} else {
		slog.Error("Timestamp is not available for the record ID:", row.RecordID)
	}

	if row.ResourceId != nil && strings.HasPrefix(*row.ResourceId, "arn:") {
		row.TpAkas = append(row.TpAkas, *row.ResourceId)
	}

	return row, nil
}

func (t *DetailedBillingReportTable) GetDescription() string {
	return "Detailed AWS billing report that includes cost breakdowns, usage quantities, tax details, and credits per product, operation, and account across various date ranges."
}
