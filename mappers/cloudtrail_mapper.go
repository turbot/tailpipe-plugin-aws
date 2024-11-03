package mappers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_mapper"
)

// CloudtrailMapper is a mapper that receives CloudTrailBatch objects and extracts CloudTrailLog records from them
type CloudtrailMapper struct {
}

// NewCloudtrailMapper creates a new CloudtrailMapper
func NewCloudtrailMapper() artifact_mapper.Mapper {
	return &CloudtrailMapper{}
}

func (c *CloudtrailMapper) Identifier() string {
	return "cloudtrail_mapper"
}

// Map casts the data item as an CloudTrailBatch and returns the CloudTrailLog records
func (c *CloudtrailMapper) Map(_ context.Context, a any) ([]any, error) {
	// the expected input type is a JSON byte[] deserializable to CloudTrailBatch
	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	// decode json ito CloudTrailBatch
	var log rows.CloudTrailBatch
	err := json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	slog.Debug("CloudwatchMapper", "record count", len(log.Records))
	var res = make([]any, len(log.Records))
	for i, record := range log.Records {
		res[i] = record
	}
	return res, nil
}
