package aws_collection

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/hcl"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
)

type LambdaLogCollection struct {
	// all collections must embed collection.CollectionBase
	collection.CollectionBase[*LambdaLogCollectionConfig]
}

func NewLambdaLogCollection() collection.Collection {
	return &LambdaLogCollection{}
}

func (c *LambdaLogCollection) Identifier() string {
	return "aws_lambda_log"
}

func (c *LambdaLogCollection) GetSourceOptions() []row_source.RowSourceOption {
	return []row_source.RowSourceOption{
		artifact_source.WithRowPerLine(),
	}
}

func (c *LambdaLogCollection) GetRowSchema() any {
	return &aws_types.AwsLambdaLog{}
}

func (c *LambdaLogCollection) GetConfigSchema() hcl.Config {
	return &LambdaLogCollectionConfig{}
}

func (c *LambdaLogCollection) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	rawRecord, ok := row.(string)
	if !ok {
		return nil, fmt.Errorf("invalid row type: %T, expected string", row)
	}

	var record aws_types.AwsLambdaLog
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
	record.TpCollection = c.Identifier()
	if record.TpConnection == "" {
		record.TpConnection = c.Identifier() // TODO: #refactor figure out how to get connection (account ID?)
	}
	record.TpYear = int32(record.Timestamp.Year())
	record.TpMonth = int32(record.Timestamp.Month())
	record.TpDay = int32(record.Timestamp.Day())

	return record, nil
}
