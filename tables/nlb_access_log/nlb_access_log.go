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

	Type                      string    `json:"type,omitempty"`
	Version                   string    `json:"version,omitempty"`
	Timestamp                 time.Time `json:"timestamp,omitempty"`
	Elb                       string    `json:"elb,omitempty"`
	Listener                  string    `json:"listener,omitempty"`
	ClientIP                  string    `json:"client_ip,omitempty"`
	ClientPort                int       `json:"client_port,omitempty"`
	DestinationIP             string    `json:"destination_ip,omitempty"`
	DestinationPort           int       `json:"destination_port,omitempty"`
	ConnectionTime            int       `json:"connection_time,omitempty"`
	TLSHandshakeTime          int       `json:"tls_handshake_time,omitempty"`
	ReceivedBytes             int64     `json:"received_bytes,omitempty"`
	SentBytes                 int64     `json:"sent_bytes,omitempty"`
	IncomingTLSAlert          string    `json:"incoming_tls_alert,omitempty"`
	ChosenCertArn             string    `json:"chosen_cert_arn,omitempty"`
	ChosenCertSerial          string    `json:"chosen_cert_serial,omitempty"`
	TLSCipher                 string    `json:"tls_cipher,omitempty"`
	TLSProtocolVersion        string    `json:"tls_protocol_version,omitempty"`
	TLSNamedGroup             string    `json:"tls_named_group,omitempty"`
	DomainName                string    `json:"domain_name,omitempty"`
	ALPNFEProtocol            string    `json:"alpn_fe_protocol,omitempty"`
	ALPNBEProtocol            string    `json:"alpn_be_protocol,omitempty"`
	ALPNClientPreferenceList  string    `json:"alpn_client_preference_list,omitempty"`
	TLSConnectionCreationTime time.Time `json:"tls_connection_creation_time,omitempty"`
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
			l.ALPNClientPreferenceList = value
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
