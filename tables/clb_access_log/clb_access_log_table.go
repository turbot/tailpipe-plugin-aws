package alb_access_log

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

const ClbAccessLogTableIdentifier = "aws_clb_access_log"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*ClbAccessLog, *ClbAccessLogTable]()
}

const albLogFormat = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason" $conn_trace_id`

type ClbAccessLogTable struct{}

func (c *ClbAccessLogTable) Identifier() string {
	return ClbAccessLogTableIdentifier
}

func (c *ClbAccessLogTable) GetSourceMetadata() []*table.SourceMetadata[*ClbAccessLog] {
	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log.gz"),
	}

	return []*table.SourceMetadata[*ClbAccessLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*ClbAccessLog](albLogFormat),
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     mappers.NewGonxMapper[*ClbAccessLog](albLogFormat),
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

func (c *ClbAccessLogTable) EnrichRow(row *ClbAccessLog, sourceEnrichmentFields schema.SourceEnrichment) (*ClbAccessLog, error) {
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

func (c *ClbAccessLogTable) GetDescription() string {
	return "AWS ALB Access logs capture detailed information about the requests that are processed by an Application Load Balancer. This table provides a structured representation of the log data, including request and response details, client and target information, processing times, and security parameters."
}
