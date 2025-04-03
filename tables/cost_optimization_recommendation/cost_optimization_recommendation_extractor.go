package cost_optimization_recommendation

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

// CostOptimizationRecommendationExtractor is an extractor that receives JSON or CSV data
// and extracts CostOptimizationRecommendation records from them
type CostOptimizationRecommendationExtractor struct {
}

// NewCostOptimizationRecommendationExtractor creates a new CostOptimizationRecommendationExtractor
func NewCostOptimizationRecommendationExtractor() artifact_source.Extractor {
	return &CostOptimizationRecommendationExtractor{}
}

func (c *CostOptimizationRecommendationExtractor) Identifier() string {
	return "cost_optimization_recommendation_extractor"
}

func (c *CostOptimizationRecommendationExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// Assert that 'a' is of type []byte
	data, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte, got %T", a)
	}

	// Validate CSV structure
	if err := validateCSV(bytes.NewReader(data)); err != nil {
		return nil, nil
	}

	// If not a ZIP, assume it is a raw CSV
	return extractFromCSV(bytes.NewReader(data))
}

// Validate that the input can be parsed as a CSV file
func validateCSV(reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	// Try to read headers
	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %v", err)
	}

	// Optionally: try to read the first data row (to confirm it's not header-only)
	_, err = csvReader.Read()
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read CSV data row: %v", err)
	}

	// Also optionally: check headers are not empty
	if len(headers) == 0 {
		return fmt.Errorf("CSV header is empty")
	}

	return nil
}

// Extract data from a CSV reader
func extractFromCSV(reader io.Reader) ([]any, error) {
	csvReader := csv.NewReader(reader)

	// Read CSV headers
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %v", err)
	}

	var records []*CostOptimizationRecommendation

	// Read all records
	for {
		// Read the next row
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("error reading CSV row", "error", err)
			continue
		}

		// Create a new CostOptimizationRecommendation record
		record := &CostOptimizationRecommendation{}

		// Create a map of column name to value
		recordMap := make(map[string]string)
		for i, value := range row {
			if i < len(headers) {
				recordMap[headers[i]] = value
			}
		}

		// Map the values to the record
		record.MapValues(recordMap)

		// We should have null value if the map value is empty
		if len(*record.RecommendedResourceDetails) == 0 {
			record.RecommendedResourceDetails = nil
		}
		if len(*record.CurrentResourceDetails) == 0 {
			record.CurrentResourceDetails = nil
		}

		// Add the record to our list
		records = append(records, record)
	}

	// Convert to a slice of any
	return convertToAny(records), nil
}

// Convert a slice of CostOptimizationRecommendation to a slice of any
func convertToAny(records []*CostOptimizationRecommendation) []any {
	result := make([]any, len(records))
	for i, record := range records {
		result[i] = record
	}
	return result
}

// MapValues maps values from a record map to a CostOptimizationRecommendation struct
func (value *CostOptimizationRecommendation) MapValues(recordMap map[string]string) {
	// Initialize nested maps
	initializeNestedMaps(value)

	// Get the type of CostOptimizationRecommendation
	t := reflect.TypeOf(*value)
	v := reflect.ValueOf(value).Elem()

	// Loop through each field in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip CommonFields
		if field.Name == "CommonFields" {
			continue
		}

		// Get the JSON tag for the field
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		// Remove ",omitempty" from the tag
		jsonTag = strings.Split(jsonTag, ",")[0]

		// Check if we have a value for this field
		strVal, ok := recordMap[jsonTag]
		if !ok || strVal == "" {
			continue
		}

		// Get the field value
		fieldVal := v.Field(i)

		// Ensure the field is settable
		if !fieldVal.CanSet() {
			continue
		}

		// Special handling for map[string]interface{} fields
		if jsonTag == "recommended_resource_details" || jsonTag == "current_resource_details" {

			var parsedMap map[string]interface{}
			err := json.Unmarshal([]byte(strVal), &parsedMap)
			if err == nil {
				// Set a pointer to the parsed map
				fieldVal.Set(reflect.New(field.Type.Elem()))
				fieldVal.Elem().Set(reflect.ValueOf(parsedMap))
			}
			continue
		}

		// For nested maps like "tags", we need to handle them specially
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Map {
			assignToNestedMap(value, jsonTag, strVal)
			continue
		}

		// Handle different field types
		switch field.Type.Kind() {
		case reflect.Ptr:
			// Get the type of the pointed-to value
			ptrType := field.Type.Elem()
			switch ptrType.Kind() {
			case reflect.String:
				val := strVal
				if val != "" {
					fieldVal.Set(reflect.ValueOf(&val))
				}
			case reflect.Int:
				val, err := strconv.Atoi(strVal)
				if err == nil && val > 0 {
					fieldVal.Set(reflect.ValueOf(&val))
				}
			case reflect.Float64:
				val, err := strconv.ParseFloat(strVal, 64)
				if err == nil && val > 0 {
					fieldVal.Set(reflect.ValueOf(&val))
				}
			case reflect.Bool:
				if strVal == "true" || strVal == "TRUE" || strVal == "True" || strVal == "1" {
					val := true
					fieldVal.Set(reflect.ValueOf(&val))
				} else if strVal == "false" || strVal == "FALSE" || strVal == "False" || strVal == "0" {
					val := false
					fieldVal.Set(reflect.ValueOf(&val))
				}
			case reflect.Struct:
				// Handle time.Time specifically
				if ptrType.String() == "time.Time" {
					layout := "Mon Jan 2 15:04:05 UTC 2006"
					val, err := time.Parse(layout, strVal)
					if err == nil {
						fieldVal.Set(reflect.ValueOf(&val))
						break
					}
				}
			}
		}
	}
}

// Initialize nested maps in a CostOptimizationRecommendation struct
func initializeNestedMaps(value *CostOptimizationRecommendation) {
	// Initialize the Tags map if it doesn't exist
	if value.Tags == nil {
		tags := make(map[string]string)
		value.Tags = &tags
	}
}

// Assign a value to a nested map in a CostOptimizationRecommendation struct
func assignToNestedMap(value *CostOptimizationRecommendation, key, strVal string) {
	// Currently, only the Tags field is a nested map
	if strings.HasPrefix(key, "tags") {
		// Parse the key to extract the tag name, e.g., "tags.Name" -> "Name"
		parts := strings.SplitN(key, ".", 2)
		if len(parts) == 2 {
			tagName := parts[1]
			(*value.Tags)[tagName] = strVal
		}
	}
}
