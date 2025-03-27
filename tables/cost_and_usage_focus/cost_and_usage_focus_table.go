package cost_and_usage_focus

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

const CostUsageFocusTableIdentifier = "aws_cost_and_usage_focus"

type CostUsageFocusTable struct{}

func (c *CostUsageFocusTable) Identifier() string {
	return CostUsageFocusTableIdentifier
}

func (c *CostUsageFocusTable) GetSourceMetadata() []*table.SourceMetadata[*CostUsageFocus] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("%{DATA:prefix}/%{DATA:export_name}/%{DATA:folder_name}/%{DATA:billing_period}/%{DATA:assembly_id}/%{DATA}.csv.(?:gz|zip)"),
	}

	return []*table.SourceMetadata[*CostUsageFocus]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewCostUsageFocusExtractor()),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewCostUsageFocusExtractor()),
			},
		},
	}
}

func (c *CostUsageFocusTable) EnrichRow(row *CostUsageFocus, sourceEnrichmentFields schema.SourceEnrichment) (*CostUsageFocus, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()


	// TpIndex
	switch {
	case typehelpers.SafeString(row.SubAccountId) != "":
		row.TpIndex = typehelpers.SafeString(row.SubAccountId)
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

func (c *CostUsageFocusTable) GetDescription() string {
	return "AWS FOCUS 1.0 (Flexible, Optimized, and Comprehensive Usage and Savings) provides a detailed breakdown of AWS service usage and cost optimization opportunities. This table structures billing and usage data, including pricing details, commitment-based discounts, capacity reservations, and SKU-level pricing metrics. It enables cost tracking, commitment analysis, and efficient cloud financial management."
}
