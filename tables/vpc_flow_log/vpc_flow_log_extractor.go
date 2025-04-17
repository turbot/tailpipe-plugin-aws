package vpc_flow_log

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/formats"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

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

// Extract unmarshalls the artifact data as a VPCFlowLogBatch and returns the VPCFlowLog records
func (c *VPCFlowLogExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// The expected input type is a JSON byte[] deserializable to VPCFlowLogBatch
	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	mappedData, err := toMapSlice(c, string(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("error in mapping the results to []map[string]string: %w", err)
	}

	fields := getValidTokensAndColumnNames()
	res := make([]any, 0, len(mappedData))

	for _, record := range mappedData {
		row := &types.DynamicRow{}
		rowColumnMap := make(map[string]string)

		for k, v := range record {
			if columnName, ok := fields[k]; ok {
				if v != VpcFlowLogTableNilValue {
					rowColumnMap[columnName] = v
				}
			}
		}

		err := row.InitialiseFromMap(rowColumnMap)
		if err != nil {
			return nil, fmt.Errorf("error initialising dynamic row: %w", err)
		}

		res = append(res, row)
	}

	return res, nil
}

// toMapSlice converts space-separated string data into a slice of map[string]string
func toMapSlice(c *VPCFlowLogExtractor, data string) ([]map[string]string, error) {
	// Split the input into lines
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("insufficient data: needs at least a header and one row")
	}

	// Extract header row as keys
	keys := strings.Fields(lines[0]) // Splitting header by spaces

	if !isValidHeader(keys) {
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

func isValidHeader(header []string) bool {
	headerAvailable := true
	allKeys := getValidTokensAndColumnNames()

	for _, field := range header {
		if _, exists := allKeys[field]; !exists {
			headerAvailable = false
			break
		}
	}

	return headerAvailable
}
