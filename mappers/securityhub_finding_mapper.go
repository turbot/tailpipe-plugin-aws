package mappers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/tailpipe-plugin-aws/rows"
)

type SecurityHubFindingMapper struct {
}

func (m *SecurityHubFindingMapper) Identifier() string {
	return "security_hub_finding_mapper"
}

func (m *SecurityHubFindingMapper) Map(_ context.Context, a any) (*rows.SecurityHubFinding, error) {
	var b rows.SecurityHubFinding

	switch data := a.(type) {
	case []byte:
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("error unmarshalling row data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &b); err != nil {
			return nil, fmt.Errorf("error unmarshalling row data: %w", err)
		}
	case rows.SecurityHubFinding:
		b = data
	default:
		return nil, fmt.Errorf("expected byte[], string or rows.SecurityHubFinding, got %T", a)
	}

	return &b, nil

}
