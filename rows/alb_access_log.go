// File: rows/alb_access_log.go

package rows

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
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
	schema.CommonFields

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

func (l *AlbAccessLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"timestamp":                "The time when the load balancer received the request from the client, in ISO 8601 format.",
		"type":                     "The type of log entry.",
		"alb":                      "The Amazon Resource Name (ARN) of the load balancer.",
		"client":                   "The IP address and port of the client that made the request.",
		"target":                   "The IP address and port of the target that the request was forwarded to.",
		"request_processing_time":  "The total time elapsed from the time the load balancer received the request to the time it sent the request to a target.",
		"target_processing_time":   "The total time elapsed from the time the load balancer sent the request to a target to the time the target started to send the response headers.",
		"response_processing_time": "The total time elapsed from the time the load balancer received the response headers from the target to the time it started to send the response to the client.",
		"alb_status_code":          "The status code of the response from the load balancer.",
		"target_status_code":       "The status code of the response from the target.",
		"received_bytes":           "The number of bytes received by the load balancer from the client.",
		"sent_bytes":               "The number of bytes sent by the load balancer to the client.",
		"request":                  "The request string.",
		"user_agent":               "The user agent of the client.",
		"ssl_cipher":               "The SSL cipher.",
		"ssl_protocol":             "The SSL protocol.",
		"target_group_arn":         "The Amazon Resource Name (ARN) of the target group.",
		"trace_id":                 "The trace ID.",
		"domain_name":              "The domain name.",
		"chosen_cert_arn":          "The Amazon Resource Name (ARN) of the chosen certificate.",
		"matched_rule_priority":    "The priority of the rule that matched the request.",
		"request_creation_time":    "The time when the request was created.",
		"actions_executed":         "The actions executed.",
		"redirect_url":             "The redirect URL.",
		"error_reason":             "The error reason.",
		"target_list":              "The target list.",
		"target_status_list":       "The target status list.",
		"classification":           "The classification.",
		"classification_reason":    "The classification reason.",
		"conn_trace_id":            "The connection trace ID.",
	}
}
