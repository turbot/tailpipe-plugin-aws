package alb_connection_log

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

// AlbConnectionLog represents a connection log entry from an AWS ALB.
type AlbConnectionLog struct {
	schema.CommonFields

	ClientIP                   string    `json:"client_ip"`
	ClientPort                 int       `json:"client_port"`
	ConnTraceID                *string   `json:"conn_trace_id"`
	LeafClientCertSerialNumber string    `json:"leaf_client_cert_serial_number,omitempty"`
	LeafClientCertSubject      string    `json:"leaf_client_cert_subject,omitempty"`
	LeafClientCertValidity     string    `json:"leaf_client_cert_validity,omitempty"`
	ListenerPort               int       `json:"listener_port,omitempty"`
	TLSCipher                  string    `json:"tls_cipher,omitempty"`
	TLSHandshakeLatency        float64   `json:"tls_handshake_latency,omitempty"`
	TLSProtocol                string    `json:"tls_protocol,omitempty"`
	TLSVerifyStatus            string    `json:"tls_verify_status,omitempty"`
	Timestamp                  time.Time `json:"timestamp"`
}

// InitialiseFromMap populates the AlbConnectionLog fields from a map.
// It expects keys that match the field names defined in the connection log format.
func (c *AlbConnectionLog) InitialiseFromMap(m map[string]string) error {
	var err error
	for key, value := range m {
		// Ignore placeholder values
		if value == "-" {
			continue
		}
		switch key {
		case "timestamp":
			ts, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("error parsing timestamp: %w", err)
			}
			c.Timestamp = ts
		case "client_ip":
			c.ClientIP = value
		case "client_port":
			c.ClientPort, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing client_port: %w", err)
			}
		case "listener_port":
			c.ListenerPort, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing listener_port: %w", err)
			}
		case "tls_protocol":
			c.TLSProtocol = value
		case "tls_cipher":
			c.TLSCipher = value
		case "tls_handshake_latency":
			// Remove any potential commas or extra spaces
			cleanValue := strings.TrimSpace(value)
			c.TLSHandshakeLatency, err = strconv.ParseFloat(cleanValue, 64)
			if err != nil {
				return fmt.Errorf("error parsing tls_handshake_latency: %w", err)
			}
		case "leaf_client_cert_subject":
			c.LeafClientCertSubject = value
		case "leaf_client_cert_validity":
			c.LeafClientCertValidity = value
		case "leaf_client_cert_serial_number":
			c.LeafClientCertSerialNumber = value
		case "tls_verify_status":
			c.TLSVerifyStatus = value
		case "conn_trace_id":
			c.ConnTraceID = &value
		}
	}
	return nil
}

// GetColumnDescriptions returns a mapping of connection log field names to their descriptions.
func (c *AlbConnectionLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"client_ip":                      "The IP address of the requesting client.",
		"client_port":                    "The port of the requesting client.",
		"conn_trace_id":                  "A unique identifier linking connection logs to subsequent access logs for the same connection.",
		"leaf_client_cert_serial_number": "The serial number of the leaf client certificate.",
		"leaf_client_cert_subject":       "The subject name of the leaf client certificate (if applicable).",
		"leaf_client_cert_validity":      "The validity period (not-before and not-after in ISO 8601 format) of the leaf client certificate.",
		"listener_port":                  "The port of the load balancer listener receiving the client request.",
		"timestamp":                      "The time when the load balancer established or failed to establish a connection (ISO 8601 format).",
		"tls_cipher":                     "The cipher used during the TLS handshake. For HTTPS listeners, the cipher used during the handshake.",
		"tls_handshake_latency":          "The total time in seconds, with millisecond precision, elapsed during the TLS handshake.",
		"tls_protocol":                   "The SSL/TLS protocol used during the handshake. For HTTPS listeners, the SSL/TLS protocol used during the handshake.",
		"tls_verify_status":              "The status of the TLS verification ('Success' if established, otherwise 'Failed:$error_code').",

		// Tailpipe-specific metadata fields
		"tp_ips":   "The IP addresses involved in requesting client.",
	}
}
