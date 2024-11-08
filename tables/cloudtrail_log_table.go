package tables

import (
	"context"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

// CloudTrailLogTable - table for CloudTrailLog logs
type CloudTrailLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[rows.CloudTrailLog, *CloudTrailLogTableConfig, *config.AwsConnection]
}

func NewCloudTrailLogTable() table.Table {
	return &CloudTrailLogTable{}
}

func (c *CloudTrailLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
	// call base init
	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
		return err
	}

	c.initMapper()
	return nil
}

func (c *CloudTrailLogTable) initMapper() {
	// TODO switch on source

	// if the source is an artifact source, we need a mapper
	c.Mapper = mappers.NewCloudtrailMapper()
}

// Identifier implements table.Table
func (c *CloudTrailLogTable) Identifier() string {
	return "aws_cloudtrail_log"
}

// GetSourceOptions returns any options which should be passed to the given source type
func (c *CloudTrailLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	var opts []row_source.RowSourceOption

	switch sourceType {
	case artifact_source.AwsS3BucketSourceIdentifier:
		// the default file layout for Cloudtrail logs in S3
		defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
			// TODO #config finalise default cloudtrail file layout
			FileLayout: utils.ToStringPointer("AWSLogs(?:/o-[a-z0-9]{8,12})?/\\d+/CloudTrail/[a-z-0-9]+/\\d{4}/\\d{2}/\\d{2}/(?P<index>\\d+)_CloudTrail_(?P<region>[a-z-0-9]+)_(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})T(?P<hour>\\d{2})(?P<minute>\\d{2})Z_.+.json.gz"),
		}
		opts = append(opts, artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig))
	}

	return opts
}

// GetRowSchema implements table.Table
func (c *CloudTrailLogTable) GetRowSchema() any {
	return rows.CloudTrailLog{}
}

func (c *CloudTrailLogTable) GetConfigSchema() parse.Config {
	return &CloudTrailLogTableConfig{}
}

// EnrichRow implements table.Table
func (c *CloudTrailLogTable) EnrichRow(row rows.CloudTrailLog, sourceEnrichmentFields *enrichment.CommonFields) (rows.CloudTrailLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpSourceType = "aws_cloudtrail_log"
	row.TpTimestamp = time.Unix(0, int64(row.EventTime)*int64(time.Millisecond))
	row.TpIngestTimestamp = time.Unix(0, int64(helpers.UnixMillis(time.Now().UnixNano()/int64(time.Millisecond)))*int64(time.Millisecond))
	if row.SourceIPAddress != nil {
		row.TpSourceIP = row.SourceIPAddress
		row.TpIps = append(row.TpIps, *row.SourceIPAddress)
	}
	for _, resource := range row.Resources {
		if resource.ARN != nil {
			newAkas := awsAkasFromArn(*resource.ARN)
			row.TpAkas = append(row.TpAkas, newAkas...)
		}
	}
	// If it's an AKIA, then record that as an identity. Do not record ASIA*
	// keys etc.
	if row.UserIdentity.AccessKeyId != nil {
		if strings.HasPrefix(*row.UserIdentity.AccessKeyId, "AKIA") {
			row.TpUsernames = append(row.TpUsernames, *row.UserIdentity.AccessKeyId)
		}
	}
	if row.UserIdentity.UserName != nil {
		row.TpUsernames = append(row.TpUsernames, *row.UserIdentity.UserName)
	}

	// Hive fields
	row.TpPartition = "default" // TODO - should be based on the definition in HCL
	row.TpIndex = row.RecipientAccountId
	// convert to date in format yy-mm-dd
	row.TpDate = time.UnixMilli(int64(row.EventTime)).Format("2006-01-02")

	return row, nil
}
