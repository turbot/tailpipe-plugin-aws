package waf_traffic_log

import (
	"strings"
	"time"

	"github.com/rs/xid"

	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/cloudwatch_log_group"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-aws/tables"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const WafTrafficLogTableIdentifier = "aws_waf_traffic_log"

// WafTrafficLogTable - table for Waf traffic logs
type WafTrafficLogTable struct{}

func (c *WafTrafficLogTable) GetSourceMetadata() ([]*table.SourceMetadata[*WafTrafficLog], error) {
	// the default file layout for Waf traffic logs in S3
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/WAFLogs/%{DATA:cloudfront_or_region}/%{DATA:cloudfront_name_or_resource_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.log.gz"),
	}

	return []*table.SourceMetadata[*WafTrafficLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     &WafMapper{},
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     &WafMapper{},
			Options:    []row_source.RowSourceOption{artifact_source.WithRowPerLine()},
		},
		{
			SourceName: cloudwatch_log_group.AwsCloudwatchLogGroupSourceIdentifier,
			Mapper:     &WafMapper{},
		},
	}, nil
}

func (c *WafTrafficLogTable) Identifier() string {
	return WafTrafficLogTableIdentifier
}

// EnrichRow implements table.Table
func (c *WafTrafficLogTable) EnrichRow(row *WafTrafficLog, sourceEnrichmentFields schema.SourceEnrichment) (*WafTrafficLog, error) { // initialize the enrichment fields to any fields provided by the source
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpTimestamp = *row.Timestamp
	row.TpIngestTimestamp = time.Now()
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)

	if row.HttpSourceId != nil {
		// For ALB and APPService we are getting ID not ARN.
		if strings.HasPrefix(*row.HttpSourceId, "arn:") {
			newAkas := tables.AwsAkasFromArn(*row.HttpSourceId)
			row.TpAkas = append(row.TpAkas, newAkas...)
		}
	}

	if row.HttpRequest != nil {
		if row.HttpRequest.ClientIp != nil {
			row.TpIps = append(row.TpIps, *row.HttpRequest.ClientIp)
		}
	}

	return row, nil
}

func (c *WafTrafficLogTable) GetDescription() string {
	return "AWS WAF traffic logs record detailed web request data, helping analyze threats, monitor rule effectiveness, and improve security posture."
}
