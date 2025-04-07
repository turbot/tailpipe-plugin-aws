package securityhub_finding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

type SecurityHubFindingMapper struct {
}

func (m *SecurityHubFindingMapper) Identifier() string {
	return "security_hub_finding_mapper"
}

func (m *SecurityHubFindingMapper) Map(_ context.Context, a any, _ ...mappers.MapOption[*SecurityHubFinding]) (*SecurityHubFinding, error) {
	var b SecurityHubFinding

	switch data := a.(type) {
	case []byte:
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("error unmarshalling row data: %w", err)
		}
	case string:
		if err := json.Unmarshal([]byte(data), &b); err != nil {
			return nil, fmt.Errorf("error unmarshalling row data: %w", err)
		}
	case SecurityHubFinding:
		b = data
	default:
		return nil, fmt.Errorf("expected byte[], string or SecurityHubFinding, got %T", a)
	}

	return &b, nil

}