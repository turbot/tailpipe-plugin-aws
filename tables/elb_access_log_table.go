package tables

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const ElbAccessLogTableIdentifier = "aws_elb_access_log"

func init() {
	// Register the table, with type parameters:
	// 1. row struct
	// 2. table config struct
	// 3. table implementation
	table.RegisterTable[*rows.ElbAccessLog, *ElbAccessLogTableConfig, *ElbAccessLogTable]()
}

const elbLogFormat = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason" $conn_trace_id`
const elbLogFormatNoConnTrace = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason"`

type ElbAccessLogTable struct{}

func (c *ElbAccessLogTable) Identifier() string {
	return ElbAccessLogTableIdentifier
}

func (c *ElbAccessLogTable) GetSourceMetadata(_ *ElbAccessLogTableConfig) []*table.SourceMetadata[*rows.ElbAccessLog] {
	return []*table.SourceMetadata[*rows.ElbAccessLog]{
		{
			SourceName: constants.ArtifactSourceIdentifier,
			Mapper:     table.NewRowPatternMapper[*rows.ElbAccessLog](elbLogFormat, elbLogFormatNoConnTrace),
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}

}

func (c *ElbAccessLogTable) EnrichRow(row *rows.ElbAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.ElbAccessLog, error) {
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()
	row.TpTimestamp = row.Timestamp
	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)

	row.TpIndex = row.Elb // TODO: #enrichment figure out how to get account id / other index

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
