package lambda_log

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	// "github.com/turbot/tailpipe-plugin-aws/rows"
	// "github.com/turbot/tailpipe-plugin-sdk/table"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

type LambdaLogMapper struct {
}

func (m *LambdaLogMapper) Identifier() string {
	return "lambda_log_mapper"
}

func (m *LambdaLogMapper) Map(_ context.Context, a any, _ ...mappers.MapOption[*LambdaLog]) (*LambdaLog, error) {
	row := &LambdaLog{}

	rawRow := ""

	switch v := a.(type) {
	case []byte:
		rawRow = string(v)
	case string:
		rawRow = v
	case *string:
		rawRow = *v
	default:
		return nil, fmt.Errorf("expected string, got %T", a)
	}

	slog.Error("rawRow ---->>>", rawRow)

	rawRow = strings.TrimSuffix(rawRow, "\n")
	fields := strings.Fields(rawRow)

	switch fields[0] {
	case "START", "END":
		row.LogType = &fields[0]
		row.RequestID = &fields[2]
	case "REPORT":
		row.LogType = &fields[0]
		row.RequestID = &fields[2]
		duration, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing duration: %w", err)
		}
		row.Duration = &duration
		billed, err := strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing billed duration: %w", err)
		}
		row.BilledDuration = &billed
		mem, err := strconv.Atoi(fields[12])
		if err != nil {
			return nil, fmt.Errorf("error parsing memory size: %w", err)
		}
		row.MemorySize = &mem
		maxMem, err := strconv.Atoi(fields[17])
		if err != nil {
			return nil, fmt.Errorf("error parsing max memory used: %w", err)
		}
		row.MaxMemoryUsed = &maxMem
	default:
		t := "LOG"
		row.LogType = &t

		ts, err := time.Parse(time.RFC3339, fields[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %w", err)
		}
		row.Timestamp = &ts

		row.RequestID = &fields[1]
		row.LogLevel = &fields[2]
		strip := fmt.Sprintf("%s%s", strings.Join(fields[:3], "\t"), "\t")
		stripped := strings.TrimPrefix(rawRow, strip)
		row.Message = &stripped
	}

	return row, nil
}
