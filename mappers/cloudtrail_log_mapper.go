package mappers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// CloudTrailLogMapper is a mapper that receives CloudTrailLogBatch objects and extracts CloudTrailLog records from them
type CloudTrailLogMapper struct {
}

// NewCloudTrailLogMapper creates a new CloudTrailLogMapper
func NewCloudTrailLogMapper() table.Mapper[*rows.CloudTrailLog] {
	return &CloudTrailLogMapper{}
}

func (c *CloudTrailLogMapper) Identifier() string {
	return "cloudtrail_log_mapper"
}

// Map casts the data item as an CloudTrailLogBatch and returns the CloudTrailLog records
func (c *CloudTrailLogMapper) Map(_ context.Context, a any) ([]*rows.CloudTrailLog, error) {
	// the expected input type is a JSON byte[] deserializable to CloudTrailLogBatch
	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	// decode json ito CloudTrailLogBatch
	var log rows.CloudTrailLogBatch
	err := json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	slog.Debug("CloudwatchMapper", "record count", len(log.Records))
	var res = make([]*rows.CloudTrailLog, len(log.Records))
	for i, record := range log.Records {
		res[i] = &record
	}
	return res, nil
}
