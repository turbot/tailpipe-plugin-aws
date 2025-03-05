package network_firewall_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type NetworkFirewallAlert struct {
	Action      string `json:"action"`
	SignatureID int    `json:"signature_id"`
	Rev         int    `json:"rev"`
	Signature   string `json:"signature"`
	Category    string `json:"category"`
	Severity    int    `json:"severity"`
}

type NetworkFirewallEvent struct {
	Timestamp       *time.Time            `json:"timestamp"`
	FlowID          int64                 `json:"flow_id"`
	EventType       string                `json:"event_type"`
	SrcIP           string                `json:"src_ip"`
	SrcPort         int                   `json:"src_port"`
	DestIP          string                `json:"dest_ip"`
	DestPort        int                   `json:"dest_port"`
	Proto           string                `json:"proto"`
	Alert           *NetworkFirewallAlert `json:"alert,omitempty"`
	TLSInspected    *bool                 `json:"tls_inspected,omitempty"`
	TLSError        *string               `json:"tls_error,omitempty"`
	RevocationCheck *string               `json:"revocation_check,omitempty"`
	Bytes           *int64                `json:"bytes,omitempty"`
	Packets         *int64                `json:"packets,omitempty"`
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
