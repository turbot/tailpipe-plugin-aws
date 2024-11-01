package models

import (
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type AwsElbAccessLog struct {
	enrichment.CommonFields

	Type                   string    `json:"type" column:"type"`
	Timestamp              time.Time `json:"timestamp" column:"timestamp"`
	Elb                    string    `json:"elb" column:"elb"`
	ClientIP               string    `json:"client_ip" column:"client"`
	ClientPort             int       `json:"client_port" column:"client_port"`
	TargetIP               *string   `json:"target_ip,omitempty" column:"target"`
	TargetPort             int       `json:"target_port,omitempty" column:"target_port"`
	RequestProcessingTime  float64   `json:"request_processing_time" column:"request_processing_time"`
	TargetProcessingTime   float64   `json:"target_processing_time" column:"target_processing_time"`
	ResponseProcessingTime float64   `json:"response_processing_time" column:"response_processing_time"`
	ElbStatusCode          *int      `json:"elb_status_code,omitempty" column:"elb_status_code"`
	TargetStatusCode       *int      `json:"target_status_code,omitempty" column:"target_status_code"`
	ReceivedBytes          *int64    `json:"received_bytes" column:"received_bytes"`
	SentBytes              *int64    `json:"sent_bytes" column:"sent_bytes"`
	Request                string    `json:"request" column:"request"`
	UserAgent              string    `json:"user_agent" column:"user_agent"`
	SslCipher              string    `json:"ssl_cipher" column:"ssl_cipher"`
	SslProtocol            string    `json:"ssl_protocol" column:"ssl_protocol"`
	TargetGroupArn         string    `json:"target_group_arn" column:"target_group_arn"`
	TraceID                string    `json:"trace_id" column:"trace_id"`
	DomainName             string    `json:"domain_name" column:"domain_name"`
	ChosenCertArn          string    `json:"chosen_cert_arn" column:"chosen_cert_arn"`
	MatchedRulePriority    int       `json:"matched_rule_priority" column:"matched_rule_priority"`
	RequestCreationTime    time.Time `json:"request_creation_time" column:"request_creation_time"`
	ActionsExecuted        string    `json:"actions_executed" column:"actions_executed"`
	RedirectURL            *string   `json:"redirect_url,omitempty"	column:"redirect_url"`
	ErrorReason            *string   `json:"error_reason,omitempty" column:"error_reason"`
	TargetList             *string   `json:"target_list,omitempty" column:"target_list"`
	TargetStatusList       *string   `json:"target_status_list,omitempty" column:"target_status_list"`
	Classification         *string   `json:"classification,omitempty" column:"classification"`
	ClassificationReason   *string   `json:"classification_reason,omitempty" column:"classification_reason"`
	ConnTraceID            string    `json:"conn_trace_id" column:"conn_trace_id"`
}

// InitialiseFromMap - initialise the struct from a map
func (l *AwsElbAccessLog) InitialiseFromMap(m map[string]string) error {
	for key, value := range m {
		switch key {
		case "timestamp":
			ts, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("error parsing timestamp: %w", err)
			}
			l.Timestamp = ts
			l.TpTimestamp = helpers.UnixMillis(ts.UnixNano() / int64(time.Millisecond))
		case "type":
			l.Type = value
		case "elb":
			l.Elb = value
		case "client":
			if value != "-" && strings.Contains(value, ":") {
				ip := strings.Split(value, ":")[0]
				l.ClientIP = ip
				// TODO MOVE TO ENRICH
				l.TpSourceIP = &ip
				l.TpIps = append(l.TpIps, ip)
				l.ClientPort, _ = strconv.Atoi(strings.Split(value, ":")[1])
			}
		case "target":
			if value != "-" && strings.Contains(value, ":") {
				ip := strings.Split(value, ":")[0]
				l.TargetIP = &ip
				l.TpDestinationIP = &ip
				l.TpIps = append(l.TpIps, ip)
				l.TargetPort, _ = strconv.Atoi(strings.Split(value, ":")[1])
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
		case "elb_status_code":
			if value != "-" {
				esc, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("error parsing elb_status_code: %w", err)
				}
				l.ElbStatusCode = &esc
			}
		case "target_status_code":
			if value != "-" {
				tsc, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("error parsing target_status_code: %w", err)
				}
				l.TargetStatusCode = &tsc
			}
		case "received_bytes":
			if value != "-" {
				rb, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("error parsing received_bytes: %w", err)
				}
				l.ReceivedBytes = &rb
			}
		case "sent_bytes":
			if value != "-" {
				sb, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("error parsing sent_bytes: %w", err)
				}
				l.SentBytes = &sb
			}
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
			l.TpDomains = append(l.TpDomains, value)
		case "chosen_cert_arn":
			l.ChosenCertArn = value
		case "matched_rule_priority":
			if value != "-" {
				mrp, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("error parsing matched_rule_priority: %w", err)
				}
				l.MatchedRulePriority = mrp
			}
		case "request_creation_time":
			rct, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("error parsing request_creation_time: %w", err)
			}
			l.RequestCreationTime = rct
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
			l.ConnTraceID = value
		}
	}
	return nil
}
