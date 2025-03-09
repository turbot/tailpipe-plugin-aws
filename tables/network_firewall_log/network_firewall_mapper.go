package network_firewall_log

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type NetworkFirewallMapper struct {
}

func (n *NetworkFirewallMapper) Identifier() string {
	return "network_firewall_log_mapper"
}

func (n *NetworkFirewallMapper) Map(_ context.Context, a any, _ ...table.MapOption[*NetworkFirewallLog]) (*NetworkFirewallLog, error) {
	var jsonBytes []byte

	switch v := a.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		return nil, fmt.Errorf("expected byte[] or string, got %T", a)
	}

	var log NetworkFirewallLog
	err := unmarshalNetworkFirewallLog(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w; partial log: %+v", err, log)
	}

	return &log, nil
}

func unmarshalNetworkFirewallLog(data []byte, log *NetworkFirewallLog) error {
	var temp struct {
		FirewallName     string `json:"firewall_name"`
		AvailabilityZone string `json:"availability_zone"`
		EventTimestamp   string `json:"event_timestamp"`
		Event            struct {
			Timestamp    string `json:"timestamp"`
			FlowID       int64  `json:"flow_id"`
			EventType    string `json:"event_type"`
			SrcIP        string `json:"src_ip"`
			SrcPort      int    `json:"src_port"`
			DestIP       string `json:"dest_ip"`
			DestPort     int    `json:"dest_port"`
			SNI          string `json:"sni"`
			TLSInspected *bool  `json:"tls_inspected,omitempty"`
			Proto        string `json:"proto"`
			AppProto     string `json:"app_proto"`
			TLSError     struct {
				ErrorMessage string `json:"error_message"`
			} `json:"tls_error,omitempty"`
			NetFlow struct {
				Packets int64  `json:"pkts"`
				Bytes   int64  `json:"bytes"`
				Start   string `json:"start"`
				End     string `json:"end"`
				Age     int64  `json:"age"`
				MinTTL  int64  `json:"min_ttl"`
				MaxTTL  int64  `json:"max_ttl"`
			} `json:"netflow"`
			Alert struct {
				Action      string `json:"action"`
				SignatureID int    `json:"signature_id"`
				Rev         int    `json:"rev"`
				Signature   string `json:"signature"`
				Category    string `json:"category"`
				Severity    int    `json:"severity"`
			} `json:"alert,omitempty"`
			RevocationCheck struct {
				LeafCertFpr string `json:"leaf_cert_fpr"`
				Status      string `json:"status"`
				Action      string `json:"action"`
			} `json:"revocation_check,omitempty"`
		} `json:"event"`
	}

	// Unmarshal JSON directly into the NetworkFirewallLog struct.
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	const layout = "2006-01-02T15:04:05.999999-0700"
	// Assign values from temp struct
	log.FirewallName = temp.FirewallName
	log.AvailabilityZone = temp.AvailabilityZone

	// Parse and assign EventTimestamp
	if temp.EventTimestamp != "" {
		// For Unix timestamps, we must convert to integer first, then to time
		unixSeconds, err := strconv.ParseInt(temp.EventTimestamp, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse event_timestamp as unix timestamp: %w", err)
		}
		timestamp := time.Unix(unixSeconds, 0)
		log.EventTimestamp = &timestamp
	}

	// Populate Event fields
	if temp.Event.Timestamp != "" {
		// Try parsing as Unix timestamp first
		timestamp, err := time.Parse(layout, temp.Event.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to parse event timestamp: %w", err)
		}
		log.Event.Timestamp = &timestamp
	}

	log.Event.FlowID = temp.Event.FlowID
	log.Event.EventType = temp.Event.EventType
	log.Event.SrcIP = temp.Event.SrcIP
	log.Event.SrcPort = temp.Event.SrcPort
	log.Event.DestIP = temp.Event.DestIP
	log.Event.DestPort = temp.Event.DestPort
	log.Event.Sni = temp.Event.SNI
	log.Event.Proto = temp.Event.Proto

	if temp.Event.AppProto != "" {
		log.Event.AppProto = &temp.Event.AppProto
	}

	log.Event.TLSInspected = temp.Event.TLSInspected

	// Handle TLSError if present
	if temp.Event.TLSError.ErrorMessage != "" {
		log.Event.TLSError = &TLSError{
			ErrorMsg: temp.Event.TLSError.ErrorMessage,
		}
	}

	// Handle Netflow
	if temp.Event.NetFlow.Packets > 0 || temp.Event.NetFlow.Bytes > 0 {
		netflow := &NetworkFirewallNetflow{
			Packets: int(temp.Event.NetFlow.Packets),
			Bytes:   int(temp.Event.NetFlow.Bytes),
			Age:     int(temp.Event.NetFlow.Age),
			MinTtl:  int(temp.Event.NetFlow.MinTTL),
			MaxTtl:  int(temp.Event.NetFlow.MaxTTL),
		}

		// Parse start and end times
		if temp.Event.NetFlow.Start != "" {
			startTime, err := time.Parse(layout, temp.Event.NetFlow.Start)
			if err != nil {
				return fmt.Errorf("failed to parse netflow start time: %w", err)
			}
			netflow.StartTime = &startTime
		}

		if temp.Event.NetFlow.End != "" {
			endTime, err := time.Parse(layout, temp.Event.NetFlow.End)
			if err != nil {
				return fmt.Errorf("failed to parse netflow end time: %w", err)
			}
			netflow.EndTime = &endTime
		}

		log.Event.Netflow = netflow
	}

	// Handle Alert if present
	if temp.Event.Alert.Action != "" || temp.Event.Alert.SignatureID != 0 {
		log.Event.Alert = &NetworkFirewallAlert{
			Action:      temp.Event.Alert.Action,
			SignatureID: temp.Event.Alert.SignatureID,
			Rev:         temp.Event.Alert.Rev,
			Signature:   temp.Event.Alert.Signature,
			Category:    temp.Event.Alert.Category,
			Severity:    temp.Event.Alert.Severity,
		}
	}

	// Handle RevocationCheck if present
	if temp.Event.RevocationCheck.LeafCertFpr != "" || temp.Event.RevocationCheck.Status != "" {
		log.Event.RevocationCheck = &RevocationCheck{
			LeafCertFpr: temp.Event.RevocationCheck.LeafCertFpr,
			Status:      temp.Event.RevocationCheck.Status,
			Action:      temp.Event.RevocationCheck.Action,
		}
	}
	return nil
}

