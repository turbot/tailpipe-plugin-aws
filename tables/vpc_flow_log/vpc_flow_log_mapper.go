package vpc_flow_log

import (
	"context"
	"fmt"

	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// VpcFlowLogMapper is a mapper that receives string objects and extracts VpcFlowLog record
type VpcFlowLogMapper struct{}

func (c *VpcFlowLogMapper) Identifier() string {
	return "vpc_flow_log_mapper"
}

// Map casts the data item as an string and returns the VpcFlowLog records
func (c *VpcFlowLogMapper) Map(_ context.Context, a any, _ ...table.MapOption[*VpcFlowLog]) (*VpcFlowLog, error) {
	flowLog, ok := a.(*VpcFlowLog)
	if !ok {
		return nil, fmt.Errorf("expected *VpcFlowLog, got %T", a)
	}

	return flowLog, nil

}