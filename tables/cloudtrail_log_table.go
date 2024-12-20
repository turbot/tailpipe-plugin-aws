package tables

import (
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/extractors"
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

const CloudTrailLogTableIdentifier = "aws_cloudtrail_log"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.CloudTrailLog, *CloudTrailLogTable]()
}

// CloudTrailLogTable - table for CloudTrailLog logs
type CloudTrailLogTable struct{}

func (t *CloudTrailLogTable) GetSourceMetadata() []*table.SourceMetadata[*rows.CloudTrailLog] {
	// the default file layout for CloudTrail logs in S3
	defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
		FileLayout: utils.ToStringPointer("AWSLogs(?:/o-[a-z0-9]{8,12})?/\\d+/CloudTrail/[a-z-0-9]+/\\d{4}/\\d{2}/\\d{2}/(?P<index>\\d+)_CloudTrail_(?P<region>[a-z-0-9]+)_(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})T(?P<hour>\\d{2})(?P<minute>\\d{2})Z_.+.json.gz"),
	}

	return []*table.SourceMetadata[*rows.CloudTrailLog]{
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig),
				artifact_source.WithArtifactExtractor(extractors.NewCloudTrailLogExtractor()),
			},
		},
		{
			SourceName: sources.AwsCloudwatchSourceIdentifier,
			Mapper:     &mappers.CloudTrailMapper{},
		},
	}
}

// Identifier implements table.Table
func (t *CloudTrailLogTable) Identifier() string {
	return CloudTrailLogTableIdentifier
}

// EnrichRow implements table.Table
func (t *CloudTrailLogTable) EnrichRow(row *rows.CloudTrailLog, sourceEnrichmentFields schema.SourceEnrichment) (*rows.CloudTrailLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpTimestamp = *row.EventTime
	row.TpIngestTimestamp = time.Now()

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
	row.TpIndex = row.RecipientAccountId
	// convert to date in format yy-mm-dd
	row.TpDate = row.EventTime.Truncate(24 * time.Hour)

	return row, nil
}
