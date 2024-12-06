package tables

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const SecurityHubFindingTableIdentifier = "aws_securityhub_finding"

// register the table from the package init function
func init() {
	table.RegisterTable[*rows.SecurityHubFinding, *SecurityHubFindingTable]()
}

type SecurityHubFindingTable struct{}

func (c *SecurityHubFindingTable) Identifier() string {
	return SecurityHubFindingTableIdentifier
}

func (c *SecurityHubFindingTable) GetSourceMetadata() []*table.SourceMetadata[*rows.SecurityHubFinding] {
	defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
		FileLayout: utils.ToStringPointer("AWSLogs(?:/o-[a-z0-9]{8,12})?/(?P<account_id>\\d+)/SecurityHub/(?P<region>[a-z0-9-]+)/(?P<year>\\d{4})/(?P<month>\\d{2})/(?P<day>\\d{2})/findings\\.json\\.gz"),
	}

	return []*table.SourceMetadata[*rows.SecurityHubFinding]{
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:  &mappers.SecurityHubFindingMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig),
			},
		},
	}
}

func (c *SecurityHubFindingTable) EnrichRow(row *rows.SecurityHubFinding, sourceEnrichmentFields enrichment.SourceEnrichment) (*rows.SecurityHubFinding, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()

	for _, resource := range row.Resources {
		newAkas := awsAkasFromArn(*resource)
		row.TpAkas = append(row.TpAkas, newAkas...)
	}

	if row.Time != nil {
		row.TpTimestamp = *row.Time
		row.TpDate = row.Time.Truncate(24 * time.Hour)
	}
	if row.Account != nil {
		row.TpIndex = *row.Account
	}

	return row, nil
}
