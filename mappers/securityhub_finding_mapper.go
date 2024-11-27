package mappers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type SecurityHubFindingsMapper struct {
}

func NewSecurityHubFindingsMapper() table.Mapper[*rows.SecurityHubFindingLog] {
	res := &SecurityHubFindingsMapper{}

	return res
}

func (m *SecurityHubFindingsMapper) Identifier() string {
	return "security_hub_finding_mapper"
}

func (m *SecurityHubFindingsMapper) Map(_ context.Context, a any) (*rows.SecurityHubFindingLog, error) {
	var b rows.SecurityHubFindingLog

	switch data := a.(type) {
	case []byte:
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("error unmarshalling row data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &b); err != nil {
			return nil, fmt.Errorf("error unmarshalling row data: %w", err)
		}
	case rows.SecurityHubFindingLog:
		b = data
	default:
		return nil, fmt.Errorf("expected byte[], string or rows.SecurityHubFindingLog, got %T", a)
	}

	return &b, nil

}
