package vpc_flow_log

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/formats"
)

var allLogFields = []string{
	"account-id", "action", "az-id", "bytes", "dstaddr", "dstport", "end", "flow-direction", "instance-id", "interface-id", "log-status", "packets", "pkt-dst-aws-service", "pkt-dstaddr", "pkt-src-aws-service", "pkt-srcaddr", "protocol", "region", "reject-reason", "srcaddr", "srcport", "start", "sublocation-id", "sublocation-type", "subnet-id", "tcp-flags", "traffic-path", "type", "version", "vpc-id", "ecs-cluster-name", "ecs-cluster-arn", "ecs-container-instance-id", "ecs-container-instance-arn", "ecs-service-name", "ecs-task-definition-arn", "ecs-task-id", "ecs-task-arn", "ecs-container-id", "ecs-second-container-id",
}

// VPCFlowLogExtractor is an extractor that receives JSON serialised VPCFlowLogBatch objects
// and extracts VPCFlowLog records from them
type VPCFlowLogExtractor struct {
	Format formats.Format
}

// NewVPCFlowLogExtractor creates a new VPCFlowLogExtractor
func NewVPCFlowLogExtractor(format formats.Format) artifact_source.Extractor {
	return &VPCFlowLogExtractor{
		Format: format,
	}
}

func (c *VPCFlowLogExtractor) Identifier() string {
	return "vpc_flow_log_extractor"
}

// Extract unmarshalls the artifact data as an VPCFlowLogBatch and returns the VPCFlowLog records
func (c *VPCFlowLogExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// the expected input type is a JSON byte[] deserializable to VPCFlowLogBatch
	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	mappedData, err := ConvertToMapSlice(c, string(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("error in mapping the results in the form of []map[string]string")
	}

	var res = make([]any, len(mappedData))
	for i, record := range mappedData {

		res[i] = record
	}
	return res, nil
}

// ConvertToMapSlice converts space-separated string data into a slice of map[string]string
func ConvertToMapSlice(c *VPCFlowLogExtractor, data string) ([]map[string]string, error) {
	// Split the input into lines
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("insufficient data: needs at least a header and one row")
	}

	// Extract header row as keys
	keys := strings.Fields(lines[0]) // Splitting header by spaces

	if !isHeaderAvailableInLog(keys) {
		// Get the layout from the format
		// If the format it specifies in config then it wil be used otherwise the default format("version account-id interface-id srcaddr dstaddr srcport dstport protocol packets bytes start end action log-status") will be used.
		// The default format has been defined under vpc_flow_log_format_presets.go
		prop := c.Format.GetProperties()
		layout := prop["layout"]
		keys = strings.Fields(layout)

		newLines := make([]string, len(lines)+1)
		newLines[0] = layout
		copy(newLines[1:], lines)
		lines = newLines
		slog.Warn("Header is not available in the log", "Using the default layout", layout)
	}

	var result []map[string]string

	// Iterate over the data rows
	for _, line := range lines[1:] {
		values := strings.Fields(line) // Splitting row by spaces
		entry := make(map[string]string)

		// Map values to keys
		for i, key := range keys {
			if i < len(values) {
				entry[key] = values[i]
			} else {
				entry[key] = "" // Assign empty string if missing
			}
		}
		result = append(result, entry)
	}

	return result, nil
}

func isHeaderAvailableInLog(header []string) bool {
	headerAvailable := true

	for _, field := range header {
		if !slices.Contains(allLogFields, field) {
			headerAvailable = false
			break
		}
	}

	return headerAvailable
}
