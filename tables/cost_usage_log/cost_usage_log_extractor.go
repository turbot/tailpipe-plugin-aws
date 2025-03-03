package cost_usage_log

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/stoewer/go-strcase"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

// CostUsageLogExtractor is an extractor that receives JSON serialised CostUsageLogBatch objects
// and extracts CostUsageLog records from them
type CostUsageLogExtractor struct {
}

// NewCostUsageLogExtractor creates a new CostUsageLogExtractor
func NewCostUsageLogExtractor() artifact_source.Extractor {
	return &CostUsageLogExtractor{}
}

func (c *CostUsageLogExtractor) Identifier() string {
	return "cost_usage_log_extractor"
}

func (c *CostUsageLogExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// Assert that 'a' is of type []byte
	zipData, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte, got %T", a)
	}

	// Open the zip archive
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("failed to open zip reader: %v", err)
	}

	var records []CostAndUsageLog

	// Iterate through the files in the archive
	for _, file := range zipReader.File {
		// Open the file
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file in zip: %v", err)
		}
		defer rc.Close()

		// Check if the file is a CSV
		if strings.HasSuffix(file.Name, ".csv") {
			// Parse the CSV file
			csvReader := csv.NewReader(rc)
			headers, err := csvReader.Read()
			if err != nil {
				return nil, fmt.Errorf("error reading CSV header: %v", err)
			}

			// Read the remaining rows
			for {
				row, err := csvReader.Read()
				if err == io.EOF {
					break // End of file reached
				}
				if err != nil {
					return nil, fmt.Errorf("error reading CSV row: %v", err)
				}

				// Create a map to hold the row data
				recordMap := make(map[string]string)
				for i, header := range headers {
					replaceStr := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(header, "/", "_"), ":", "_"), " ", "_")
					snakeCaseKey := strcase.SnakeCase(replaceStr)
					if i < len(row) {
						recordMap[snakeCaseKey] = row[i]
					} else {
						recordMap[snakeCaseKey] = ""
					}
				}

				// Marshal the map to JSON
				jsonData, err := json.Marshal(recordMap)
				if err != nil {
					return nil, fmt.Errorf("error marshaling record to JSON: %v", err)
				}

				// Unmarshal the JSON into a CostAndUsageLog struct
				var record CostAndUsageLog
				if err := json.Unmarshal(jsonData, &record); err != nil {
					return nil, fmt.Errorf("error unmarshaling JSON to struct: %v", err)
				}

				// Append the record to the slice
				records = append(records, record)
			}
		}
	}

	// Convert the slice of structs to a slice of empty interfaces
	res := make([]any, len(records))
	for i, record := range records {
		res[i] = &record
	}

	return res, nil
}
