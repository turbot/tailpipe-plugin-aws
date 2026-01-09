package vpc_flow_log

import (
	"context"
	"errors"
	"fmt"

	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

// ErrNoDataRecord is returned when a NODATA record is encountered.
// NODATA records indicate no traffic was recorded during the aggregation interval
// and don't have timestamp data by design - they should be skipped.
var ErrNoDataRecord = errors.New("skipping NODATA record: no traffic was recorded during aggregation interval")

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

// Map processes CloudWatch event data by extracting the message field and processing it with the regex mapper.
// Returns ErrNoDataRecord for NODATA records which should be skipped (they have no traffic data by design).
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
	row, err := m.mapper.Map(ctx, string(input), opts...)
	if err != nil {
		return nil, err
	}

	// Check for NODATA records - these are valid VPC flow log entries that indicate
	// no traffic was recorded during the aggregation interval. They don't have
	// timestamp data by design and should be skipped rather than treated as errors.
	// See: https://docs.aws.amazon.com/vpc/latest/userguide/flow-log-records.html
	if logStatus, ok := row.GetSourceValue("log_status"); ok && logStatus == VpcFlowLogStatusNoData {
		return nil, ErrNoDataRecord
	}

	return row, nil
}
