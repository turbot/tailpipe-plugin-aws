package elb_access_log

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

// ElbAccessLogExtractor is an extractor that receives JSON serialised ElbAccessLogBatch objects
// and extracts ElbAccessLog records from them
type ElbAccessLogExtractor struct {
}

// NewElbAccessLogExtractor creates a new ElbAccessLogExtractor
func NewElbAccessLogExtractor() artifact_source.Extractor {
	return &ElbAccessLogExtractor{}
}

func (c *ElbAccessLogExtractor) Identifier() string {
	return "elb_access_log_extractor"
}

// Extract unmarshalls the artifact data as an ElbAccessLogBatch and returns the ElbAccessLog records
func (c *ElbAccessLogExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// the expected input type is a JSON byte[] deserializable to ElbAccessLogBatch
	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	// decode json ito ElbAccessLogBatch
	var log ElbAccessLogBatch
	err := json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	slog.Debug("ElbAccessLogExtractor", "record count", len(log.Records))
	var res = make([]any, len(log.Records))
	for i, record := range log.Records {
		res[i] = &record
	}
	return res, nil
}
