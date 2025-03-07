package network_firewall_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type NetworkFirewallAlert struct {
	Action      string `json:"action,omitempty"`
	SignatureID int    `json:"signature_id,omitempty"`
	Rev         int    `json:"rev,omitempty"`
	Signature   string `json:"signature,omitempty"`
	Category    string `json:"category,omitempty"`
	Severity    int    `json:"severity,omitempty"`
}

type NetworkFirewallNetflow struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Bytes     int        `json:"bytes,omitempty"`
	Packets   int        `json:"packets,omitempty"`
	Age       int        `json:"age,omitempty"`
	MinTtl    int        `json:"min_ttl,omitempty"`
	MaxTtl    int        `json:"max_ttl,omitempty"`
}

type RevocationCheck struct {
	LeafCertFpr string `json:"leaf_cert_fpr,omitempty"`
	Status        string `json:"status,omitempty"`
	Action        string `json:"action,omitempty"`
}

type TLSError struct {
	ErrorMsg string `json:"error_message"`
}

type NetworkFirewallEvent struct {
	Timestamp       *time.Time              `json:"timestamp"`
	FlowID          int64                   `json:"flow_id"`
	EventType       string                  `json:"event_type"`
	SrcIP           string                  `json:"src_ip"`
	SrcPort         int                     `json:"src_port"`
	DestIP          string                  `json:"dest_ip"`
	DestPort        int                     `json:"dest_port"`
	Sni             string                  `json:"sni,omitempty"`
	Proto           string                  `json:"proto,omitempty"`
	AppProto        *string                 `json:"app_proto,omitempty"`
	Alert           *NetworkFirewallAlert   `json:"alert,omitempty"`
	Netflow         *NetworkFirewallNetflow `json:"netflow,omitempty"`
	TLSInspected    *bool                   `json:"tls_inspected,omitempty"`
	TLSError        *TLSError               `json:"tls_error,omitempty"`
	RevocationCheck *RevocationCheck        `json:"revocation_check,omitempty"`
}

type NetworkFirewallLog struct {
	schema.CommonFields

	FirewallName     string               `json:"firewall_name"`
	AvailabilityZone string               `json:"availability_zone"`
	EventTimestamp   *time.Time           `json:"event_timestamp"`
	Event            NetworkFirewallEvent `json:"event"`
}

func (n *NetworkFirewallLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"firewall_name":     "The name of the firewall associated with the log entry.",
		"availability_zone": "The Availability Zone of the firewall endpoint that generated the log entry.",
		"event_timestamp":   "The epoch timestamp (in seconds) when the log was created (UTC).",
		"event":             "Detailed event information in Suricata EVE JSON format, including human-readable timestamp, event type, network packet details, and alert details if applicable.",
		// Override table specific tp_* column descriptions
		"tp_index":     "The name of the firewall associated with the log entry.",
		"tp_ips":       "IP addresses extracted from the event (e.g., the source IP).",
		"tp_usernames": "Not applicable for firewall logs.",
	}
}
