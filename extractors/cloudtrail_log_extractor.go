package extractors

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

// CloudTrailLogExtractor is an extractor that receives JSON serialised CloudTrailLogBatch objects
// and extracts CloudTrailLog records from them
type CloudTrailLogExtractor struct {
}

// NewCloudTrailLogExtractor creates a new CloudTrailLogExtractor
func NewCloudTrailLogExtractor() artifact_source.Extractor {
	return &CloudTrailLogExtractor{}
}

func (c *CloudTrailLogExtractor) Identifier() string {
	return "cloudtrail_log_extractor"
}

// Extract unmarshalls the artifact data as an CloudTrailLogBatch and returns the CloudTrailLog records
func (c *CloudTrailLogExtractor) Extract(_ context.Context, a any) ([]any, error) {
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

	slog.Debug("CloudTrailLogExtractor", "record count", len(log.Records))
	var res = make([]any, len(log.Records))
	for i, record := range log.Records {
		res[i] = &record
	}
	return res, nil
}
