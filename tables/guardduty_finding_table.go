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

const GuardDutyFindingTableIdentifier = "aws_guardduty_finding"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.GuardDutyFinding, *GuardDutyFindingTable]()
}

// GuardDutyFindingTable - table for GuardDuty Findings
type GuardDutyFindingTable struct{}

func (c *GuardDutyFindingTable) Identifier() string {
	return GuardDutyFindingTableIdentifier
}

func (c *GuardDutyFindingTable) GetSourceMetadata() []*table.SourceMetadata[*rows.GuardDutyFinding] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/%{NUMBER:account_id}/GuardDuty/%{DATA:region_path}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA:unique_id}.jsonl.gz"),
	}

	return []*table.SourceMetadata[*rows.GuardDutyFinding]{
		{
			// S3 artifact source
			SourceName: sources.AwsS3BucketSourceIdentifier,
			Mapper:     &mappers.GuardDutyMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &mappers.GuardDutyMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *GuardDutyFindingTable) EnrichRow(row *rows.GuardDutyFinding, sourceEnrichmentFields schema.SourceEnrichment) (*rows.GuardDutyFinding, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields
	row.TpID = xid.New().String()
	row.TpTimestamp = row.CreatedAt
	row.TpIngestTimestamp = time.Now()
	row.TpDate = row.CreatedAt.Truncate(24 * time.Hour)
	row.TpIndex = *row.AccountId

	row.TpAkas = append(row.TpAkas, *row.Arn)

	if row.Resource != nil && row.Resource.AccessKeyDetails != nil {
		// usernames
		if row.Resource.AccessKeyDetails.AccessKeyId != nil {
			row.TpUsernames = append(row.TpUsernames, *row.Resource.AccessKeyDetails.AccessKeyId)
		}
		if row.Resource.AccessKeyDetails.UserName != nil {
			row.TpUsernames = append(row.TpUsernames, *row.Resource.AccessKeyDetails.UserName)
		}
		if row.Resource.AccessKeyDetails.PrincipalId != nil {
			row.TpUsernames = append(row.TpUsernames, *row.Resource.AccessKeyDetails.PrincipalId)
		}

	}

	return row, nil
}
