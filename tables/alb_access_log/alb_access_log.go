package alb_access_log

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type AlbAccessLog struct {
	schema.CommonFields

	ActionsExecuted        []string  `json:"actions_executed,omitempty"`
	ChosenCertArn          string    `json:"chosen_cert_arn,omitempty"`
	Classification         *string   `json:"classification,omitempty"`
	ClassificationReason   *string   `json:"classification_reason,omitempty"`
	ClientIP               string    `json:"client_ip"`
	ClientPort             int       `json:"client_port"`
	ConnTraceID            *string   `json:"conn_trace_id,omitempty"`
	DomainName             string    `json:"domain_name,omitempty"`
	Elb                    string    `json:"elb"`
	ElbStatusCode          *int      `json:"elb_status_code,omitempty"`
	ErrorReason            *string   `json:"error_reason,omitempty"`
	MatchedRulePriority    int       `json:"matched_rule_priority,omitempty"`
	ReceivedBytes          *int64    `json:"received_bytes"`
	RedirectURL            *string   `json:"redirect_url,omitempty"`
	RequestHTTPVersion     string    `json:"request_http_version,omitempty"`
	RequestHTTPMethod      string    `json:"request_http_method,omitempty"`
	RequestUrl             string    `json:"request_url,omitempty"`
	RequestCreationTime    time.Time `json:"request_creation_time"`
	RequestProcessingTime  float64   `json:"request_processing_time,omitempty"`
	ResponseProcessingTime float64   `json:"response_processing_time,omitempty"`
	SentBytes              *int64    `json:"sent_bytes"`
	SslCipher              string    `json:"ssl_cipher,omitempty"`
	SslProtocol            string    `json:"ssl_protocol,omitempty"`
	TargetGroupArn         *string   `json:"target_group_arn,omitempty"`
	TargetIP               *string   `json:"target_ip,omitempty"`
	TargetList             *string   `json:"target_list,omitempty"`
	TargetPort             int       `json:"target_port,omitempty"`
	TargetProcessingTime   float64   `json:"target_processing_time,omitempty"`
	TargetStatusCode       *int      `json:"target_status_code,omitempty"`
	TargetStatusList       *string   `json:"target_status_list,omitempty"`
	Timestamp              time.Time `json:"timestamp"`
	TraceID                string    `json:"trace_id,omitempty"`
	Type                   string    `json:"type"`
	UserAgent              string    `json:"user_agent,omitempty"`
}

// InitialiseFromMap - initialise the struct from a map
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
		case "elb":
			l.Elb = value
		case "client_ip":
			l.ClientIP = value
		case "client_port":
			l.ClientPort, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing client_port: %w", err)
			}
		case "target":
				parts := strings.Split(value, ":")
				ip := parts[0]
				l.TargetIP = &ip
				l.TargetPort, err = strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("error parsing target_port: %w", err)
				}
		case "request_processing_time":
			if value != "-1" { // -1 if the load balancer can't dispatch the request to a target
				l.RequestProcessingTime, err = strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("error parsing request_processing_time: %w", err)
				}
			}
		case "target_processing_time":
			if value != "-1" { // -1 if the load balancer can't dispatch the request to a target
				l.TargetProcessingTime, err = strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("error parsing target_processing_time: %w", err)
				}
			}
		case "response_processing_time":
			if value != "-1" { // -1 if the load balancer doesn't receive a response from a target.
				l.ResponseProcessingTime, err = strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("error parsing response_processing_time: %w", err)
				}
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
		case "request_http_method":
			l.RequestHTTPMethod = value
		case "request_url":
			l.RequestUrl = value
		case "request_http_version":
			l.RequestHTTPVersion = value
		case "user_agent":
			l.UserAgent = value
		case "ssl_cipher":
			l.SslCipher = value
		case "ssl_protocol":
			l.SslProtocol = value
		case "target_group_arn":
			l.TargetGroupArn = &value
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
			l.ActionsExecuted = strings.Split(value, ",")
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
func (c *AlbAccessLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"actions_executed":         "The actions taken when processing the request, such as forwarding, redirecting, or fixed responses.",
		"chosen_cert_arn":          "The ARN of the certificate presented to the client during the TLS handshake.",
		"classification":           "The classification for desync mitigation, indicating compliance with RFC 7230.",
		"classification_reason":    "The classification reason code, describing why a request was classified as Acceptable, Ambiguous, or Severe.",
		"client_ip":                "The IP address of the client making the request.",
		"client_port":              "The source port used by the client for the connection.",
		"conn_trace_id":            "A unique identifier linking multiple log entries to the same connection.",
		"domain_name":              "The SNI domain provided by the client during the TLS handshake.",
		"elb":                      "The resource ID of the load balancer handling the request.",
		"elb_status_code":          "The HTTP status code returned by the load balancer.",
		"error_reason":             "The reason for request failure, including errors from AWS WAF or connection issues.",
		"matched_rule_priority":    "The priority of the rule that matched the request. Default rules are assigned 0.",
		"received_bytes":           "The number of bytes received from the client, including headers.",
		"redirect_url":             "The target URL for redirection if a redirect action was executed.",
		"request_http_method":      "The HTTP method used in the request.",
		"request_http_version":     "The HTTP version used in the request.",
		"request_url":              "The URL path requested by the client.",
		"request_creation_time":    "The timestamp when the load balancer received the request from the client.",
		"request_processing_time":  "The time elapsed from receiving the request to sending it to a target, in seconds.",
		"response_processing_time": "The time taken from receiving the response header to sending the response to the client.",
		"sent_bytes":               "The number of bytes sent to the client, including headers and body.",
		"ssl_cipher":               "The SSL cipher used for encrypting the connection.",
		"ssl_protocol":             "The SSL/TLS version used for the connection.",
		"target_group_arn":         "The ARN of the target group that handled the request.",
		"target_ip":                "The IP address of the target that processed the request.",
		"target_list":              "A space-delimited list of target IP addresses and ports.",
		"target_port":              "The port on the target that processed the request.",
		"target_processing_time":   "The time from when the request was sent to the target to when the response started.",
		"target_status_code":       "The HTTP status code returned by the target.",
		"target_status_list":       "A space-delimited list of status codes from targets that processed the request.",
		"timestamp":                "The timestamp when the load balancer generated a response.",
		"trace_id":                 "A unique identifier for tracing requests across AWS services.",
		"type":                     "The type of request (http, https, h2, grpcs, ws, wss).",
		"user_agent":               "The User-Agent string from the client making the request.",

		// Tailpipe-specific metadata fields
		"tp_ips":     "A list of IP addresses involved in the request, including the client IP and target IP.",
		"tp_domains": "A list of domains involved in the request, including the SNI domain from TLS connections.",
		"tp_akas":    "A list of AWS ARNs associated with the request, including the target group ARN.",
	}
}
