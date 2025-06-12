package securityhub_finding

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type SecurityHubFinding struct {
	schema.CommonFields

	// Top level fields
	Version    *string    `json:"version,omitempty"`
	ID         *string    `json:"id,omitempty"`
	DetailType *string    `json:"detail_type,omitempty"`
	Source     *string    `json:"source,omitempty"`
	Account    *string    `json:"account,omitempty"`
	Time       *time.Time `json:"time,omitempty"`
	Region     *string    `json:"region,omitempty"`

	// Finding array schema
	AwsAccountName        *string                      `json:"aws_account_name" parquet:"name=aws_account_name"`
	CompanyName           *string                      `json:"company_name" parquet:"name=company_name"`
	Compliance            *types.Compliance            `json:"compliance" parquet:"name=compliance"`
	Confidence            *int32                       `json:"confidence" parquet:"name=confidence"`
	CreatedAt             *string                      `json:"createdAt" parquet:"name=created_at"`
	Criticality           *int32                       `json:"criticality" parquet:"name=criticality"`
	Description           *string                      `json:"description" parquet:"name=description"`
	FirstObservedAt       *string                      `json:"first_observed_at" parquet:"name=first_observed_at"`
	GeneratorId           *string                      `json:"generatorId" parquet:"name=generator_id"`
	GeneratorDetails      *types.GeneratorDetails      `json:"generator_details" parquet:"name=generator_details"`
	FindingId             *string                      `json:"findingId" parquet:"name=finding_id"`
	FindingRegion         *string                      `json:"findingRegion" parquet:"name=finding_region"`
	LastObservedAt        *string                      `json:"last_observed_at" parquet:"name=last_observed_at"`
	Malware               []types.Malware              `json:"malware" parquet:"name=malware"`
	Network               *types.Network               `json:"network" parquet:"name=network"`
	NetworkPath           []types.NetworkPathComponent `json:"network_path" parquet:"name=network_path"`
	Note                  *types.Note                  `json:"note" parquet:"name=note"`
	PatchSummary          *types.PatchSummary          `json:"patch_summary" parquet:"name=patch_summary"`
	Process               *types.ProcessDetails        `json:"process" parquet:"name=process"`
	ProcessedAt           *string                      `json:"processed_at" parquet:"name=processed_at"`
	ProductArn            *string                      `json:"product_arn" parquet:"name=product_arn"`
	ProductFields         map[string]string            `json:"product_fields" parquet:"name=product_fields"`
	ProductName           *string                      `json:"product_name" parquet:"name=product_name"`
	RecordState           types.RecordState            `json:"record_state" parquet:"name=record_state"`
	RelatedFindings       []types.RelatedFinding       `json:"related_findings" parquet:"name=related_findings"`
	Remediation           *types.Remediation           `json:"remediation" parquet:"name=remediation"`
	Resources             []types.Resource             `json:"resources" parquet:"name=resources"`
	Action                *types.Action                `json:"action" parquet:"name=action"`
	Sample                *bool                        `json:"sample" parquet:"name=sample"`
	SchemaVersion         *string                      `json:"schema_version" parquet:"name=schema_version"`
	Severity              *types.Severity              `json:"severity" parquet:"name=severity"`
	SourceUrl             *string                      `json:"source_url" parquet:"name=source_url"`
	ThreatIntelIndicators []types.ThreatIntelIndicator `json:"threat_intel_indicators" parquet:"name=threat_intel_indicators"`
	Threats               []types.Threat               `json:"threats" parquet:"name=threats"`
	Title                 *string                      `json:"title" parquet:"name=title"`
	Types                 []string                     `json:"types" parquet:"name=types"`
	UpdatedAt             *string                      `json:"updated_at" parquet:"name=updated_at"`
	UserDefinedFields     map[string]string            `json:"user_defined_fields" parquet:"name=user_defined_fields"`
	VerificationState     types.VerificationState      `json:"verification_state" parquet:"name=verification_state"`
	Vulnerabilities       []types.Vulnerability        `json:"vulnerabilities" parquet:"name=vulnerabilities"`
	Workflow              *types.Workflow              `json:"workflow" parquet:"name=workflow"`
	WorkflowState         types.WorkflowState          `json:"workflow_state" parquet:"name=workflow_state"`
}

// DetailFindingsData maps the `detail` field containing findings
// The following struct will be used for only parse the log lines
type DetailFindingsData struct {
	Version    *string    `json:"version,omitempty"`
	ID         *string    `json:"id,omitempty"`
	DetailType *string    `json:"detail-type,omitempty"`
	Source     *string    `json:"source,omitempty"`
	Account    *string    `json:"account,omitempty"`
	Time       *time.Time `json:"time,omitempty"`
	Region     *string    `json:"region,omitempty"`
	Detail     struct {
		Findings []types.AwsSecurityFinding `json:"findings" parquet:"name=findings, type=JSON"`
	} `json:"detail" parquet:"name=detail, type=JSON"`
}

func (c *SecurityHubFinding) GetColumnDescriptions() map[string]string {
	return map[string]string{
		// Top level fields
		"version":     "The version of the event format.",
		"id":          "The unique identifier for the event.",
		"detail_type": "The type of the event detail.",
		"source":      "The service or system that generated the event.",
		"account":     "The AWS account ID where the finding was generated.",
		"time":        "The timestamp when the event was generated.",
		"region":      "The AWS region where the finding was generated.",

		// Finding fields
		"aws_account_name":        "The name of the AWS account from which a finding was generated.",
		"company_name":            "The name of the company for the product that generated the finding. Security Hub populates this attribute automatically for each finding.",
		"compliance":              "Contains security standard-related finding details for findings generated from compliance checks against specific rules in supported security standards.",
		"confidence":              "The likelihood that a finding accurately identifies the behavior or issue that it was intended to identify. Scored on a 0-100 basis.",
		"created_at":              "The timestamp when the security findings provider created the potential security issue that a finding captured.",
		"criticality":             "The level of importance assigned to the resources associated with the finding. A score of 0 means no criticality, and 100 is reserved for the most critical resources.",
		"description":             "A detailed description of the security finding.",
		"first_observed_at":       "The timestamp when the security findings provider first observed the potential security issue that a finding captured.",
		"generator_id":            "The identifier for the solution-specific component (a discrete unit of logic) that generated a finding.",
		"generator_details":       "Metadata for the Amazon CodeGuru detector associated with a finding, particularly for Lambda function-related findings.",
		"finding_id":              "The security findings provider-specific identifier for a finding.",
		"finding_region":          "The AWS region from which the finding was generated.",
		"last_observed_at":        "The timestamp when the security findings provider most recently observed the potential security issue that a finding captured.",
		"malware":                 "A list of malware related to a finding.",
		"network":                 "The details of network-related information about a finding.",
		"network_path":            "Information about a network path that is relevant to a finding, with each entry representing a component of that path.",
		"note":                    "A user-defined note added to a finding.",
		"patch_summary":           "An overview of the patch compliance status for an instance against a selected compliance standard.",
		"process":                 "The details of process-related information about a finding.",
		"processed_at":            "The timestamp when Security Hub received a finding and began to process it.",
		"product_arn":             "The ARN generated by Security Hub that uniquely identifies a product that generates findings.",
		"product_fields":          "Additional solution-specific details that aren't part of the defined AwsSecurityFinding format. Can contain up to 50 key-value pairs.",
		"product_name":            "The name of the product that generated the finding. Security Hub populates this attribute automatically for each finding.",
		"record_state":            "The record state of a finding.",
		"related_findings":        "A list of related findings.",
		"remediation":             "A data type that describes the remediation options for a finding.",
		"resources":               "A set of resource data types that describe the resources that the finding refers to.",
		"action":                  "Details about an action that affects or that was taken on a resource.",
		"sample":                  "Indicates whether the finding is a sample finding.",
		"schema_version":          "The schema version that a finding is formatted for.",
		"severity":                "The severity level of the finding.",
		"source_url":              "A URL that links to a page about the current finding in the security findings provider's solution.",
		"threat_intel_indicators": "Threat intelligence details related to a finding.",
		"threats":                 "Details about the threat detected in a security finding and the file paths that were affected by the threat.",
		"title":                   "A short human-readable title for the finding.",
		"types":                   "One or more finding types in the format of namespace/category/classifier that classify a finding.",
		"updated_at":              "The timestamp when the security findings provider last updated the finding record.",
		"user_defined_fields":     "A list of name/value string pairs associated with the finding. These are custom, user-defined fields added to a finding.",
		"verification_state":      "Indicates the veracity of a finding.",
		"vulnerabilities":         "A list of vulnerabilities associated with the findings.",
		"workflow":                "Information about the status of the investigation into a finding.",
		"workflow_state":          "The workflow state of a finding.",

		// Tailpipe-specific metadata fields
		"tp_akas":      "The list of AWS ARNs associated with the finding.",
		"tp_timestamp": "The timestamp when the finding was generated.",
		"tp_date":      "The date when the finding was generated, truncated to day.",
	}
}
