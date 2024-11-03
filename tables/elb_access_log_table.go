package tables

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_mapper"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type ElbAccessLogTable struct {
	table.TableImpl[map[string]string, *ElbAccessLogTableConfig, *config.AwsConnection]
}

func NewElbAccessLogTable() table.Table {
	return &ElbAccessLogTable{}
}

func (c *ElbAccessLogTable) Identifier() string {
	return "aws_elb_access_log"
}

func (c *ElbAccessLogTable) GetRowSchema() any {
	return &rows.AwsElbAccessLog{}
}

func (c *ElbAccessLogTable) GetConfigSchema() parse.Config {
	return &ElbAccessLogTableConfig{}
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
	// TODO switch on source

	// TODO KAI make sure tables add NewCloudwatchMapper if needed
	// NOTE: add the cloudwatch mapper to ensure rows are in correct format
	//s.AddMappers(artifact_mapper.NewCloudwatchMapper())

	// if the source is an artifact source, we need a mapper
	c.Mappers = []artifact_mapper.Mapper{mappers.NewElbAccessLogMapper()}
}

func (c *ElbAccessLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
	}
}

func (c *ElbAccessLogTable) EnrichRow(rawRecord map[string]string, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// TODO: #validate ensure we have a timestamp field

	// Build record and add any source enrichment fields
	var row rows.AwsElbAccessLog
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

	for key, value := range rawRecord {
		switch key {
		case "type":
			row.Type = value
		case "timestamp":
			ts, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return nil, fmt.Errorf("error parsing timestamp: %w", err)
			}
			row.Timestamp = ts
			row.TpTimestamp = helpers.UnixMillis(ts.UnixNano() / int64(time.Millisecond))
		case "elb":
			row.Elb = value
		case "client":
			if value != "-" && strings.Contains(value, ":") {
				ip := strings.Split(value, ":")[0]
				row.ClientIP = ip
				row.TpSourceIP = &ip
				row.TpIps = append(row.TpIps, ip)
				row.ClientPort, _ = strconv.Atoi(strings.Split(value, ":")[1])
			}
		case "target":
			if value != "-" && strings.Contains(value, ":") {
				ip := strings.Split(value, ":")[0]
				row.TargetIP = &ip
				row.TpDestinationIP = &ip
				row.TpIps = append(row.TpIps, ip)
				row.TargetPort, _ = strconv.Atoi(strings.Split(value, ":")[1])
			}
		case "request_processing_time":
			rpt, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing request_processing_time: %w", err)
			}
			row.RequestProcessingTime = rpt
		case "target_processing_time":
			tpt, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing target_processing_time: %w", err)
			}
			row.TargetProcessingTime = tpt
		case "response_processing_time":
			rpt, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing response_processing_time: %w", err)
			}
			row.ResponseProcessingTime = rpt
		case "elb_status_code":
			if value != "-" {
				esc, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing elb_status_code: %w", err)
				}
				row.ElbStatusCode = &esc
			}
		case "target_status_code":
			if value != "-" {
				tsc, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing target_status_code: %w", err)
				}
				row.TargetStatusCode = &tsc
			}
		case "received_bytes":
			if value != "-" {
				rb, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing received_bytes: %w", err)
				}
				row.ReceivedBytes = &rb
			}
		case "sent_bytes":
			if value != "-" {
				sb, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing sent_bytes: %w", err)
				}
				row.SentBytes = &sb
			}
		case "request":
			row.Request = value
		case "user_agent":
			row.UserAgent = value
		case "ssl_cipher":
			row.SslCipher = value
		case "ssl_protocol":
			row.SslProtocol = value
		case "target_group_arn":
			row.TargetGroupArn = value
		case "trace_id":
			row.TraceID = value
		case "domain_name":
			row.DomainName = value
			row.TpDomains = append(row.TpDomains, value)
		case "chosen_cert_arn":
			row.ChosenCertArn = value
		case "matched_rule_priority":
			if value != "-" {
				mrp, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing matched_rule_priority: %w", err)
				}
				row.MatchedRulePriority = mrp
			}
		case "request_creation_time":
			rct, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return nil, fmt.Errorf("error parsing request_creation_time: %w", err)
			}
			row.RequestCreationTime = rct
		case "actions_executed":
			row.ActionsExecuted = value
		case "redirect_url":
			row.RedirectURL = &value
		case "error_reason":
			row.ErrorReason = &value
		case "target_list":
			row.TargetList = &value
		case "target_status_list":
			row.TargetStatusList = &value
		case "classification":
			row.Classification = &value
		case "classification_reason":
			row.ClassificationReason = &value
		case "conn_trace_id":
			row.ConnTraceID = value
		}

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
