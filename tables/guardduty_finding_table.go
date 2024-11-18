package tables

//import (
//	"github.com/rs/xid"
//	"github.com/turbot/pipe-fittings/utils"
//	"github.com/turbot/tailpipe-plugin-aws/config"
//	"github.com/turbot/tailpipe-plugin-aws/mappers"
//	"github.com/turbot/tailpipe-plugin-aws/rows"
//	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
//	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
//	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
//	"github.com/turbot/tailpipe-plugin-sdk/parse"
//	"github.com/turbot/tailpipe-plugin-sdk/row_source"
//	"github.com/turbot/tailpipe-plugin-sdk/table"
//	"github.com/turbot/tailpipe-plugin-sdk/types"
//	"time"
//)
//
//// register the table from the package init function
//func init() {
//	table.RegisterTable(NewGuardDutyFindingTable)
//}
//
//// GuardDutyFindingTable - table for GuardDuty Findings
//type GuardDutyFindingTable struct {
//	table.TableImpl[*rows.GuardDutyFinding, *GuardDutyFindingTableConfig, *config.AwsConnection]
//}
//
//func NewGuardDutyFindingTable() table.Table[*rows.GuardDutyFinding] {
//	return &GuardDutyFindingTable{}
//}
//
//func (c *GuardDutyFindingTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
//	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
//		return err
//	}
//	c.initMapper()
//	return nil
//}
//
//func (c *GuardDutyFindingTable) initMapper() {
//	c.Mapper = mappers.NewGuardDutyMapper()
//}
//
//func (c *GuardDutyFindingTable) Identifier() string {
//	return "aws_guardduty_finding"
//}
//
//func (c *GuardDutyFindingTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
//	var opts []row_source.RowSourceOption
//	switch sourceType {
//	case artifact_source.AwsS3BucketSourceIdentifier:
//		defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
//			FileLayout: utils.ToStringPointer("AWSLogs/[0-9]+/GuardDuty/[a-z0-9-]+/(?P<year>\\d{4})/(?P<month>\\d{2})/(?P<day>\\d{2})/[0-9a-fA-F-]+\\.jsonl\\.gz"),
//		}
//		opts = append(opts, artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig), artifact_source.WithRowPerLine())
//	}
//	return opts
//}
//
//func (c *GuardDutyFindingTable) GetRowSchema() types.RowStruct {
//	return rows.GuardDutyFinding{}
//}
//
//func (c *GuardDutyFindingTable) GetConfigSchema() parse.Config {
//	return &GuardDutyFindingTableConfig{}
//}
//
//func (c *GuardDutyFindingTable) EnrichRow(row *rows.GuardDutyFinding, sourceEnrichmentFields *enrichment.CommonFields) (*rows.GuardDutyFinding, error) {
//	if sourceEnrichmentFields != nil {
//		row.CommonFields = *sourceEnrichmentFields
//	}
//	row.TpID = xid.New().String()
//	row.TpSourceType = "aws_guardduty_finding"
//	row.TpTimestamp = row.CreatedAt
//	row.TpIngestTimestamp = time.Now()
//
//	if row.IpAddressV4 != nil {
//		if row.IpAddressV4 != nil {
//			row.TpIps = append(row.TpIps, *row.IpAddressV4)
//		}
//	}
//
//	row.TpIndex = *row.AccountId
//	row.TpDate = row.CreatedAt.Truncate(24 * time.Hour)
//
//	return row, nil
//}
