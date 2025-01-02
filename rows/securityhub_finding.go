package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type SecurityHubFinding struct {
	schema.CommonFields

	// Top level fields
	Version    *string             `json:"version,omitempty"`
	ID         *string             `json:"id,omitempty"`
	DetailType *string             `json:"detail_type,omitempty"`
	Source     *string             `json:"source,omitempty"`
	Account    *string             `json:"account,omitempty"`
	Time       *time.Time          `json:"time,omitempty"`
	Region     *string             `json:"region,omitempty"`
	Resources  []*string           `json:"resources,omitempty"`
	Detail     *DetailFindingsData `json:"detail,omitempty"`
}

// DetailFindingsData maps the `detail` field containing findings
type DetailFindingsData struct {
	Findings           []*Finding              `json:"findings,omitempty" parquet:"name=findings, type=JSON"`
	AwsRegion          *string                 `json:"awsRegion" parquet:"name=aws_region"`
	EventCategory      *string                 `json:"eventCategory" parquet:"name=event_category"`
	EventID            *string                 `json:"eventID" parquet:"name=event_id"`
	EventName          *string                 `json:"eventName" parquet:"name=event_name"`
	EventSource        *string                 `json:"eventSource" parquet:"name=event_source"`
	EventTime          time.Time               `json:"eventTime" parquet:"name=event_time"`
	EventType          *string                 `json:"eventType" parquet:"name=event_type"`
	EventVersion       *string                 `json:"eventVersion" parquet:"name=event_version"`
	ManagementEvent    bool                    `json:"managementEvent" parquet:"name=management_event"`
	ReadOnly           bool                    `json:"readOnly" parquet:"name=read_only"`
	RecipientAccountID *string                 `json:"recipientAccountId" parquet:"name=recipient_account_id"`
	RequestID          *string                 `json:"requestID" parquet:"name=request_id"`
	RequestParameters  *map[string]interface{} `json:"requestParameters" parquet:"name=request_parameters, type=JSON"`
	ResponseElements   *map[string]interface{} `json:"responseElements" parquet:"name=response_elements, type=JSON"`
	SourceIPAddress    *string                 `json:"sourceIPAddress" parquet:"name=source_ip_address"`
	UserAgent          *string                 `json:"userAgent" parquet:"name=user_agent"`
	UserIdentity       SecurityHubUserIdentity `json:"userIdentity" parquet:"name=user_identity, type=JSON"`
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

type SecurityHubUserIdentity struct {
	AccessKeyID    *string                   `json:"accessKeyId"`
	AccountID      *string                   `json:"accountId"`
	Arn            *string                   `json:"arn"`
	PrincipalID    *string                   `json:"principalId"`
	SessionContext SecurityHubSessionContext `json:"sessionContext"`
	Type           *string                   `json:"type"`
}

type SecurityHubSessionContext struct {
	Attributes    Attributes               `json:"attributes"`
	SessionIssuer SecurityHubSessionIssuer `json:"sessionIssuer"`
}

type Attributes struct {
	CreationDate     time.Time `json:"creationDate"`
	MfaAuthenticated *string   `json:"mfaAuthenticated"`
}

type SecurityHubSessionIssuer struct {
	AccountID   string `json:"accountId"`
	Arn         string `json:"arn"`
	PrincipalID string `json:"principalId"`
	Type        string `json:"type"`
	UserName    string `json:"userName"`
}
