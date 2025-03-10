package cost_and_usage_report

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"

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
	return "cost_and_usage_report_extractor"
}

func (c *CostUsageLogExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// Assert that 'a' is of type []byte
	data, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte, got %T", a)
	}

	// Check if the data is a ZIP archive
	// We can have the zip(.csv.zip) and non-zip(.csv.gz) file object stored in S3 bucket
	if isZip(data) {
		return extractFromZip(data)
	}

	// If not a ZIP, assume it is a raw CSV
	return extractFromCSV(bytes.NewReader(data))
}

// Function to check if the provided data is a ZIP archive
func isZip(data []byte) bool {
	// ZIP files start with "PK" signature (0x50 0x4B)
	return len(data) > 4 && data[0] == 'P' && data[1] == 'K'
}

// Extract data from a ZIP archive
func extractFromZip(zipData []byte) ([]any, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("failed to open zip reader: %v", err)
	}

	var records []*CostAndUsageReport

	// Iterate through the files in the archive
	for _, file := range zipReader.File {
		// Open the file
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file in zip: %v", err)
		}
		defer rc.Close()

		// Process only CSV files
		if strings.HasSuffix(file.Name, ".csv") {
			csvRecords, err := extractFromCSV(rc)
			if err != nil {
				return nil, fmt.Errorf("error processing CSV file %s: %v", file.Name, err)
			}

			// Append extracted records
			for _, r := range csvRecords {
				if rec, ok := r.(*CostAndUsageReport); ok {
					records = append(records, rec)
				}
			}
		}
	}

	// Convert to a slice of empty interfaces
	return convertToAny(records), nil
}

// Extract data from a CSV reader
func extractFromCSV(reader io.Reader) ([]any, error) {
	csvReader := csv.NewReader(reader)

	// Read CSV headers
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %v", err)
	}

	var records []*CostAndUsageReport

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
			// Normalize header names
			replaceStr := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(header, "/", "_"), ":", "_"), " ", "_"), ".", "_"), "a_r_n", "arn")
			snakeCaseKey := strcase.SnakeCase(replaceStr)

			// Assign row data
			if i < len(row) {
				recordMap[snakeCaseKey] = row[i]
			} else {
				recordMap[snakeCaseKey] = ""
			}
		}

		record := &CostAndUsageReport{}
		record.MapValues(recordMap)

		// Append the record
		records = append(records, record)
	}

	// Convert to a slice of empty interfaces
	return convertToAny(records), nil
}

// Convert a slice of structs to a slice of empty interfaces
func convertToAny(records []*CostAndUsageReport) []any {
	res := make([]any, len(records))
	for i, record := range records {
		res[i] = record
	}
	return res
}

// MapValues dynamically assigns values based on JSON tags, iterating over the map instead of struct fields
func (value *CostAndUsageReport) MapValues(recordMap map[string]string) {
	if value == nil {
		panic("CostAndUsageReport is nil") // Prevent nil dereference
	}

	v := reflect.ValueOf(value).Elem()
	t := v.Type()

	// Create a map to store json tag → field index mapping
	jsonTagToField := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		jsonKey := strings.Split(jsonTag, ",")[0]
		if jsonKey != "" && jsonKey != "-" {
			jsonTagToField[jsonKey] = i
		}
	}

	// Ensure required maps are initialized
	initializeNestedMaps(value)

	// Iterate over map keys instead of struct fields
	for key, strVal := range recordMap {
		if fieldIndex, exists := jsonTagToField[key]; exists {
			structField := v.Field(fieldIndex)

			// Ensure the field is addressable and settable
			if !structField.CanSet() {
				slog.Debug("Skipping field %s: not settable", key)
				continue
			}

			// Assign value based on field type
			if structField.Kind() == reflect.Ptr {
				elemType := structField.Type().Elem()

				switch elemType.Kind() {
				case reflect.String:
					val := strVal
					if val == "" {
						structField.Set(reflect.Zero(structField.Type())) // Set to nil
					} else {
						structField.Set(reflect.ValueOf(&val))
					}
				case reflect.Float64:
					if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
						structField.Set(reflect.ValueOf(&floatVal))
					} else {
						defaultVal := float64(0)
						structField.Set(reflect.ValueOf(&defaultVal))
					}
				case reflect.Int64:
					if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
						structField.Set(reflect.ValueOf(&intVal))
					} else {
						defaultVal := int64(0)
						structField.Set(reflect.ValueOf(&defaultVal))
					}
				case reflect.Struct:
					if elemType == reflect.TypeOf(time.Time{}) {
						if timeVal, err := time.Parse(time.RFC3339, strVal); err == nil {
							structField.Set(reflect.ValueOf(&timeVal))
						} else {
							defaultVal := time.Time{}
							structField.Set(reflect.ValueOf(&defaultVal))
						}
					}
				}
			}
		} else {
			// Handle additional attributes
			assignToNestedMap(value, key, strVal)
		}
	}
}

// initializeNestedMaps ensures that all map fields in the struct are initialized
func initializeNestedMaps(value *CostAndUsageReport) {
	if value.Product == nil {
		value.Product = new(map[string]interface{})
		*value.Product = make(map[string]interface{})
	}
	if value.Reservation == nil {
		value.Reservation = new(map[string]interface{})
		*value.Reservation = make(map[string]interface{})
	}
	if value.ResourceTags == nil {
		value.ResourceTags = new(map[string]interface{})
		*value.ResourceTags = make(map[string]interface{})
	}
	if value.CostCategory == nil {
		value.CostCategory = new(map[string]interface{})
		*value.CostCategory = make(map[string]interface{})
	}
	if value.Discount == nil {
		value.Discount = new(map[string]interface{})
		*value.Discount = make(map[string]interface{})
	}
}

// assignToNestedMap assigns values to dynamically mapped attributes
func assignToNestedMap(value *CostAndUsageReport, key, strVal string) {
	switch {
	case strings.HasPrefix(key, "product_"):
		key = strings.Replace(key, "product_", "", 1)
		(*value.Product)[key] = strVal
	case strings.HasPrefix(key, "cost_category_"):
		key = strings.Replace(key, "cost_category_", "", 1)
		(*value.CostCategory)[key] = strVal
	case strings.HasPrefix(key, "reservation_"):
		key = strings.Replace(key, "reservation_", "", 1)
		(*value.Reservation)[key] = strVal
	case strings.HasPrefix(key, "resource_tags_"):
		key = strings.Replace(key, "resource_tags_", "", 1)
		(*value.ResourceTags)[key] = strVal
	case strings.HasPrefix(key, "discount_"):
		key = strings.Replace(key, "discount_", "", 1)
		(*value.Discount)[key] = strVal
	}
}
