package cost_and_usage_focus_1_0

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

const Focus1_0TableIdentifier = "aws_cost_and_usage_focus_1_0"

type Focus1_0Table struct{}

func (c *Focus1_0Table) Identifier() string {
	return Focus1_0TableIdentifier
}

func (c *Focus1_0Table) GetSourceMetadata() []*table.SourceMetadata[*Focus1_0] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("%{DATA:prefix}/%{DATA:exportName}/%{DATA:folderName}/%{DATA:billing_period}/%{DATA:assembly_id}/%{DATA}.csv.(?:gz|zip)"),
	}

	return []*table.SourceMetadata[*Focus1_0]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewCostUsageFocus_1_0_Extractor()),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewCostUsageFocus_1_0_Extractor()),
			},
		},
	}
}

func (c *Focus1_0Table) EnrichRow(row *Focus1_0, sourceEnrichmentFields schema.SourceEnrichment) (*Focus1_0, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()

	// TpIndex
	switch {
	case typehelpers.SafeString(row.SubAccountId) != "":
		row.TpIndex = typehelpers.SafeString(row.SubAccountId)
	case typehelpers.SafeString(row.BillingAccountId) != "":
		row.TpIndex = typehelpers.SafeString(row.BillingAccountId)
	default:
		row.TpIndex = schema.DefaultIndex
	}

	if row.ChargePeriodStart != nil {
		row.TpTimestamp = *row.ChargePeriodStart

		row.TpDate = row.ChargePeriodStart.Truncate(24 * time.Hour)
	} else if row.ChargePeriodEnd != nil {
		row.TpTimestamp = *row.ChargePeriodEnd

		row.TpDate = row.ChargePeriodEnd.Truncate(24 * time.Hour)
	} else if row.BillingPeriodStart != nil {
		row.TpTimestamp = *row.BillingPeriodStart

		row.TpDate = row.BillingPeriodStart.Truncate(24 * time.Hour)
	} else if row.BillingPeriodEnd != nil {
		row.TpTimestamp = *row.BillingPeriodEnd

		row.TpDate = row.BillingPeriodEnd.Truncate(24 * time.Hour)
	}

	if row.ResourceId != nil && strings.HasPrefix(*row.ResourceId, "arn:") {
		row.TpAkas = append(row.TpAkas, *row.ResourceId)
	}

	return row, nil
}

func (c *Focus1_0Table) GetDescription() string {
	return "AWS FOCUS 1.0 (Flexible, Optimized, and Comprehensive Usage and Savings) provides a detailed breakdown of AWS service usage and cost optimization opportunities. This table structures billing and usage data, including pricing details, commitment-based discounts, capacity reservations, and SKU-level pricing metrics. It enables cost tracking, commitment analysis, and efficient cloud financial management."
}
