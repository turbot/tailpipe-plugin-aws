package vpc_flow_log

import (
	"context"
	"fmt"

	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

// VPCFlowLogCloudWatchMapper is a custom mapper for VPC flow logs from CloudWatch sources.
// It extracts the message field from the complete CloudWatch event JSON and processes it
// using the existing regex mapper for VPC flow logs.
type VPCFlowLogCloudWatchMapper struct {
	// The underlying regex mapper that processes the VPC flow log message
	mapper mappers.Mapper[*types.DynamicRow]
}

// NewVPCFlowLogCloudWatchMapper creates a new CloudWatch mapper for VPC flow logs
func NewVPCFlowLogCloudWatchMapper(format *VPCFlowLogTableFormat) (*VPCFlowLogCloudWatchMapper, error) {
	// Get the regex mapper from the format
	mapper, err := format.GetMapper()
	if err != nil {
		return nil, fmt.Errorf("failed to create regex mapper: %w", err)
	}

	return &VPCFlowLogCloudWatchMapper{
		mapper: mapper,
	}, nil
}

// Identifier returns the mapper identifier
func (m *VPCFlowLogCloudWatchMapper) Identifier() string {
	return "vpc_flow_log_cloudwatch_mapper"
}

// Map processes CloudWatch event data by extracting the message field and processing it with the regex mapper
func (m *VPCFlowLogCloudWatchMapper) Map(ctx context.Context, a any, opts ...mappers.MapOption[*types.DynamicRow]) (*types.DynamicRow, error) {
	var input []byte

	// Handle different input types
	switch v := a.(type) {
	case []byte:
		input = v
	case string:
		input = []byte(v)
	case *string:
		if v != nil {
			input = []byte(*v)
		} else {
			return nil, fmt.Errorf("nil string input")
		}
	case cwTypes.FilteredLogEvent:
		input = []byte(*v.Message)
	default:
		return nil, fmt.Errorf("expected byte[], string, or *string, got %T", a)
	}

	// Use the regex mapper to process the message
	return m.mapper.Map(ctx, string(input), opts...)
}
