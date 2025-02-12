package elb_access_log

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-aws/tables"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const ElbAccessLogTableIdentifier = "aws_elb_access_log"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*ElbAccessLog, *ElbAccessLogTable]()
}

const elbLogFormat = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason" $conn_trace_id`
const elbLogFormatNoConnTrace = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason"`

type ElbAccessLogTable struct{}

func (c *ElbAccessLogTable) Identifier() string {
	return ElbAccessLogTableIdentifier
}

func (c *ElbAccessLogTable) GetSourceMetadata() []*table.SourceMetadata[*ElbAccessLog] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log.gz"),
	}

	return []*table.SourceMetadata[*ElbAccessLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*ElbAccessLog](elbLogFormat, elbLogFormatNoConnTrace),
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*ElbAccessLog](elbLogFormat, elbLogFormatNoConnTrace),
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *ElbAccessLogTable) EnrichRow(row *ElbAccessLog, sourceEnrichmentFields schema.SourceEnrichment) (*ElbAccessLog, error) {
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	row.TpTimestamp = row.Timestamp
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)

	callerIdentityData, err := tables.GetCallerIdentityData()
	if err != nil {
		return nil, err
	}
	row.TpIndex = *callerIdentityData.Account

	row.TpSourceIP = &row.ClientIP
	row.TpIps = append(row.TpIps, row.ClientIP)
	if row.TargetIP != nil {
		row.TpDestinationIP = row.TargetIP
		row.TpIps = append(row.TpIps, *row.TargetIP)
	}

	if row.DomainName != "" {
		row.TpDomains = append(row.TpDomains, row.DomainName)
	}

	if row.TargetGroupArn != "" {
		row.TpAkas = append(row.TpAkas, row.TargetGroupArn)
	}

	return row, nil
}

func (c *ElbAccessLogTable) GetDescription() string {
	return "AWS ELB Access logs capture detailed information about the requests that are processed by an Elastic Load Balancer. This table provides a structured representation of the log data, including request and response details, client and target information, processing times, and security parameters."
}