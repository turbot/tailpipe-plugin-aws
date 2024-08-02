package aws_types

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type AwsElbAccessLog struct {
	enrichment.CommonFields

	Type                   string    `json:"type"`
	Timestamp              time.Time `json:"timestamp"`
	Elb                    string    `json:"elb"`
	ClientIP               string    `json:"client_ip"`
	ClientPort             int       `json:"client_port"`
	TargetIP               string    `json:"target_ip"`
	TargetPort             int       `json:"target_port"`
	RequestProcessingTime  float64   `json:"request_processing_time"`
	TargetProcessingTime   float64   `json:"target_processing_time"`
	ResponseProcessingTime float64   `json:"response_processing_time"`
	ElbStatusCode          int       `json:"elb_status_code"`
	TargetStatusCode       int       `json:"target_status_code"`
	ReceivedBytes          int64     `json:"received_bytes"`
	SentBytes              int64     `json:"sent_bytes"`
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
}
