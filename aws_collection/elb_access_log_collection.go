package aws_collection

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_source"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/hcl"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
)

type ElbAccessLogCollection struct {
	collection.CollectionBase[*ElbAccessLogCollectionConfig]
}

func NewElbAccessLogCollection() collection.Collection {
	return &ElbAccessLogCollection{}
}

func (c *ElbAccessLogCollection) Identifier() string {
	return "aws_elb_access_log"
}

func (c *ElbAccessLogCollection) GetRowSchema() any {
	return &aws_types.AwsElbAccessLog{}
}

func (c *ElbAccessLogCollection) GetConfigSchema() hcl.Config {
	return &ElbAccessLogCollectionConfig{}
}

func (c *ElbAccessLogCollection) GetSourceOptions() []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
		artifact_source.WithMapper(aws_source.NewElbAccessLogMapper()),
	}
}

func (c *ElbAccessLogCollection) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// short-circuit for unexpected row type
	rawRecord, ok := row.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("invalid row type: %T, expected map[string]string", row)
	}

	// TODO: #validate ensure we have a timestamp field

	// Build record and add any source enrichment fields
	var record aws_types.AwsElbAccessLog
	if sourceEnrichmentFields != nil {
		record.CommonFields = *sourceEnrichmentFields
	}

	for key, value := range rawRecord {
		switch key {
		case "type":
			record.Type = value
		case "timestamp":
			ts, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return nil, fmt.Errorf("error parsing timestamp: %w", err)
			}
			record.Timestamp = ts
		case "elb":
			record.Elb = value
		case "client":
			if value != "-" && strings.Contains(value, ":") {
				ip := strings.Split(value, ":")[0]
				record.ClientIP = ip
				record.TpSourceIP = &ip
				record.TpIps = append(record.TpIps, ip)
				record.ClientPort, _ = strconv.Atoi(strings.Split(value, ":")[1])
			}
		case "target":
			if value != "-" && strings.Contains(value, ":") {
				ip := strings.Split(value, ":")[0]
				record.TargetIP = &ip
				record.TpDestinationIP = &ip
				record.TpIps = append(record.TpIps, ip)
				record.TargetPort, _ = strconv.Atoi(strings.Split(value, ":")[1])
			}
		case "request_processing_time":
			rpt, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing request_processing_time: %w", err)
			}
			record.RequestProcessingTime = rpt
		case "target_processing_time":
			tpt, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing target_processing_time: %w", err)
			}
			record.TargetProcessingTime = tpt
		case "response_processing_time":
			rpt, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing response_processing_time: %w", err)
			}
			record.ResponseProcessingTime = rpt
		case "elb_status_code":
			if value != "-" {
				esc, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing elb_status_code: %w", err)
				}
				record.ElbStatusCode = &esc
			}
		case "target_status_code":
			if value != "-" {
				tsc, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing target_status_code: %w", err)
				}
				record.TargetStatusCode = &tsc
			}
		case "received_bytes":
			if value != "-" {
				rb, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing received_bytes: %w", err)
				}
				record.ReceivedBytes = &rb
			}
		case "sent_bytes":
			if value != "-" {
				sb, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing sent_bytes: %w", err)
				}
				record.SentBytes = &sb
			}
		case "request":
			record.Request = value
		case "user_agent":
			record.UserAgent = value
		case "ssl_cipher":
			record.SslCipher = value
		case "ssl_protocol":
			record.SslProtocol = value
		case "target_group_arn":
			record.TargetGroupArn = value
		case "trace_id":
			record.TraceID = value
		case "domain_name":
			record.DomainName = value
			record.TpDomains = append(record.TpDomains, value)
		case "chosen_cert_arn":
			record.ChosenCertArn = value
		case "matched_rule_priority":
			if value != "-" {
				mrp, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("error parsing matched_rule_priority: %w", err)
				}
				record.MatchedRulePriority = mrp
			}
		case "request_creation_time":
			rct, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return nil, fmt.Errorf("error parsing request_creation_time: %w", err)
			}
			record.RequestCreationTime = rct
		case "actions_executed":
			record.ActionsExecuted = value
		case "redirect_url":
			record.RedirectURL = &value
		case "error_reason":
			record.ErrorReason = &value
		case "target_list":
			record.TargetList = &value
		case "target_status_list":
			record.TargetStatusList = &value
		case "classification":
			record.Classification = &value
		case "classification_reason":
			record.ClassificationReason = &value
		case "conn_trace_id":
			record.ConnTraceID = value
		}

	}

	// Record standardization
	record.TpID = xid.New().String()
	record.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))
	record.TpSourceType = "aws_elb_access_log" // TODO: #refactor move to source?

	// Hive Fields
	record.TpCollection = c.Identifier()
	if record.TpConnection == "" {
		record.TpConnection = c.Identifier() // TODO: #refactor figure out how to get connection (account ID?)
	}
	record.TpYear = int32(record.Timestamp.Year())
	record.TpMonth = int32(record.Timestamp.Month())
	record.TpDay = int32(record.Timestamp.Day())

	return record, nil
}
