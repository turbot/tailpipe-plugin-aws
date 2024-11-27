package mappers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type WafMapper struct {
}

// NewWafMapper creates a new WafMapper
func NewWafMapper() table.Mapper[*rows.WafTrafficLog] {
	return &WafMapper{}
}

func (c *WafMapper) Identifier() string {
	return "waf_mapper"
}

func (c *WafMapper) Map(_ context.Context, a any) (*rows.WafTrafficLog, error) {
	var jsonBytes []byte

	switch v := a.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		return nil, fmt.Errorf("expected byte[] or string, got %T", a)
	}

	// decode JSON into WafTrafficLog
	var log rows.WafTrafficLog
	err := json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w; partial log: %+v", err, log)
	}

	return &log, nil
}
