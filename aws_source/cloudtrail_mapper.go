package aws_source

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_mapper"
	"github.com/turbot/tailpipe-plugin-sdk/types"
	"log/slog"
)

// CloudtrailMapper is an Mapper that receives AWSCloudTrailBatch objects and extracts AWSCloudTrail records from them
type CloudtrailMapper struct {
}

// NewCloudtrailMapper creates a new CloudtrailMapper
func NewCloudtrailMapper() artifact_mapper.Mapper {
	return &CloudtrailMapper{}
}

func (c *CloudtrailMapper) Identifier() string {
	return "cloudtrail_mapper"
}

// Map casts the data item as an AWSCloudTrailBatch and returns the records
func (c *CloudtrailMapper) Map(_ context.Context, a *types.RowData) ([]*types.RowData, error) {
	// the expected input type is a JSON byte[] deserializable to AWSCloudTrailBatch
	jsonBytes, ok := a.Data.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	// decode json ito AWSCloudTrailBatch
	var log aws_types.AWSCloudTrailBatch
	err := json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	slog.Debug("CloudwatchMapper", "record count", len(log.Records))
	var res = make([]*types.RowData, len(log.Records))
	for i, record := range log.Records {
		res[i] = types.NewData(record)
	}
	return res, nil
}
