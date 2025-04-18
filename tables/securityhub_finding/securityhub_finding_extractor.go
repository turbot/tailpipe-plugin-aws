package securityhub_finding

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

// SecurityHubFindingExtractor is an extractor that receives JSON serialised CloudTrailLogBatch objects
// and extracts SecurityHubFinding records from them
type SecurityHubFindingExtractor struct {
}

// NewCloudTrailLogExtractor creates a new SecurityHubFindingExtractor
func NewCloudTrailLogExtractor() artifact_source.Extractor {
	return &SecurityHubFindingExtractor{}
}

func (c *SecurityHubFindingExtractor) Identifier() string {
	return "cloudtrail_log_extractor"
}

// Extract unmarshalls the artifact data as an CloudTrailLogBatch and returns the SecurityHubFinding records
func (c *SecurityHubFindingExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// the expected input type is a JSON byte[] deserializable to CloudTrailLogBatch
	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	// decode json ito CloudTrailLogBatch
	var log SecurityHubFinding
	err := json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	slog.Debug("SecurityHubFindingExtractor", "record count", len(log.Findings))
	var res = make([]any, len(log.Findings))
	for i, record := range log.Findings {
		res[i] = &record
	}
	return res, nil
}
