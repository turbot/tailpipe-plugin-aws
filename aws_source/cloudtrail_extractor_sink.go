package aws_source

import (
	"context"
	"fmt"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/artifact"
	"github.com/turbot/tailpipe-plugin-sdk/events"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

// CloudtrailExtractorSink is an ExtractorSinkBase that receives AWSCloudTrailBatch objects and extracts rows from them
type CloudtrailExtractorSink struct {
	artifact.ExtractorSinkBase
}

// NewCloudtrailExtractorSink creates a new CloudtrailExtractorSink
func NewCloudtrailExtractorSink() *CloudtrailExtractorSink {
	return &CloudtrailExtractorSink{}
}

// Notify implements observable.Observer
func (s *CloudtrailExtractorSink) Notify(event events.Event) error {
	switch e := event.(type) {
	case *events.ArtifactExtracted:
		return s.ExtractArtifactRows(context.Background(), e.Request, e.Artifact)
	default:
		return fmt.Errorf("CloudtrailExtractorSink received unexpected event type: %T", e)
	}
}

func (c *CloudtrailExtractorSink) ExtractArtifactRows(ctx context.Context, req *proto.CollectRequest, a *types.Artifact) error {
	log, ok := a.Data.(*aws_types.AWSCloudTrailBatch)
	if !ok {
		return fmt.Errorf("expected AWSCloudTrailBatch, got %T", a.Data)
	}

	for _, record := range log.Records {
		// check context cancellation
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// call base OnRow
		err := c.OnRow(req, record, a.EnrichmentFields)
		if err != nil {
			return fmt.Errorf("error sending row: %w", err)
		}
	}

	return nil
}
