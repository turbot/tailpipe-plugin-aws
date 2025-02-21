package clb_access_log

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type ClbAccessLogBatch struct {
	Records []ClbAccessLog `json:"Records"`
}

type ClbAccessLog struct {
	schema.CommonFields

	ClientIP               string    `json:"client_ip,omitempty"`
	ClientPort             int       `json:"client_port,omitempty"`
	Elb                    string    `json:"elb,omitempty"`
	ElbStatusCode          *int      `json:"elb_status_code,omitempty"`
	BackendStatusCode      *int      `json:"backend_status_code,omitempty"`
	ReceivedBytes          *int64    `json:"received_bytes,omitempty"`
	SentBytes              *int64    `json:"sent_bytes,omitempty"`
	BackendIP              *string   `json:"backend_ip,omitempty"`
	Request                string    `json:"request,omitempty"`
	RequestProcessingTime  float64   `json:"request_processing_time,omitempty"`
	BackendProcessingTime  float64   `json:"backend_processing_time,omitempty"`
	ResponseProcessingTime float64   `json:"response_processing_time,omitempty"`
	SslCipher              string    `json:"ssl_cipher,omitempty"`
	SslProtocol            string    `json:"ssl_protocol,omitempty"`
	Timestamp              time.Time `json:"timestamp,omitempty"`
	UserAgent              string    `json:"user_agent,omitempty"`
}

// InitialiseFromMap - initialise the struct from a map
func (l *ClbAccessLog) InitialiseFromMap(m map[string]string) error {
	var err error
	var method, path, httpVersion string

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
		case "client":
			if strings.Contains(value, ":") {
				parts := strings.Split(value, ":")
				l.ClientIP = parts[0]
				l.ClientPort, err = strconv.Atoi(parts[1])
				if err != nil {
					return fmt.Errorf("error parsing client_port: %w", err)
				}
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
		case "method":
			method = value
		case "path":
			path = value
		case "http_version":
			httpVersion = value
		case "user_agent":
			l.UserAgent = value
		case "ssl_cipher":
			l.SslCipher = value
		case "ssl_protocol":
			l.SslProtocol = value
		}
	}

	// Construct request string in the correct order
	if method != "" && path != "" && httpVersion != "" {
		l.Request = fmt.Sprintf("%s %s %s", method, path, httpVersion)
	}

	return nil
}

func (c *ClbAccessLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"timestamp":                 "The time when the load balancer received the request from the client, in ISO 8601 format.",
		"elb":                       "The name of the load balancer.",
		"client_ip":                 "The IP address of the requesting client.",
		"client_port":               "The source port used by the client for the connection.",
		"backend_ip":                "The IP address of the registered instance that processed the request.",
		"request":                   "The full request line from the client, including method, protocol, and URI.",
		"request_processing_time":   "The time elapsed from receiving the request to sending it to a registered instance, in seconds.",
		"backend_processing_time":   "The time elapsed from the load balancer sending the request to the registered instance until the instance starts sending response headers.",
		"response_processing_time":  "The time elapsed from the load balancer receiving the response headers to sending the response to the client.",
		"elb_status_code":           "The HTTP status code returned by the load balancer.",
		"backend_status_code":       "The HTTP status code returned by the registered instance.",
		"received_bytes":            "The size of the request in bytes received from the client.",
		"sent_bytes":                "The size of the response in bytes sent to the client.",
		"user_agent":                "A User-Agent string that identifies the client that originated the request.",
		"ssl_cipher":                "The SSL cipher used for encrypting the connection.",
		"ssl_protocol":              "The SSL/TLS version used for the connection.",
	}
}
