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
		case "request":
			l.Request = value
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
