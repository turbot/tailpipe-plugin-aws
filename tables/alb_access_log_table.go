package tables

import (
	"time"

	"github.com/rs/xid"

	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-aws/sources"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const AlbAccessLogTableIdentifier = "aws_alb_access_log"

const albLogFormat = `$type $timestamp $alb $client $target $request_processing_time $target_processing_time $response_processing_time $alb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason" $conn_trace_id`
const albLogFormatNoConnTrace = `$type $timestamp $alb $client $target $request_processing_time $target_processing_time $response_processing_time $alb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason"`

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.AlbAccessLog, *AlbAccessLogTable]()
}

type AlbAccessLogTable struct{}

func (c *AlbAccessLogTable) GetSourceMetadata() []*table.SourceMetadata[*rows.AlbAccessLog] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA:elb_name}_%{TIMESTAMP_ISO8601:end_time}_%{DATA:suffix}.log"),
	}

	return []*table.SourceMetadata[*rows.AlbAccessLog]{
		{
			// S3 artifact source
			SourceName: sources.AwsS3BucketSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*rows.AlbAccessLog](albLogFormat, albLogFormatNoConnTrace),
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*rows.AlbAccessLog](albLogFormat, albLogFormatNoConnTrace),
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *AlbAccessLogTable) Identifier() string {
	return AlbAccessLogTableIdentifier
}

func (c *AlbAccessLogTable) EnrichRow(row *rows.AlbAccessLog, sourceEnrichmentFields schema.SourceEnrichment) (*rows.AlbAccessLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Standard record enrichment
	row.TpID = xid.New().String()
	row.TpTimestamp = row.Timestamp
	row.TpIngestTimestamp = time.Now()
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)

	row.TpIndex = row.Alb // TODO: #enrichment figure out how to get the account id / better index

	// IP-related enrichment
	if row.ClientIP != "" {
		row.TpSourceIP = &row.ClientIP
		row.TpIps = append(row.TpIps, row.ClientIP)
	}
	if row.TargetIP != nil {
		row.TpDestinationIP = row.TargetIP
		row.TpIps = append(row.TpIps, *row.TargetIP)
	}
	// Domain enrichment
	if row.DomainName != "" {
		row.TpDomains = append(row.TpDomains, row.DomainName)
	}

	// AWS resource linking
	if row.TargetGroupArn != "" {
		row.TpAkas = append(row.TpAkas, row.TargetGroupArn)
	}
	return row, nil
}
