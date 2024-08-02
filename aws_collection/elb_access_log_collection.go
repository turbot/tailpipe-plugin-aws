package aws_collection

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/hcl"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"strconv"
	"time"
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
		case "client_ip":
			record.ClientIP = value
			record.TpSourceIP = &value
			record.TpIps = append(record.TpIps, value)
		case "client_port":
			cp, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("error parsing client_port: %w", err)
			}
			record.ClientPort = cp
		case "target_ip":
			record.TargetIP = value
			record.TpDestinationIP = &value
			record.TpIps = append(record.TpIps, value)
		case "target_port":
			tp, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("error parsing target_port: %w", err)
			}
			record.TargetPort = tp
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
			esc, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("error parsing elb_status_code: %w", err)
			}
			record.ElbStatusCode = esc
		case "target_status_code":
			tsc, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("error parsing target_status_code: %w", err)
			}
			record.TargetStatusCode = tsc
		case "received_bytes":
			rb, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing received_bytes: %w", err)
			}
			record.ReceivedBytes = rb
		case "sent_bytes":
			sb, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing sent_bytes: %w", err)
			}
			record.SentBytes = sb
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
			mrp, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("error parsing matched_rule_priority: %w", err)
			}
			record.MatchedRulePriority = mrp
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
