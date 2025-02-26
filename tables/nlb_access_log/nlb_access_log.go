package nlb_access_log

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type NlbAccessLog struct {
	schema.CommonFields

	ALPNBEProtocol            string    `json:"alpn_be_protocol,omitempty"`
	ALPNClientPreferenceList  []string  `json:"alpn_client_preference_list,omitempty"`
	ALPNFEProtocol            string    `json:"alpn_fe_protocol,omitempty"`
	ChosenCertArn             string    `json:"chosen_cert_arn,omitempty"`
	ChosenCertSerial          string    `json:"chosen_cert_serial,omitempty"`
	ClientIP                  string    `json:"client_ip,omitempty"`
	ClientPort                int       `json:"client_port,omitempty"`
	ConnectionTime            int       `json:"connection_time,omitempty"`
	DestinationIP             string    `json:"destination_ip,omitempty"`
	DestinationPort           int       `json:"destination_port,omitempty"`
	DomainName                string    `json:"domain_name,omitempty"`
	Elb                       string    `json:"elb"`
	IncomingTLSAlert          string    `json:"incoming_tls_alert,omitempty"`
	Listener                  string    `json:"listener,omitempty"`
	ReceivedBytes             int64     `json:"received_bytes,omitempty"`
	SentBytes                 int64     `json:"sent_bytes,omitempty"`
	TLSCipher                 string    `json:"tls_cipher,omitempty"`
	TLSConnectionCreationTime time.Time `json:"tls_connection_creation_time"`
	TLSHandshakeTime          int       `json:"tls_handshake_time,omitempty"`
	TLSNamedGroup             string    `json:"tls_named_group,omitempty"`
	TLSProtocolVersion        string    `json:"tls_protocol_version,omitempty"`
	Timestamp                 time.Time `json:"timestamp"`
	Type                      string    `json:"type"`
	Version                   string    `json:"version"`
}

// InitialiseFromMap - initialise the struct from a map
func (l *NlbAccessLog) InitialiseFromMap(m map[string]string) error {
	var err error
	ISO8601 := "2006-01-02T15:04:05"
	for key, value := range m {
		if value == "-" {
			continue
		}
		switch key {
		case "timestamp":
			ts, err := time.Parse(ISO8601, value)
			if err != nil {
				return fmt.Errorf("error parsing timestamp: %w", err)
			}
			l.Timestamp = ts
		case "type":
			l.Type = value
		case "version":
			l.Version = value
		case "elb":
			l.Elb = value
		case "listener":
			l.Listener = value
		case "incoming_tls_alert":
			l.IncomingTLSAlert = value
		case "chosen_cert_serial":
			l.ChosenCertSerial = value
		case "client":
			parts := strings.Split(value, ":")
			l.ClientIP = parts[0]
			l.ClientPort, err = strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("error parsing client_port: %w", err)
			}
		case "destination":
			parts := strings.Split(value, ":")
			l.DestinationIP = parts[0]
			l.DestinationPort, err = strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("error parsing destination_port: %w", err)
			}
		case "connection_time":
			l.ConnectionTime, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing connection_time: %w", err)
			}
		case "tls_handshake_time":
			l.TLSHandshakeTime, err = strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing tls_handshake_time: %w", err)
			}
		case "received_bytes":
			l.ReceivedBytes, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing received_bytes: %w", err)
			}
		case "sent_bytes":
			l.SentBytes, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing sent_bytes: %w", err)
			}
		case "chosen_cert_arn":
			l.ChosenCertArn = value
		case "tls_cipher":
			l.TLSCipher = value
		case "tls_protocol_version":
			l.TLSProtocolVersion = value
		case "tls_named_group":
			l.TLSNamedGroup = value
		case "domain_name":
			l.DomainName = value
		case "alpn_fe_protocol":
			l.ALPNFEProtocol = value
		case "alpn_be_protocol":
			l.ALPNBEProtocol = value
		case "alpn_client_preference_list":
			l.ALPNClientPreferenceList = strings.Split(value, ",")
		case "tls_connection_creation_time":
			ts, err := time.Parse(ISO8601, value)
			if err != nil {
				return fmt.Errorf("error parsing timestamp: %w", err)
			}
			l.TLSConnectionCreationTime = ts
		}
	}
	return nil
}

func (l *NlbAccessLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"alpn_be_protocol":             "The application protocol negotiated with the target. Possible values: h2, http/1.1, http/1.0.",
		"alpn_client_preference_list":  "The value of the application_layer_protocol_negotiation extension in the client hello message, URL-encoded.",
		"alpn_fe_protocol":             "The application protocol negotiated with the client. Possible values: h2, http/1.1, http/1.0.",
		"chosen_cert_arn":              "The ARN of the certificate served to the client.",
		"chosen_cert_serial":           "Reserved for future use. This value is always set to -.",
		"client_ip":                    "The IP address of the client initiating the connection.",
		"client_port":                  "The port number used by the client to establish the connection.",
		"connection_time":              "The total time for the connection to complete, from start to closure, in milliseconds.",
		"destination_ip":               "The IP address of the destination, which could be the listener or VPC endpoint.",
		"destination_port":             "The port on the destination receiving the connection.",
		"domain_name":                  "The value of the server_name extension in the client hello message, URL-encoded.",
		"elb":                          "The resource ID of the load balancer.",
		"incoming_tls_alert":           "The integer value of TLS alerts received by the load balancer from the client, if present.",
		"listener":                     "The resource ID of the TLS listener for the connection.",
		"received_bytes":               "The count of bytes received by the load balancer from the client, after decryption.",
		"sent_bytes":                   "The count of bytes sent by the load balancer to the client, before encryption.",
		"timestamp":                    "The time recorded at the end of the TLS connection, in ISO 8601 format.",
		"tls_cipher":                   "The cipher suite negotiated with the client, in OpenSSL format.",
		"tls_connection_creation_time": "The time recorded at the beginning of the TLS connection, in ISO 8601 format.",
		"tls_handshake_time":           "The total time for the TLS handshake to complete after the TCP connection is established, including client-side delays, in milliseconds.",
		"tls_named_group":              "Reserved for future use. This value is always set to -.",
		"tls_protocol_version":         "The TLS protocol negotiated with the client. Possible values: tlsv10, tlsv11, tlsv12, tlsv13.",
		"type":                         "The type of listener. The supported value is tls.",
		"version":                      "The version of the log entry. The current version is 2.0.",

		// Tailpipe-specific metadata fields
		"tp_akas":      "The ARN of the certificate served to the client.",
		"tp_index":     "The resource ID of the load balancer handling the request.",
		"tp_ips":       "All IP addresses associated with the request, including the client IP and destination IP.",
		"tp_source_ip": "The IP address of the client initiating the connection.",
	}
}
