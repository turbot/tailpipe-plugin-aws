package guardduty_finding

import (
	"time"

	"github.com/rs/xid"

	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const GuardDutyFindingTableIdentifier = "aws_guardduty_finding"

// GuardDutyFindingTable - table for GuardDuty Findings
type GuardDutyFindingTable struct{}

func (c *GuardDutyFindingTable) Identifier() string {
	return GuardDutyFindingTableIdentifier
}

func (c *GuardDutyFindingTable) GetSourceMetadata() ([]*table.SourceMetadata[*GuardDutyFinding], error) {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/GuardDuty/%{DATA:region_path}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.jsonl.gz"),
	}

	return []*table.SourceMetadata[*GuardDutyFinding]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     &GuardDutyMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &GuardDutyMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}, nil
}

func (c *GuardDutyFindingTable) EnrichRow(row *GuardDutyFinding, sourceEnrichmentFields schema.SourceEnrichment) (*GuardDutyFinding, error) {
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

func (c *GuardDutyFindingTable) GetDescription() string {
	return "AWS GuardDuty findings provide detailed security alerts about potential threats and suspicious activities detected in your AWS environment. This table captures comprehensive information about each finding, including threat details, affected resources, and severity levels to help security teams identify and respond to potential security issues."
}
