package mappers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// WafMapper is a mapper that receives CloudTrailBatch objects and extracts CloudTrailLog records from them
type WafMapper struct {
}

// NewWafMapper creates a new WafMapper
func NewWafMapper() table.Mapper[rows.WafTrafficLog] {
	return &WafMapper{}
}

func (c *WafMapper) Identifier() string {
	return "waf_mapper"
}

func (c *WafMapper) Map(_ context.Context, a any) ([]rows.WafTrafficLog, error) {
	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	// decode json ito CloudTrailBatch
	var log rows.WafTrafficLog
	err := json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return []rows.WafTrafficLog{log}, nil
}
