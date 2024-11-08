// File: rows/alb_access_log.go

package rows

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
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
	AlbName                string    `json:"alb_name"`
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
	RedirectUrl            *string   `json:"redirect_url,omitempty"`
	ErrorReason            *string   `json:"error_reason,omitempty"`
	TargetList             *string   `json:"target_list,omitempty"`
	TargetStatusList       *string   `json:"target_status_list,omitempty"`
	Classification         *string   `json:"classification,omitempty"`
	ClassificationReason   *string   `json:"classification_reason,omitempty"`
}

// NewAlbAccessLog creates a new ALB access log entry with initialized fields.
// Used by the mapper when creating new log entries from raw log lines.
func NewAlbAccessLog() *AlbAccessLog {
	return &AlbAccessLog{}
}

// parseField handles parsing of individual fields with proper error handling
func (l *AlbAccessLog) parseField(fieldName, value string) error {
	if value == "-" {
		return nil
	}

	switch fieldName {
	case "timestamp":
		ts, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return fmt.Errorf("error parsing timestamp: %w", err)
		}
		l.Timestamp = ts
	case "type":
		l.Type = value
	case "alb":
		l.AlbName = value
	case "client":
		parts := strings.Split(value, ":")
		if len(parts) == 2 {
			l.ClientIP = parts[0]
			port, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("error parsing client port: %w", err)
			}
			l.ClientPort = port
		}
	case "target":
		parts := strings.Split(value, ":")
		if len(parts) == 2 {
			ip := parts[0]
			l.TargetIP = &ip
			port, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("error parsing target port: %w", err)
			}
			l.TargetPort = port
		}
	case "request_processing_time":
		rpt, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("error parsing request_processing_time: %w", err)
		}
		l.RequestProcessingTime = rpt
	case "target_processing_time":
		tpt, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("error parsing target_processing_time: %w", err)
		}
		l.TargetProcessingTime = tpt
	case "response_processing_time":
		rpt, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("error parsing response_processing_time: %w", err)
		}
		l.ResponseProcessingTime = rpt
	case "alb_status_code":
		code, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("error parsing alb_status_code: %w", err)
		}
		l.AlbStatusCode = &code
	case "target_status_code":
		code, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("error parsing target_status_code: %w", err)
		}
		l.TargetStatusCode = &code
	case "received_bytes":
		bytes, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing received_bytes: %w", err)
		}
		l.ReceivedBytes = &bytes
	case "sent_bytes":
		bytes, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing sent_bytes: %w", err)
		}
		l.SentBytes = &bytes
	case "request":
		l.Request = strings.Trim(value, "\"")
	case "user_agent":
		l.UserAgent = strings.Trim(value, "\"")
	case "ssl_cipher":
		l.SslCipher = value
	case "ssl_protocol":
		l.SslProtocol = value
	case "target_group_arn":
		l.TargetGroupArn = value
	case "trace_id":
		l.TraceId = strings.Trim(value, "\"")
	case "domain_name":
		l.DomainName = strings.Trim(value, "\"")
	case "chosen_cert_arn":
		l.ChosenCertArn = strings.Trim(value, "\"")
	case "matched_rule_priority":
		priority, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("error parsing matched_rule_priority: %w", err)
		}
		l.MatchedRulePriority = priority
	case "request_creation_time":
		rct, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return fmt.Errorf("error parsing request_creation_time: %w", err)
		}
		l.RequestCreationTime = rct
	case "actions_executed":
		l.ActionsExecuted = strings.Trim(value, "\"")
	case "redirect_url":
		if value != "-" {
			v := strings.Trim(value, "\"")
			l.RedirectUrl = &v
		}
	case "error_reason":
		if value != "-" {
			v := strings.Trim(value, "\"")
			l.ErrorReason = &v
		}
	case "target_list":
		if value != "-" {
			v := strings.Trim(value, "\"")
			l.TargetList = &v
		}
	case "target_status_list":
		if value != "-" {
			v := strings.Trim(value, "\"")
			l.TargetStatusList = &v
		}
	case "classification":
		if value != "-" {
			v := strings.Trim(value, "\"")
			l.Classification = &v
		}
	case "classification_reason":
		if value != "-" {
			v := strings.Trim(value, "\"")
			l.ClassificationReason = &v
		}
	}
	return nil
}

// InitialiseFromMap initializes the struct from a map of string values
func (l *AlbAccessLog) InitialiseFromMap(m map[string]string) error {
	for fieldName, value := range m {
		if err := l.parseField(fieldName, value); err != nil {
			return fmt.Errorf("error parsing field %s: %w", fieldName, err)
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
	l.TpTimestamp = time.Unix(0, int64(helpers.UnixMillis(l.Timestamp.UnixNano()/int64(time.Millisecond)))*int64(time.Millisecond))
	l.TpIngestTimestamp = time.Unix(0, int64(helpers.UnixMillis(time.Now().UnixNano()/int64(time.Millisecond)))*int64(time.Millisecond))
	l.TpDate = l.Timestamp.Format("2006-01-02")
	// Use ALB name as the index
	l.TpIndex = l.AlbName

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
