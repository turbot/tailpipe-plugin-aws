package clb_access_log

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type ClbAccessLog struct {
	schema.CommonFields

	BackendIP              *string   `json:"backend_ip,omitempty"`
	BackendPort            int       `json:"backend_port,omitempty"`
	BackendProcessingTime  float64   `json:"backend_processing_time,omitempty"`
	BackendStatusCode      *int      `json:"backend_status_code,omitempty"`
	ClientIP               string    `json:"client_ip"`
	ClientPort             int       `json:"client_port"`
	Elb                    string    `json:"elb"`
	ElbStatusCode          *int      `json:"elb_status_code,omitempty"`
	ReceivedBytes          *int64    `json:"received_bytes,omitempty"`
	RequestHTTPVersion     string    `json:"request_http_version,omitempty"`
	RequestHTTPMethod      string    `json:"request_http_method,omitempty"`
	RequestUrl             string    `json:"request_url,omitempty"`
	RequestProcessingTime  float64   `json:"request_processing_time,omitempty"`
	ResponseProcessingTime float64   `json:"response_processing_time,omitempty"`
	SentBytes              *int64    `json:"sent_bytes,omitempty"`
	SslCipher              string    `json:"ssl_cipher,omitempty"`
	SslProtocol            string    `json:"ssl_protocol,omitempty"`
	Timestamp              time.Time `json:"timestamp"`
	UserAgent              string    `json:"user_agent,omitempty"`
}

// InitialiseFromMap - initialise the struct from a map
func (l *ClbAccessLog) InitialiseFromMap(m map[string]string) error {
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
		case "elb":
			l.Elb = value
		case "client_ip":
			l.ClientIP = value
		case "client_port":
			l.ClientPort, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing client_port: %w", err)
			}
		case "backend":
				parts := strings.Split(value, ":")
				ip := parts[0]
				l.BackendIP = &ip
				l.BackendPort, err = strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("error parsing backend_port: %w", err)
				}
		case "request_processing_time":
			l.RequestProcessingTime, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("error parsing request_processing_time: %w", err)
			}
		case "backend_processing_time":
			l.BackendProcessingTime, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("error parsing backend_processing_time: %w", err)
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
		case "backend_status_code":
			bsc, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing backend_status_code: %w", err)
			}
			l.BackendStatusCode = &bsc
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
		}
	}

	return nil
}

func (c *ClbAccessLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"backend_ip":               "The IP address of the registered instance that processed the request.",
		"backend_port":             "The port on the registered instance that processed the request.",
		"backend_processing_time":  "The time elapsed from the load balancer sending the request to the registered instance until the instance starts sending response headers.",
		"backend_status_code":      "The HTTP status code returned by the registered instance.",
		"client_ip":                "The IP address of the requesting client.",
		"client_port":              "The source port used by the client for the connection.",
		"elb":                      "The name of the load balancer.",
		"elb_status_code":          "The HTTP status code returned by the load balancer.",
		"received_bytes":           "The size of the request in bytes received from the client.",
		"request_http_version":     "The HTTP version of the request.",
		"request_http_method":      "The HTTP method of the request.",
		"request_url":              "The URL of the request.",
		"request_processing_time":  "The time elapsed from receiving the request to sending it to a registered instance, in seconds.",
		"response_processing_time": "The time elapsed from the load balancer receiving the response headers to sending the response to the client.",
		"sent_bytes":               "The size of the response in bytes sent to the client.",
		"ssl_cipher":               "The SSL cipher used for encrypting the connection.",
		"ssl_protocol":             "The SSL/TLS version used for the connection.",
		"timestamp":                "The time when the load balancer received the request from the client, in ISO 8601 format.",
		"user_agent":               "A User-Agent string that identifies the client that originated the request.",

		// Tailpipe-specific metadata fields
		"tp_index":          "The name of the load balancer.",
		"tp_source_ip":      "The IP address of the requesting client.",
		"tp_ips":            "The IP addresses of the requesting client and the registered instance that processed the request.",
		"tp_destination_ip": "The IP address of the registered instance that processed the request.",
	}
}
