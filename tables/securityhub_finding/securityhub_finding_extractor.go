package securityhub_finding

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

// SecurityHubFindingExtractor is an extractor that receives JSON serialised SecurityHub findings
// and extracts SecurityHubFinding records from them
type SecurityHubFindingExtractor struct {
}

// NewSecurityHubFindingExtractor creates a new SecurityHubFindingExtractor
func NewSecurityHubFindingExtractor() artifact_source.Extractor {
	return &SecurityHubFindingExtractor{}
}

func (c *SecurityHubFindingExtractor) Identifier() string {
	return "securityhub_finding_extractor"
}

// Extract unmarshalls the artifact data as SecurityHub findings and returns the SecurityHubFinding records
func (c *SecurityHubFindingExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// the expected input type is a JSON byte[] deserializable to DetailFindingsData
	var jsonBytes []byte

	switch v := a.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		return nil, fmt.Errorf("expected []byte or string, got %T", a)
	}

	// First, we need to remap certain JSON fields due to naming conventions
	// DetailFindingsData expects "detail-type" to be mapped to "detail_type"
	var rawEvent map[string]json.RawMessage
	if err := json.Unmarshal(jsonBytes, &rawEvent); err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	// Re-encode the modified JSON
	modifiedJSON, err := json.Marshal(rawEvent)
	if err != nil {
		return nil, fmt.Errorf("error re-encoding json: %w", err)
	}

	// decode json into DetailFindingsData
	var event DetailFindingsData
	err = json.Unmarshal(modifiedJSON, &event)
	if err != nil {
		slog.Debug("Error decoding SecurityHub finding", "error", err, "sample_start", string(jsonBytes[:min(len(jsonBytes), 500)]))
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	slog.Debug("SecurityHubFindingExtractor", "record count", len(event.Detail.Findings))

	findings := toMapSecurityHubFinding(event)
	var res = make([]any, len(findings))
	for i, record := range findings {
		res[i] = &record
	}
	return res, nil
}

func toMapSecurityHubFinding(event DetailFindingsData) []SecurityHubFinding {
	var findings []SecurityHubFinding

	for _, finding := range event.Detail.Findings {
		f := SecurityHubFinding{}

		// Event metadata
		f.Version = event.Version
		f.ID = event.ID
		f.DetailType = event.DetailType
		f.Source = event.Source
		f.Account = event.Account
		f.Time = event.Time
		f.Region = event.Region

		// Finding details from AWS security finding
		if finding.AwsAccountName != nil {
			f.AwsAccountName = finding.AwsAccountName
		}
		if finding.CompanyName != nil {
			f.CompanyName = finding.CompanyName
		}
		if finding.Compliance != nil {
			f.Compliance = finding.Compliance
		}
		if finding.Confidence != nil {
			f.Confidence = finding.Confidence
		}
		if finding.CreatedAt != nil {
			createdAtStr := *finding.CreatedAt
			f.CreatedAt = &createdAtStr
		}
		if finding.Criticality != nil {
			f.Criticality = finding.Criticality
		}
		if finding.Description != nil {
			f.Description = finding.Description
		}
		if finding.FirstObservedAt != nil {
			f.FirstObservedAt = finding.FirstObservedAt
		}
		if finding.GeneratorId != nil {
			f.GeneratorId = finding.GeneratorId
		}
		if finding.GeneratorDetails != nil {
			f.GeneratorDetails = finding.GeneratorDetails
		}
		if finding.Id != nil {
			f.FindingId = finding.Id
		}
		if finding.Region != nil {
			f.FindingRegion = finding.Region
		}
		if finding.LastObservedAt != nil {
			f.LastObservedAt = finding.LastObservedAt
		}
		if finding.Malware != nil {
			f.Malware = finding.Malware
		}
		if finding.Network != nil {
			f.Network = finding.Network
		}
		if finding.NetworkPath != nil {
			f.NetworkPath = finding.NetworkPath
		}
		if finding.Note != nil {
			f.Note = finding.Note
		}
		if finding.PatchSummary != nil {
			f.PatchSummary = finding.PatchSummary
		}
		if finding.Process != nil {
			f.Process = finding.Process
		}
		if finding.ProcessedAt != nil {
			f.ProcessedAt = finding.ProcessedAt
		}
		if finding.ProductArn != nil {
			f.ProductArn = finding.ProductArn
		}
		if finding.ProductName != nil {
			f.ProductName = finding.ProductName
		}
		if finding.RecordState != "" {
			f.RecordState = finding.RecordState
		}
		if finding.RelatedFindings != nil {
			f.RelatedFindings = finding.RelatedFindings
		}
		if finding.Action != nil {
			f.Action = finding.Action
		}
		if finding.Sample != nil {
			f.Sample = finding.Sample
		}
		if finding.SchemaVersion != nil {
			f.SchemaVersion = finding.SchemaVersion
		}
		if finding.Severity != nil {
			f.Severity = finding.Severity
		}
		if finding.SourceUrl != nil {
			f.SourceUrl = finding.SourceUrl
		}
		if finding.ThreatIntelIndicators != nil {
			f.ThreatIntelIndicators = finding.ThreatIntelIndicators
		}
		if finding.Threats != nil {
			f.Threats = finding.Threats
		}
		if finding.Title != nil {
			f.Title = finding.Title
		}
		if finding.Types != nil {
			f.Types = finding.Types
		}
		if finding.UpdatedAt != nil {
			f.UpdatedAt = finding.UpdatedAt
		}
		if finding.UserDefinedFields != nil {
			userDefinedFields := make(map[string]string)
			for k, v := range finding.UserDefinedFields {
				userDefinedFields[k] = v
			}
			f.UserDefinedFields = userDefinedFields
		}
		if finding.VerificationState != "" {
			f.VerificationState = finding.VerificationState
		}
		if finding.Vulnerabilities != nil {
			f.Vulnerabilities = finding.Vulnerabilities
		}
		if finding.Workflow != nil {
			f.Workflow = finding.Workflow
		}

		// Map ProductFields
		if finding.ProductFields != nil {
			productFields := make(map[string]string)
			for k, v := range finding.ProductFields {
				productFields[k] = v
			}
			f.ProductFields = productFields
		}

		// Map Resources
		if len(finding.Resources) > 0 {
			f.Resources = finding.Resources
		}

		// Map Remediation
		if finding.Remediation != nil && finding.Remediation.Recommendation != nil {
			f.Remediation = finding.Remediation
		}

		findings = append(findings, f)
	}

	return findings
}
