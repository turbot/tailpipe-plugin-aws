package tables

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/util"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// CloudTrailLogTable - table for CloudTrail logs
type CloudTrailLogTable struct {
	// all tables must embed table.TableBase
	table.TableBase[*CloudTrailLogTableConfig]
}

func NewCloudTrailLogTable() table.Table {
	return &CloudTrailLogTable{}
}

// Identifier implements table.Table
func (c *CloudTrailLogTable) Identifier() string {
	return "aws_cloudtrail_log"
}

// GetSourceOptions returns any options which should be passed to the given source type
func (c *CloudTrailLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	var opts = []row_source.RowSourceOption{
		// if the source is an artifact source, we need a mapper
		// NOTE: WithArtifactMapper option will ONLY apply if the RowSource IS an ArtifactSource
		// TODO #design we may be able to remove the need for this if we can handle JSON generically
		artifact_source.WithArtifactMapper(mappers.NewCloudtrailMapper()),
	}

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
	return AWSCloudTrail{}
}

func (c *CloudTrailLogTable) GetConfigSchema() parse.Config {
	return &CloudTrailLogTableConfig{}
}

// EnrichRow implements table.Table
func (c *CloudTrailLogTable) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// row must be an AWSCloudTrail
	record, ok := row.(AWSCloudTrail)
	if !ok {
		return nil, fmt.Errorf("invalid row type %T, expected AWSCloudTrail", row)
	}

	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		record.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	record.TpID = xid.New().String()
	record.TpSourceType = "aws_cloudtrail_log"
	record.TpTimestamp = record.EventTime
	record.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))
	if record.SourceIPAddress != nil {
		record.TpSourceIP = record.SourceIPAddress
		record.TpIps = append(record.TpIps, *record.SourceIPAddress)
	}
	for _, resource := range record.Resources {
		if resource.ARN != nil {
			newAkas := util.AwsAkasFromArn(*resource.ARN)
			record.TpAkas = append(record.TpAkas, newAkas...)
		}
	}
	// If it's an AKIA, then record that as an identity. Do not record ASIA*
	// keys etc.
	if record.UserIdentity.AccessKeyId != nil {
		if strings.HasPrefix(*record.UserIdentity.AccessKeyId, "AKIA") {
			record.TpUsernames = append(record.TpUsernames, *record.UserIdentity.AccessKeyId)
		}
	}
	if record.UserIdentity.UserName != nil {
		record.TpUsernames = append(record.TpUsernames, *record.UserIdentity.UserName)
	}

	// Hive fields
	record.TpPartition = "default" // TODO - should be based on the definition in HCL
	record.TpIndex = record.RecipientAccountId
	// convert to date in format yy-mm-dd
	record.TpDate = time.UnixMilli(int64(record.EventTime)).Format("2006-01-02")

	return record, nil
}
