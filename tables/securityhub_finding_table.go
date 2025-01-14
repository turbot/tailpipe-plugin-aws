package tables

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-aws/sources"
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
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/%{NUMBER:account_id}/SecurityHub/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/findings.json.gz"),
	}

	return []*table.SourceMetadata[*rows.SecurityHubFinding]{
		{
			// S3 artifact source
			SourceName: sources.AwsS3BucketSourceIdentifier,
			Mapper:     &mappers.SecurityHubFindingMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &mappers.SecurityHubFindingMapper{},
		},
	}
}

func (c *SecurityHubFindingTable) EnrichRow(row *rows.SecurityHubFinding, sourceEnrichmentFields schema.SourceEnrichment) (*rows.SecurityHubFinding, error) {
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
