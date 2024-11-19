package tables

import (
	"time"

	"github.com/rs/xid"

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

const GuardDutyFindingTableIdentifier = "aws_guardduty_finding"

// register the table from the package init function
func init() {
	table.RegisterTable[*rows.GuardDutyFinding, *GuardDutyFindingTable]()
}

// GuardDutyFindingTable - table for GuardDuty Findings
type GuardDutyFindingTable struct {
	table.TableImpl[*rows.GuardDutyFinding, *GuardDutyFindingTableConfig, *config.AwsConnection]
}

func (c *GuardDutyFindingTable) Identifier() string {
	return GuardDutyFindingTableIdentifier
}

func (c *GuardDutyFindingTable) SupportedSources() []*table.SourceMetadata[*rows.GuardDutyFinding] {
	defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
		FileLayout: utils.ToStringPointer("AWSLogs(?:/o-[a-z0-9]{8,12})?/[0-9]+/GuardDuty/[a-z0-9-]+/(?P<year>\\d{4})/(?P<month>\\d{2})/(?P<day>\\d{2})/[0-9a-fA-F-]+\\.jsonl\\.gz"),
	}

	return []*table.SourceMetadata[*rows.GuardDutyFinding]{
		{
			SourceName: constants.ArtifactSourceIdentifier,
			MapperFunc: mappers.NewGuardDutyMapper,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *GuardDutyFindingTable) EnrichRow(row *rows.GuardDutyFinding, sourceEnrichmentFields *enrichment.CommonFields) (*rows.GuardDutyFinding, error) {
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}
	row.TpID = xid.New().String()
	row.TpTimestamp = row.CreatedAt
	row.TpIngestTimestamp = time.Now()
	row.TpDate = row.CreatedAt.Truncate(24 * time.Hour)
	row.TpIndex = *row.AccountId

	if row.IpAddressV4 != nil {
		row.TpIps = append(row.TpIps, *row.IpAddressV4)
	}
	if row.IpAddressV6 != nil {
		row.TpIps = append(row.TpIps, *row.IpAddressV6)
	}
	if row.Ipv6Addresses != nil {
		row.TpIps = append(row.TpIps, row.Ipv6Addresses...)
	}

	if row.UserName != nil {
		row.TpUsernames = append(row.TpUsernames, *row.UserName)
	}

	return row, nil
}
