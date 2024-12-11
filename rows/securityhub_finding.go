package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type SecurityHubFinding struct {
	enrichment.CommonFields

	// Top level fields
	Version        *string             `json:"version,omitempty"`
	ID             *string             `json:"id,omitempty"`
	DetailType     *string             `json:"detail_type,omitempty"`
	Source         *string             `json:"source,omitempty"`
	Account        *string             `json:"account,omitempty"`
	Time           *time.Time          `json:"time,omitempty"`
	Region         *string             `json:"region,omitempty"`
	Resources      []*string           `json:"resources,omitempty"`
	DetailFindings *DetailFindingsData `json:"detail,omitempty"` // Updated to map the detailed findings
}

// DetailFindingsData maps the `detail` field containing findings
type DetailFindingsData struct {
	Findings []*Finding `json:"findings,omitempty" parquet:"name=findings, type=JSON"`
}

// Finding maps the individual findings in the detail
type Finding struct {
	ProductArn            *string                `json:"ProductArn,omitempty"`
	Types                 []*string              `json:"Types,omitempty"`
	Description           *string                `json:"Description,omitempty"`
	Compliance            *Compliance            `json:"Compliance,omitempty"`
	ProductName           *string                `json:"ProductName,omitempty"`
	FirstObservedAt       *time.Time             `json:"FirstObservedAt,omitempty"`
	CreatedAt             *time.Time             `json:"CreatedAt,omitempty"`
	LastObservedAt        *time.Time             `json:"LastObservedAt,omitempty"`
	CompanyName           *string                `json:"CompanyName,omitempty"`
	FindingProviderFields *FindingProviderFields `json:"FindingProviderFields,omitempty"`
	ProductFields         map[string]string      `json:"ProductFields,omitempty"`
	Remediation           *Remediation           `json:"Remediation,omitempty"`
	SchemaVersion         *string                `json:"SchemaVersion,omitempty"`
	GeneratorId           *string                `json:"GeneratorId,omitempty"`
	RecordState           *string                `json:"RecordState,omitempty"`
	Title                 *string                `json:"Title,omitempty"`
	Workflow              *Workflow              `json:"Workflow,omitempty"`
	Severity              *Severity              `json:"Severity,omitempty"`
	UpdatedAt             *time.Time             `json:"UpdatedAt,omitempty"`
	WorkflowState         *string                `json:"WorkflowState,omitempty"`
	AwsAccountId          *string                `json:"AwsAccountId,omitempty"`
	Region                *string                `json:"Region,omitempty"`
	Id                    *string                `json:"Id,omitempty"`
	Resources             []*FindingResource     `json:"Resources,omitempty"`
	ProcessedAt           *time.Time             `json:"ProcessedAt,omitempty"`
}

// Supporting structs for nested fields
type Compliance struct {
	Status                    *string                     `json:"Status,omitempty"`
	SecurityControlId         *string                     `json:"SecurityControlId,omitempty"`
	AssociatedStandards       []*AssociatedStandard       `json:"AssociatedStandards,omitempty"`
	SecurityControlParameters []*SecurityControlParameter `json:"SecurityControlParameters,omitempty"`
}

type AssociatedStandard struct {
	StandardsId *string `json:"StandardsId,omitempty"`
}

type SecurityControlParameter struct {
	Value []string `json:"Value,omitempty"`
	Name  *string  `json:"Name,omitempty"`
}

type FindingProviderFields struct {
	Types    []*string `json:"Types,omitempty"`
	Severity *Severity `json:"Severity,omitempty"`
}

type Remediation struct {
	Recommendation *Recommendation `json:"Recommendation,omitempty"`
}

type Recommendation struct {
	Text *string `json:"Text,omitempty"`
	Url  *string `json:"Url,omitempty"`
}

type Workflow struct {
	Status *string `json:"Status,omitempty"`
}

type Severity struct {
	Normalized *int    `json:"Normalized,omitempty"`
	Label      *string `json:"Label,omitempty"`
	Original   *string `json:"Original,omitempty"`
}

type FindingResource struct {
	Partition *string          `json:"Partition,omitempty"`
	Type      *string          `json:"Type,omitempty"`
	Details   *ResourceDetails `json:"Details,omitempty"`
	Region    *string          `json:"Region,omitempty"`
	Id        *string          `json:"Id,omitempty"`
}

type ResourceDetails struct {
	AwsLambdaFunction *AwsLambdaFunction `json:"AwsLambdaFunction,omitempty"`
}

type AwsLambdaFunction struct {
	LastModified  *time.Time     `json:"LastModified,omitempty"`
	Role          *string        `json:"Role,omitempty"`
	FunctionName  *string        `json:"FunctionName,omitempty"`
	MemorySize    *int           `json:"MemorySize,omitempty"`
	Runtime       *string        `json:"Runtime,omitempty"`
	TracingConfig *TracingConfig `json:"TracingConfig,omitempty"`
	Version       *string        `json:"Version,omitempty"`
	Timeout       *int           `json:"Timeout,omitempty"`
	Handler       *string        `json:"Handler,omitempty"`
	CodeSha256    *string        `json:"CodeSha256,omitempty"`
	RevisionId    *string        `json:"RevisionId,omitempty"`
}

type TracingConfig struct {
	Mode *string `json:"Mode,omitempty"`
}
