package aws_source

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"log/slog"
)

// CloudtrailMapper is an Mapper that receives AWSCloudTrailBatch objects and extracts AWSCloudTrail records from them
type CloudtrailMapper struct {
}

// NewCloudtrailMapper creates a new CloudtrailMapper
func NewCloudtrailMapper() *CloudtrailMapper {
	return &CloudtrailMapper{}
}

func (c *CloudtrailMapper) Identifier() string {
	return "cloudtrail_mapper"
}

// Map casts the data item as an AWSCloudTrailBatch and returns the records
func (c *CloudtrailMapper) Map(_ context.Context, _ *proto.CollectRequest, a any) ([]any, error) {
	// the expected input type is a JSON string deserializable to  AWSCloudTrailBatch
	// convert from char[] to string

	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}
	jsonString := string(jsonBytes)

	// decode json ito AWSCloudTrailBatch
	var log aws_types.AWSCloudTrailBatch
	err := json.Unmarshal([]byte(jsonString), &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	slog.Debug("CloudtrailMapper", "record count", len(log.Records))
	var res = make([]any, len(log.Records))
	for i, record := range log.Records {
		res[i] = record
	}
	return res, nil
}
