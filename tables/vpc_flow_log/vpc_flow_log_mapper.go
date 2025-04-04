package vpc_flow_log

import (
	"context"
	"fmt"

	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type VPCFlowLogMapper struct {
}

func (c *VPCFlowLogMapper) Identifier() string {
	return "aws_waf_traffic_log_mapper"
}

func (c *VPCFlowLogMapper) Map(_ context.Context, a any, _ ...mappers.MapOption[*types.DynamicRow]) (*types.DynamicRow, error) {
	mappedData, ok := a.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("expected map[string]string, got %T", a)
	}

	row := &types.DynamicRow{}

	fields := getValidTokensAndColumnNames()

	rowColumnMap := make(map[string]string)

	for k, v := range mappedData {
		if columnName, ok := fields[k]; ok {
			if v != VpcFlowLogTableSkippedData && v != VpcFlowLogTableNoData && v != VpcFlowLogTableNilValue {
				rowColumnMap[columnName] = v
			}
		}
	}

	err := row.InitialiseFromMap(rowColumnMap)
	if err != nil {
		return nil, fmt.Errorf("error initialising dynamic row: %w", err)
	}

	return row, nil
}
