package detailed_billing_report

import (
	"bufio"
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

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

type DetailedBillingReportExtractor struct{}

func NewDetailedBillingReportExtractor() artifact_source.Extractor {
	return &DetailedBillingReportExtractor{}
}

func (e *DetailedBillingReportExtractor) Identifier() string {
	return "detailed_billing_report_extractor"
}

func (e *DetailedBillingReportExtractor) Extract(_ context.Context, a any) ([]any, error) {
	data, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte, got %T", a)
	}
	return extractFromCSV(bytes.NewReader(data))
}

func extractFromCSV(reader io.Reader) ([]any, error) {
	var headers []string

	// Preprocess to skip invalid header/comment lines
	// Handle cases where the CSV file may start with one or more non-CSV comment or metadata lines.
	// Since we can't assume it's always a single comment line, we read lines in a loop
	// until we find a valid header row — typically identified by having more than one column.
	// Once found, treat it as the actual CSV header and proceed with parsing.
	scanner := bufio.NewScanner(reader)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) < 2 {
		return nil, fmt.Errorf("not enough lines to determine valid CSV content")
	}

	headerIndex := -1
	for i := 0; i < len(lines)-1; i++ {
		h := parseCSVLine(lines[i])
		d := parseCSVLine(lines[i+1])
		if len(h) > 1 && len(h) == len(d) {
			headers = h
			headerIndex = i
			break
		} else {
			slog.Debug("skipping potential header line:", "line", lines[i])
		}
	}

	if headerIndex == -1 {
		return nil, fmt.Errorf("could not identify valid CSV header")
	}

	cleanCSV := strings.Join(lines[headerIndex:], "\n")
	csvReader := csv.NewReader(strings.NewReader(cleanCSV))
	csvReader.FieldsPerRecord = -1

	// Read the confirmed header
	_, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %v", err)
	}

	var records []*DetailedBillingReport
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row: %v", err)
		}

		recordMap := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				recordMap[header] = row[i]
			} else {
				recordMap[header] = ""
			}
		}

		record := &DetailedBillingReport{}
		err = mapToStruct(recordMap, record)
		if err != nil {
			return nil, fmt.Errorf("error mapping values: %w", err)
		}

		if record.RecordID != nil && strings.HasPrefix(*record.RecordID, "AccountTotal:") {
			continue
		}

		records = append(records, record)
	}

	result := make([]any, len(records))
	for i, r := range records {
		result[i] = r
	}
	return result, nil
}

func parseCSVLine(line string) []string {
	r := csv.NewReader(strings.NewReader(line))
	r.FieldsPerRecord = -1
	record, err := r.Read()
	if err != nil {
		return []string{}
	}
	return record
}

func mapToStruct(data map[string]string, target *DetailedBillingReport) error {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	tags := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonKey := trimOmitempty(field.Tag.Get("json"))
		if jsonKey == "" {
			continue
		}

		fieldVal := v.Field(i)

		// Handle tag fields (e.g., aws:*, user:*) dynamically from data map keys
		if field.Name == "Tags" && field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Map {
			for k, v := range data {
				if strings.Contains(k, ":") && v != "" {
					tags[k] = v
				}
			}
			if len(tags) > 0 && fieldVal.CanSet() {
				fieldVal.Set(reflect.ValueOf(&tags))
			}
			continue
		}

		// Match normal struct fields by json key
		value, exists := data[jsonKey]
		if !exists || value == "" {
			continue
		}

		if fieldVal.Kind() == reflect.Ptr {
			switch fieldVal.Type().Elem().Kind() {
			case reflect.String:
				fieldVal.Set(reflect.ValueOf(&value))

			case reflect.Float64:
				if f, err := strconv.ParseFloat(value, 64); err == nil {
					fieldVal.Set(reflect.ValueOf(&f))
				}

			case reflect.Struct:
				if fieldVal.Type().Elem() == reflect.TypeOf(time.Time{}) {
					formats := []string{"2006-01-02 15:04:05", "2006/01/02 15:04:05"}
					for _, format := range formats {
						if t, err := time.Parse(format, value); err == nil {
							fieldVal.Set(reflect.ValueOf(&t))
							break
						}
					}
				}
			}
		}
	}

	return nil
}


func trimOmitempty(tag string) string {
	if idx := len(tag) - len(",omitempty"); idx > 0 && tag[idx:] == ",omitempty" {
		return tag[:idx]
	}
	return tag
}
