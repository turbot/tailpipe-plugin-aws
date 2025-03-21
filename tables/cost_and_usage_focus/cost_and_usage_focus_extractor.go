package cost_and_usage_focus

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

type CostUsageFocusExtractor struct {
}

// NewCostUsageReportExtractor creates a new CostUsageReportExtractor
func NewCostUsageFocusExtractor() artifact_source.Extractor {
	return &CostUsageFocusExtractor{}
}

func (m *CostUsageFocusExtractor) Identifier() string {
	return "cost_and_usage_focus_extractor"
}

func (c *CostUsageFocusExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// Assert that 'a' is of type []byte
	data, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte, got %T", a)
	}

	// If not a ZIP, assume it is a raw CSV
	return extractFromCSV(bytes.NewReader(data))
}

// Extract data from a CSV reader
func extractFromCSV(reader io.Reader) ([]any, error) {
	csvReader := csv.NewReader(reader)

	// Read CSV headers
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %v", err)
	}

	var records []*CostUsageFocus

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

			// Assign row data
			if i < len(row) {
				recordMap[header] = row[i]
			} else {
				recordMap[header] = ""
			}
		}

		record := &CostUsageFocus{}
		err = record.mapValues(recordMap)
		if err != nil {
			return nil, fmt.Errorf("error in mapping the value to struct: %w", err)
		}

		// Append the record
		records = append(records, record)
	}

	// Convert to a slice of empty interfaces
	return convertToAny(records), nil
}

// Convert a slice of structs to a slice of empty interfaces
func convertToAny(records []*CostUsageFocus) []any {
	res := make([]any, len(records))
	for i, record := range records {
		res[i] = record
	}
	return res
}

// mapValues maps the values from a map[string]string into the Focus1_0 struct
func (log *CostUsageFocus) mapValues(row map[string]string) error {
	v := reflect.ValueOf(log).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		// Extract the JSON key, ignoring omitempty if present
		jsonKey := jsonTag
		if idx := len(jsonTag) - len(",omitempty"); idx > 0 && jsonTag[idx:] == ",omitempty" {
			jsonKey = jsonTag[:idx]
		}

		if value, exists := row[jsonKey]; exists {
			structField := v.Field(i)

			if structField.Kind() == reflect.Ptr {
				elemType := structField.Type().Elem()

				switch elemType.Kind() {
				case reflect.String:
					if value != "" {
						structField.Set(reflect.ValueOf(&value))
					}

				case reflect.Float64:
					if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
						structField.Set(reflect.ValueOf(&floatVal))
					}

				case reflect.Int64:
					if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
						structField.Set(reflect.ValueOf(&intVal))
					}

				case reflect.Struct:
					if elemType == reflect.TypeOf(time.Time{}) {
						if timeVal, err := time.Parse(time.RFC3339, value); err == nil {
							structField.Set(reflect.ValueOf(&timeVal))
						}
					}
				case reflect.Map:
					// Handle map[string]string parsing
					if elemType.Key().Kind() == reflect.String && elemType.Elem().Kind() == reflect.String {
						var parsedMap map[string]string
						if err := json.Unmarshal([]byte(value), &parsedMap); err == nil {
							structField.Set(reflect.ValueOf(&parsedMap))
						}
					}
				}
			}
		}
	}

	return nil
}
