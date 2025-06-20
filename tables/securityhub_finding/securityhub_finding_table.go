package securityhub_finding

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-aws/tables"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const SecurityHubFindingTableIdentifier = "aws_securityhub_finding"

type SecurityHubFindingTable struct{}

func (c *SecurityHubFindingTable) Identifier() string {
	return SecurityHubFindingTableIdentifier
}

func (c *SecurityHubFindingTable) GetSourceMetadata() ([]*table.SourceMetadata[*SecurityHubFinding], error) {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		// AWS does not provide a default or recommended S3 bucket file path for Security Hub Findings export.
		// The default file layout defined here is based on the implementation from:
		// https://github.com/aws-samples/aws-security-hub-findings-export
		FileLayout: utils.ToStringPointer("AWSLogs/%{NUMBER:account_id}/%{DATA:security_hub_integrrated_product_name}/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json"),
	}

	return []*table.SourceMetadata[*SecurityHubFinding]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     &SecurityHubFindingMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewSecurityHubFindingExtractor()),
			},
		},
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &SecurityHubFindingMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithArtifactExtractor(NewSecurityHubFindingExtractor()),
			},
		},
	}, nil
}

func (c *SecurityHubFindingTable) EnrichRow(row *SecurityHubFinding, sourceEnrichmentFields schema.SourceEnrichment) (*SecurityHubFinding, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()

	if row.ProductArn != nil {
		akas := tables.AwsAkasFromArn(*row.ProductArn)
		row.TpAkas = append(row.TpAkas, akas...)
	}

	if row.Time != nil {
		row.TpTimestamp = *row.Time
		row.TpDate = row.Time.Truncate(24 * time.Hour)
	}

	return row, nil
}

func (c *SecurityHubFindingTable) GetDescription() string {
	return "AWS Security Hub findings provide detailed information about potential security issues and compliance violations detected across your AWS accounts and resources. This table captures comprehensive security findings from various AWS security services and partner integrations, including details about the affected resources, severity levels, compliance status, and recommended remediation steps."
}
