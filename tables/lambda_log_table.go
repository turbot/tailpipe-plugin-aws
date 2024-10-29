package tables

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type LambdaLogTable struct {
	// all tables must embed table.TableBase
	table.TableBase[*LambdaLogTableConfig]
}

func NewLambdaLogTable() table.Table {
	return &LambdaLogTable{}
}

func (c *LambdaLogTable) Identifier() string {
	return "aws_lambda_log"
}

func (c *LambdaLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
	}
}

func (c *LambdaLogTable) GetRowSchema() any {
	return &AwsLambdaLog{}
}

func (c *LambdaLogTable) GetConfigSchema() parse.Config {
	return &LambdaLogTableConfig{}
}

func (c *LambdaLogTable) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	rawRecord, ok := row.(string)
	if !ok {
		return nil, fmt.Errorf("invalid row type: %T, expected string", row)
	}

	var record AwsLambdaLog
	if sourceEnrichmentFields != nil {
		record.CommonFields = *sourceEnrichmentFields

		ts := time.UnixMilli(int64(sourceEnrichmentFields.TpTimestamp))
		record.Timestamp = &ts
	}

	// remove trailing newline
	rawRecord = strings.TrimSuffix(rawRecord, "\n")
	fields := strings.Fields(rawRecord)

	switch fields[0] {
	case "START", "END":
		record.LogType = &fields[0]
		record.RequestID = &fields[2]
	case "REPORT":
		record.LogType = &fields[0]
		record.RequestID = &fields[2]
		duration, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing duration: %w", err)
		}
		record.Duration = &duration
		billed, err := strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing billed duration: %w", err)
		}
		record.BilledDuration = &billed
		mem, err := strconv.Atoi(fields[12])
		if err != nil {
			return nil, fmt.Errorf("error parsing memory size: %w", err)
		}
		record.MemorySize = &mem
		maxMem, err := strconv.Atoi(fields[17])
		if err != nil {
			return nil, fmt.Errorf("error parsing max memory used: %w", err)
		}
		record.MaxMemoryUsed = &maxMem
	default:
		t := "LOG"
		record.LogType = &t
		// TODO: #enrich should we overwrite the timestamp with that in the log entry?
		record.RequestID = &fields[1]
		record.LogLevel = &fields[2]
		strip := fmt.Sprintf("%s%s", strings.Join(fields[:3], "\t"), "\t")
		stripped := strings.TrimPrefix(rawRecord, strip)
		record.Message = &stripped
	}

	// Record standardization
	record.TpID = xid.New().String()
	record.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))

	// Hive fields
	record.TpPartition = c.Identifier()
	if record.TpIndex == "" {
		record.TpIndex = c.Identifier() // TODO: #refactor figure out how to get connection (account ID?)
	}
	// convert to date in format yy-mm-dd
	record.TpDate = record.Timestamp.Format("2006-01-02")

	return record, nil
}
