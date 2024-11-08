package tables

import (
	"context"

	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type SecurityHubFindingLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[rows.SecurityHubFindingLog, *SecurityHubFindingLogTableConfig, *config.AwsConnection]
}

func NewSecurityHubFindingLogTable() table.Table {
	return &SecurityHubFindingLogTable{}
}

func (c *SecurityHubFindingLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
	// call base init
	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
		return err
	}

	c.initMapper()
	return nil
}

func (c *SecurityHubFindingLogTable) Identifier() string {
	return "aws_securityhub_finding_log"
}

func (c *SecurityHubFindingLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	var opts []row_source.RowSourceOption

	switch sourceType {
	case artifact_source.AwsS3BucketSourceIdentifier:
		// the default file layout for Cloudtrail logs in S3
		defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
			// TODO #config finalise default cloudtrail file layout
			FileLayout: utils.ToStringPointer("security_hub_findings_(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})(?P<hour>\\d{2})(?P<minute>\\d{2})(?P<second>\\d{2})"),
		}
		opts = append(opts, artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig))
	}

	return opts
}

func (c *SecurityHubFindingLogTable) initMapper() {
	// TODO switch on source
	c.Mapper = mappers.NewSecurityHubFindingsMapper()
}

func (c *SecurityHubFindingLogTable) GetRowSchema() any {
	return &rows.SecurityHubFindingLog{}
}

func (c *SecurityHubFindingLogTable) GetConfigSchema() parse.Config {
	return &SecurityHubFindingLogTableConfig{}
}

// TODO K why does this not fail to compile?
func (c *SecurityHubFindingLogTable) EnrichRow(row rows.SecurityHubFindingLog, sourceEnrichmentFields *enrichment.CommonFields) (rows.SecurityHubFindingLog, error) {
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	for _, resource := range row.Resources {
		newAkas := awsAkasFromArn(*resource)
		row.TpAkas = append(row.TpAkas, newAkas...)
	}

	row.TpSourceType = "aws_securityhub_finding_log"
	if row.Time != nil {
		row.TpTimestamp = *row.Time
	}
	if row.Account != nil {
		row.TpIndex = *row.Account
	}

	return row, nil
}
