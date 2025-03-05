package network_firewall_log

import (
	"context"
	"encoding/json"
	"fmt"

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
	
	// Unmarshal JSON directly into the NetworkFirewallLog struct.
	if err := json.Unmarshal(data, log); err != nil {
		return err
	}
	return nil
}