package rows

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type ElbAccessLog struct {
	enrichment.CommonFields

	Type                   string    `json:"type"`
	Timestamp              time.Time `json:"timestamp"`
	Elb                    string    `json:"elb"`
	ClientIP               string    `json:"client_ip"`
	ClientPort             int       `json:"client_port"`
	TargetIP               *string   `json:"target_ip,omitempty"`
	TargetPort             int       `json:"target_port,omitempty"`
	RequestProcessingTime  float64   `json:"request_processing_time"`
	TargetProcessingTime   float64   `json:"target_processing_time"`
	ResponseProcessingTime float64   `json:"response_processing_time"`
	ElbStatusCode          *int      `json:"elb_status_code,omitempty"`
	TargetStatusCode       *int      `json:"target_status_code,omitempty"`
	ReceivedBytes          *int64    `json:"received_bytes"`
	SentBytes              *int64    `json:"sent_bytes"`
	Request                string    `json:"request"`
	UserAgent              string    `json:"user_agent"`
	SslCipher              string    `json:"ssl_cipher"`
	SslProtocol            string    `json:"ssl_protocol"`
	TargetGroupArn         string    `json:"target_group_arn"`
	TraceID                string    `json:"trace_id"`
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

// InitialiseFromMap - initialise the struct from a map
func (l *ElbAccessLog) InitialiseFromMap(m map[string]string) error {
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
		case "elb":
			l.Elb = value
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
		case "elb_status_code":
			esc, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing elb_status_code: %w", err)
			}
			l.ElbStatusCode = &esc
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
			l.TraceID = value
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
