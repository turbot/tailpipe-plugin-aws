package tables

import (
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const SecurityHubFindingLogTableIdentifier = "aws_security_hub_finding_log"

// register the table from the package init function
func init() {
	table.RegisterTable[*rows.SecurityHubFindingLog, *SecurityHubFindingLogTable]()
}

type SecurityHubFindingLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[*rows.SecurityHubFindingLog, *SecurityHubFindingLogTableConfig, *config.AwsConnection]
}

func (c *SecurityHubFindingLogTable) Identifier() string {
	return SecurityHubFindingLogTableIdentifier
}

func (c *SecurityHubFindingLogTable) SupportedSources() []*table.SourceMetadata[*rows.SecurityHubFindingLog] {
	defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
		FileLayout: utils.ToStringPointer("AWSLogs(?:/o-[a-z0-9]{8,12})?/(?P<account_id>\\d+)/SecurityHub/(?P<region>[a-z0-9-]+)/(?P<year>\\d{4})/(?P<month>\\d{2})/(?P<day>\\d{2})/findings\\.json\\.gz"),
	}

	return []*table.SourceMetadata[*rows.SecurityHubFindingLog]{
		{
			SourceName: constants.ArtifactSourceIdentifier,
			MapperFunc: mappers.NewSecurityHubFindingsMapper,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *SecurityHubFindingLogTable) EnrichRow(row *rows.SecurityHubFindingLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.SecurityHubFindingLog, error) {
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	for _, resource := range row.Resources {
		newAkas := awsAkasFromArn(*resource)
		row.TpAkas = append(row.TpAkas, newAkas...)
	}

	if row.Time != nil {
		row.TpTimestamp = *row.Time
	}
	if row.Account != nil {
		row.TpIndex = *row.Account
	}

	return row, nil
}
