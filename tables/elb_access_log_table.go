package tables

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const elbLogFormat = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason" $conn_trace_id`
const elbLogFormatNoConnTrace = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason"`

type ElbAccessLogTable struct {
	table.TableImpl[*rows.ElbAccessLog, *ElbAccessLogTableConfig, *config.AwsConnection]
}

func NewElbAccessLogTable() table.Table {
	return &ElbAccessLogTable{}
}

func (c *ElbAccessLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
	// call base init
	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
		return err
	}

	c.initMappers()
	return nil
}

func (c *ElbAccessLogTable) initMappers() {
	// todo switch on source
	c.Mapper = table.NewDelimitedLineMapper(rows.NewElbAccessLog, elbLogFormat, elbLogFormatNoConnTrace)
}

func (c *ElbAccessLogTable) Identifier() string {
	return "aws_elb_access_log"
}

func (c *ElbAccessLogTable) GetRowSchema() any {
	return &rows.ElbAccessLog{}
}

func (c *ElbAccessLogTable) GetConfigSchema() parse.Config {
	return &ElbAccessLogTableConfig{}
}

func (c *ElbAccessLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
	}
}

func (c *ElbAccessLogTable) EnrichRow(row *rows.ElbAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.ElbAccessLog, error) {
	// TODO: #validate ensure we have a timestamp field

	// add any source enrichment fields
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))
	row.TpSourceType = "aws_elb_access_log" // TODO: #refactor move to source?

	// Hive Fields
	row.TpPartition = c.Identifier()
	if row.TpIndex == "" {
		row.TpIndex = c.Identifier() // TODO: #refactor figure out how to get connection (account ID?)
	}

	// convert to date in format yy-mm-dd
	row.TpDate = row.Timestamp.Format("2006-01-02")

	return row, nil
}
