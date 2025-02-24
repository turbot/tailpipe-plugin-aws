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

	Type                      string    `json:"type"`
	Version                   string    `json:"version"`
	Time                      time.Time `json:"time"`
	Elb                       string    `json:"elb"`
	Listener                  string    `json:"listener"`
	ClientIP                  string    `json:"client_ip"`
	ClientPort                int       `json:"client_port"`
	DestinationIP             string    `json:"destination_ip"`
	DestinationPort           int       `json:"destination_port"`
	ConnectionTime            int       `json:"connection_time"`
	TLSHandshakeTime          int       `json:"tls_handshake_time"`
	ReceivedBytes             int64     `json:"received_bytes"`
	SentBytes                 int64     `json:"sent_bytes"`
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
	TLSConnectionCreationTime time.Time `json:"tls_connection_creation_time"`
}

// InitialiseFromMap - initialise the struct from a map
func (l *NlbAccessLog) InitialiseFromMap(m map[string]string) error {
	var err error
	for key, value := range m {
		if value == "-" {
			continue
		}
		switch key {
		case "time":
			l.Time, err = time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("error parsing time: %w", err)
			}
		case "type":
			l.Type = value
		case "version":
			l.Version = value
		case "elb":
			l.Elb = value
		case "listener":
			l.Listener = value
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
		case "domain_name":
			l.DomainName = value
		case "tls_connection_creation_time":
			l.TLSConnectionCreationTime, err = time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("error parsing tls_connection_creation_time: %w", err)
			}
		}
	}
	return nil
}
