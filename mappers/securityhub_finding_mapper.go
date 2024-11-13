package mappers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type SecurityHubFindingsMapper struct {
}

func NewSecurityHubFindingsMapper() table.Mapper[*rows.SecurityHubFindingLog] {
	res := &SecurityHubFindingsMapper{}

	return res
}

func (c *SecurityHubFindingsMapper) Identifier() string {
	return "securityhub_finding_mapper"
}

func (m *SecurityHubFindingsMapper) Map(_ context.Context, a any) ([]*rows.SecurityHubFindingLog, error) {
	var b rows.SecurityHubFindingLog

	// If `a` is a []uint8, unmarshal it into the expected struct
	if data, ok := a.([]uint8); ok {
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("error unmarshalling row data: %w", err)
		}
	} else {
		// Type assertion as usual if `a` is already of type rows.SecurityHubFindingLog
		var ok bool
		b, ok = a.(rows.SecurityHubFindingLog)
		if !ok {
			return nil, fmt.Errorf("expected rows.SecurityHubFindingLog, got %T", a)
		}
	}

	// Proceed with existing logic
	s, err := json.Marshal(b.DetailFindings)
	if err != nil {
		return nil, fmt.Errorf("error marshalling row data: %w", err)
	}

	dataJsonString := types.JSONString(s)
	b.DetailFindings = &dataJsonString

	return []*rows.SecurityHubFindingLog{&b}, nil

}
