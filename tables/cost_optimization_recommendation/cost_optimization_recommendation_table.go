package cost_optimization_recommendation

import (
	"time"

	"github.com/rs/xid"
	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const CostOptimizationRecommendationsTableIdentifier = "aws_cost_optimization_recommendation"

// CostOptimizationRecommendationsTable - table for CostOptimizationRecommendations
type CostOptimizationRecommendationsTable struct{}

// Identifier implements table.Table
func (t *CostOptimizationRecommendationsTable) Identifier() string {
	return CostOptimizationRecommendationsTableIdentifier
}

// GetSourceMetadata implements table.Table
func (t *CostOptimizationRecommendationsTable) GetSourceMetadata() []*table.SourceMetadata[*CostOptimizationRecommendation] {
	defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("%{DATA:export_name}/(?:data/%{DATA:partition}/)?(?:%{INT:from_date}-%{INT:to_date}/)?(?:%{DATA:assembly_id}/)?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.(?:zip|gz)"),
	}

	return []*table.SourceMetadata[*CostOptimizationRecommendation]{
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig),
				artifact_source.WithArtifactExtractor(NewCostOptimizationRecommendationExtractor()),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewCostOptimizationRecommendationExtractor()),
			},
		},
	}
}

// EnrichRow implements table.Table
func (t *CostOptimizationRecommendationsTable) EnrichRow(row *CostOptimizationRecommendation, sourceEnrichmentFields schema.SourceEnrichment) (*CostOptimizationRecommendation, error) {
	// initialize the enrichment fields to any fields provided by the source
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	
	row.TpTimestamp = *row.LastRefreshTimestamp
	// convert to date in format yyyy-mm-dd
	row.TpDate = row.LastRefreshTimestamp.Truncate(24 * time.Hour)

	// TpIndex
	if typehelpers.SafeString(row.AccountID) != "" {
		row.TpIndex = typehelpers.SafeString(row.AccountID)
	} else {
		row.TpIndex = schema.DefaultIndex
	}

	if row.ResourceARN != nil {
		row.TpAkas = append(row.TpAkas, *row.ResourceARN)
	}

	return row, nil
}

func (c *CostOptimizationRecommendationsTable) GetDescription() string {
	return "AWS Cost Optimization Recommendations provide insights into opportunities to reduce AWS spending through various actions such as rightsizing, reserved instances, savings plans, and idle resource cleanup."
}
