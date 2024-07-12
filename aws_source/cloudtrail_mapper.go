package aws_source

import (
	"context"
	"fmt"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
)

// CloudtrailMapper is an Mapper that receives AWSCloudTrailBatch objects and extracts rows from them
type CloudtrailMapper struct {
}

// NewCloudtrailMapper creates a new CloudtrailMapper
func NewCloudtrailMapper() *CloudtrailMapper {
	return &CloudtrailMapper{}
}

// Map casts the data item as an AWSCloudTrailBatch and returns the records
func (c *CloudtrailMapper) Map(ctx context.Context, req *proto.CollectRequest, a any) ([]any, error) {
	log, ok := a.(*aws_types.AWSCloudTrailBatch)
	if !ok {
		return nil, fmt.Errorf("expected AWSCloudTrailBatch, got %T", a)
	}
	var res = make([]any, len(log.Records))
	for i, record := range log.Records {
		res[i] = record
	}
	return res, nil
}
