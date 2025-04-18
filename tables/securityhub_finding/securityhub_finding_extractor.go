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
func NewSecurityHubFindingExtractor() artifact_source.Extractor {
	return &SecurityHubFindingExtractor{}
}

func (c *SecurityHubFindingExtractor) Identifier() string {
	return "cloudtrail_log_extractor"
}

// Extract unmarshalls the artifact data as an CloudTrailLogBatch and returns the SecurityHubFinding records
func (c *SecurityHubFindingExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// the expected input type is a JSON byte[] deserializable to CloudTrailLogBatch
	var jsonBytes []byte

	switch v := a.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		return nil, fmt.Errorf("expected []byte or string, got %T", a)
	}

	// decode json ito CloudTrailLogBatch
	var log DetailFindingsData
	err := json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	slog.Debug("SecurityHubFindingExtractor", "record count", len(log.Detail.Findings))
	findings := toMapSecurityHubFinding(log)
	var res = make([]any, len(findings))
	for i, record := range findings {
		res[i] = &record
	}
	return res, nil
}

func toMapSecurityHubFinding(findingsRow DetailFindingsData) []SecurityHubFinding {
	var findings []SecurityHubFinding

	for _, finding := range findingsRow.Detail.Findings {
		f := &SecurityHubFinding{}
		f.Version = findingsRow.Version
		f.ID = findingsRow.ID
		f.DetailType = findingsRow.DetailType
		f.Source = findingsRow.Source
		f.Account = findingsRow.Account
		f.Time = findingsRow.Time

		// Findings field
		f.AwsAccountId = finding.AwsAccountId
		f.CreatedAt = finding.CreatedAt
		f.Description = finding.Description
		f.GeneratorId = finding.GeneratorId
		f.FindingId = finding.Id
		f.ProductArn = finding.ProductArn
		f.ProductFields = finding.ProductFields
		f.ProductName = finding.ProductName
		f.Remediation = finding.Remediation
		f.Resources = finding.Resources
		// findings = append(findings, *f)
	}
	return findings
}
