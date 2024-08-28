package aws_collection

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/aws_source"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-aws/util"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
)

// CloudTrailLogCollection - collection for CloudTrail logs
type CloudTrailLogCollection struct {
	// all collections must embed collection.CollectionBase
	collection.CollectionBase[*CloudTrailLogCollectionConfig]
}

func NewCloudTrailLogCollection() collection.Collection {
	return &CloudTrailLogCollection{}
}

// Identifier implements collection.Collection
func (c *CloudTrailLogCollection) Identifier() string {
	return "aws_cloudtrail_log"
}

// GetSourceOptions returns any options which should be passed to the given source type
func (c *CloudTrailLogCollection) GetSourceOptions() []row_source.RowSourceOption {
	// the defaulkt file layout for Cloudtrail logs in S3
	defaultArtifactConfig := &artifact_source.ArtifactSourceConfigBase{
		FileLayout: utils.ToStringPointer("AWSLogs/o-z3cf4qoe7m/\\d+/CloudTrail/[a-z-0-9]+/\\d{4}/\\d{2}/\\d{2}/(?P<index>\\d+)_CloudTrail_(?P<region>[a-z-0-9]+)_(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})T(?P<hour>\\d{2})(?P<minute>\\d{2})Z_.+.json.gz"),
		//JsonPath:    nil,
	}

	return []row_source.RowSourceOption{
		// if the source is an artifact source, we need a mapper
		// NOTE: WithArtifactMapper option will ONLY apply if the RowSource IS an ArtifactSource
		// TODO #design we may be able to remove the need for this if we can handle JSON generically
		artifact_source.WithArtifactMapper(aws_source.NewCloudtrailMapper()),

		// default file layout for CloudTrail logs in S3
		// TODO check if source is S3???
		artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig),
	}
}

// GetRowSchema implements collection.Collection
func (c *CloudTrailLogCollection) GetRowSchema() any {
	return aws_types.AWSCloudTrail{}
}

func (c *CloudTrailLogCollection) GetConfigSchema() parse.Config {
	return &CloudTrailLogCollectionConfig{}
}

// EnrichRow implements collection.Collection
func (c *CloudTrailLogCollection) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// row must be an AWSCloudTrail
	record, ok := row.(aws_types.AWSCloudTrail)
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
	record.TpCollection = "default" // TODO - should be based on the definition in HCL
	record.TpConnection = record.RecipientAccountId
	record.TpYear = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Year())
	record.TpMonth = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Month())
	record.TpDay = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Day())

	return record, nil
}
