// File: rows/alb_access_log.go

package rows

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

// AlbAccessLog represents a single Application Load Balancer (ALB) access log entry.
// It includes both AWS ALB-specific fields and Tailpipe enrichment fields for enhanced analysis.
// The struct maps directly to the format of ALB access logs as documented by AWS:
// https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-access-logs.html
//
// Data flow:
// 1. Raw log parsing: InitialiseFromMap parses raw string values into typed struct fields
// 2. Enrichment: EnrichRow adds derived and computed fields for analysis
//
// Enrichment fields are populated as follows:
// - tp_source_ip: from client IP
// - tp_destination_ip: from target IP
// - tp_ips: array containing both client and target IPs
// - tp_domains: from domain_name field
// - tp_akas: from target_group_arn for AWS resource linking
// - tp_timestamp: from timestamp field
// - tp_id: generated unique identifier
// - tp_source_type: set to "aws_alb_access_log"
// - tp_ingest_timestamp: set at processing time
// - tp_partition: set to "aws_alb_access_log"
// - tp_index: set to ALB name
// - tp_date: derived from timestamp in yyyy-mm-dd format

type AlbAccessLog struct {
	enrichment.CommonFields

	// Standard ALB fields
	Type                   string    `json:"type"`
	Timestamp              time.Time `json:"timestamp"`
	Alb                    string    `json:"alb"`
	ClientIP               string    `json:"client_ip"`
	ClientPort             int       `json:"client_port"`
	TargetIP               *string   `json:"target_ip,omitempty"`
	TargetPort             int       `json:"target_port"`
	RequestProcessingTime  float64   `json:"request_processing_time"`
	TargetProcessingTime   float64   `json:"target_processing_time"`
	ResponseProcessingTime float64   `json:"response_processing_time"`
	AlbStatusCode          *int      `json:"alb_status_code,omitempty"`
	TargetStatusCode       *int      `json:"target_status_code,omitempty"`
	ReceivedBytes          *int64    `json:"received_bytes,omitempty"`
	SentBytes              *int64    `json:"sent_bytes,omitempty"`
	Request                string    `json:"request"`
	UserAgent              string    `json:"user_agent"`
	SslCipher              string    `json:"ssl_cipher"`
	SslProtocol            string    `json:"ssl_protocol"`
	TargetGroupArn         string    `json:"target_group_arn"`
	TraceId                string    `json:"trace_id"`
	DomainName             string    `json:"domain_name"`
	ChosenCertArn          string    `json:"chosen_cert_arn"`
	MatchedRulePriority    int       `json:"matched_rule_priority"`
	RequestCreationTime    time.Time `json:"request_creation_time"`
	ActionsExecuted        string    `json:"actions_executed"`
	RedirectURL            *string   `json:"redirect_url,omitempty"`
	ErrorReason            *string   `json:"error_reason,omitempty"`
	TargetList             *string   `json:"target_list,omitempty"`
	TargetStatusList       *string   `json:"target_status_list,omitempty"`
	Classification         *string   `json:"classification,omitempty"`
	ClassificationReason   *string   `json:"classification_reason,omitempty"`
	ConnTraceID            *string   `json:"conn_trace_id,omitempty"`
}

// InitialiseFromMap initializes the struct from a map of string values
func (l *AlbAccessLog) InitialiseFromMap(m map[string]string) error {
	var err error
	for key, value := range m {
		if value == "-" {
			continue
		}
		switch key {
		case "timestamp":
			ts, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("error parsing timestamp: %w", err)
			}
			l.Timestamp = ts
		case "type":
			l.Type = value
		case "alb":
			l.Alb = value
		case "client":
			if strings.Contains(value, ":") {
				parts := strings.Split(value, ":")
				ip := parts[0]
				l.ClientIP = ip
				l.ClientPort, err = strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("error parsing client_port: %w", err)
				}
			}
		case "target":
			if strings.Contains(value, ":") {
				parts := strings.Split(value, ":")
				ip := parts[0]
				l.TargetIP = &ip
				l.TargetPort, err = strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("error parsing target_port: %w", err)
				}
			}
		case "request_processing_time":
			l.RequestProcessingTime, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("error parsing request_processing_time: %w", err)
			}
		case "target_processing_time":
			l.TargetProcessingTime, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("error parsing target_processing_time: %w", err)
			}
		case "response_processing_time":
			l.ResponseProcessingTime, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("error parsing response_processing_time: %w", err)
			}
		case "alb_status_code":
			asc, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing elb_status_code: %w", err)
			}
			l.AlbStatusCode = &asc
		case "target_status_code":
			tsc, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing target_status_code: %w", err)
			}
			l.TargetStatusCode = &tsc
		case "received_bytes":
			rb, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing received_bytes: %w", err)
			}
			l.ReceivedBytes = &rb
		case "sent_bytes":
			sb, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing sent_bytes: %w", err)
			}
			l.SentBytes = &sb
		case "request":
			l.Request = value
		case "user_agent":
			l.UserAgent = value
		case "ssl_cipher":
			l.SslCipher = value
		case "ssl_protocol":
			l.SslProtocol = value
		case "target_group_arn":
			l.TargetGroupArn = value
		case "trace_id":
			l.TraceId = value
		case "domain_name":
			l.DomainName = value
		case "chosen_cert_arn":
			l.ChosenCertArn = value
		case "matched_rule_priority":
			l.MatchedRulePriority, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing matched_rule_priority: %w", err)
			}
		case "request_creation_time":
			l.RequestCreationTime, err = time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("error parsing request_creation_time: %w", err)
			}
		case "actions_executed":
			l.ActionsExecuted = value
		case "redirect_url":
			l.RedirectURL = &value
		case "error_reason":
			l.ErrorReason = &value
		case "target_list":
			l.TargetList = &value
		case "target_status_list":
			l.TargetStatusList = &value
		case "classification":
			l.Classification = &value
		case "classification_reason":
			l.ClassificationReason = &value
		case "conn_trace_id":
			l.ConnTraceID = &value
		}
	}
	return nil
}

// EnrichRow handles all enrichment operations for the ALB access log entry
func (l *AlbAccessLog) EnrichRow(sourceEnrichmentFields *enrichment.CommonFields) error {
	// Add source enrichment fields if provided
	if sourceEnrichmentFields != nil {
		l.CommonFields = *sourceEnrichmentFields
	}

	// Standard record enrichment
	l.TpID = xid.New().String()
	l.TpTimestamp = l.Timestamp
	l.TpIngestTimestamp = time.Now()
	// truncate timestamp to date
	l.TpDate = l.Timestamp.Truncate(24 * time.Hour)
	// Use ALB name as the index
	l.TpIndex = l.Alb

	// IP-related enrichment
	if l.ClientIP != "" {
		l.TpSourceIP = &l.ClientIP
		l.TpIps = append(l.TpIps, l.ClientIP)
	}
	if l.TargetIP != nil {
		l.TpDestinationIP = l.TargetIP
		l.TpIps = append(l.TpIps, *l.TargetIP)
	}

	// Domain enrichment
	if l.DomainName != "" {
		l.TpDomains = append(l.TpDomains, l.DomainName)
	}

	// AWS resource linking
	if l.TargetGroupArn != "" {
		l.TpAkas = append(l.TpAkas, l.TargetGroupArn)
	}

	return nil
}
