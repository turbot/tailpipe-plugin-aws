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
	// Detail     *DetailFindingsData `json:"detail,omitempty" parquet:"name=detail, type=JSON"`
	// Finding array schema
	AwsAccountId  *string            `json:"awsAccountId" parquet:"name=aws_account_id"`
	CreatedAt     *string            `json:"createdAt" parquet:"name=created_at"`
	Description   *string            `json:"description" parquet:"name=description"`
	GeneratorId   *string            `json:"generatorId" parquet:"name=generator_id"`
	FindingId     *string            `json:"findingId" parquet:"name=finding_id"`
	ProductArn    *string            `json:"productArn" parquet:"name=product_arn"`
	ProductFields map[string]string  `json:"productFields" parquet:"name=product_fields, type=JSON"`
	ProductName   *string            `json:"productName" parquet:"name=product_name"`
	Remediation   *types.Remediation `json:"remediation" parquet:"name=remediation"`
	Resources     []types.Resource   `json:"resources" parquet:"name=resources, type=JSON"`
	SchemaVersion *string            `json:"schemaVersion" parquet:"name=schema_version"`
	Title         *string            `json:"title" parquet:"name=title"`

	// It is for schema only
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

// DetailFindingsData maps the `detail` field containing findings
// The following struct will be used for only parse the log lines
type DetailFindingsData struct {
	Version    *string    `json:"version,omitempty"`
	ID         *string    `json:"id,omitempty"`
	DetailType *string    `json:"detail_type,omitempty"`
	Source     *string    `json:"source,omitempty"`
	Account    *string    `json:"account,omitempty"`
	Time       *time.Time `json:"time,omitempty"`
	Region     *string    `json:"region,omitempty"`
	Detail     struct {
		Findings           []types.AwsSecurityFinding `json:"findings" parquet:"name=findings, type=JSON"`
		AwsRegion          *string                    `json:"awsRegion" parquet:"name=aws_region"`
		EventCategory      *string                    `json:"eventCategory" parquet:"name=event_category"`
		EventID            *string                    `json:"eventID" parquet:"name=event_id"`
		EventName          *string                    `json:"eventName" parquet:"name=event_name"`
		EventSource        *string                    `json:"eventSource" parquet:"name=event_source"`
		EventTime          time.Time                  `json:"eventTime" parquet:"name=event_time"`
		EventType          *string                    `json:"eventType" parquet:"name=event_type"`
		EventVersion       *string                    `json:"eventVersion" parquet:"name=event_version"`
		ManagementEvent    bool                       `json:"managementEvent" parquet:"name=management_event"`
		ReadOnly           bool                       `json:"readOnly" parquet:"name=read_only"`
		RecipientAccountID *string                    `json:"recipientAccountId" parquet:"name=recipient_account_id"`
		RequestID          *string                    `json:"requestID" parquet:"name=request_id"`
		RequestParameters  *map[string]interface{}    `json:"requestParameters" parquet:"name=request_parameters, type=JSON"`
		ResponseElements   *map[string]interface{}    `json:"responseElements" parquet:"name=response_elements, type=JSON"`
		SourceIPAddress    *string                    `json:"sourceIPAddress" parquet:"name=source_ip_address"`
		UserAgent          *string                    `json:"userAgent" parquet:"name=user_agent"`
		UserIdentity       SecurityHubUserIdentity    `json:"userIdentity" parquet:"name=user_identity, type=JSON"`
	} `json:"detail" parquet:"name=detail, type=JSON"`
}

type AssociatedStandard struct {
	StandardsId *string `json:"StandardsId,omitempty"`
}

type SecurityControlParameter struct {
	Value []string `json:"Value,omitempty"`
	Name  *string  `json:"Name,omitempty"`
}

type Recommendation struct {
	Text *string `json:"Text,omitempty"`
	Url  *string `json:"Url,omitempty"`
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
		"resources":   "The list of AWS resources associated with the finding.",

		// Detail fields
		"detail": "The detailed information about the security finding.",

		// Finding fields
		"product_arn":       "The ARN of the AWS security product that generated the finding.",
		"types":             "The list of types assigned to the finding.",
		"description":       "A detailed description of the security finding.",
		"compliance":        "Information about the finding's compliance status.",
		"product_name":      "The name of the security product that generated the finding.",
		"first_observed_at": "The timestamp when the finding was first observed.",
		"created_at":        "The timestamp when the finding was created.",
		"last_observed_at":  "The timestamp when the finding was last observed.",
		"company_name":      "The name of the company that provides the security product.",
		"product_fields":    "Additional fields provided by the security product.",
		"remediation":       "Recommended steps to remediate the finding.",
		"schema_version":    "The version of the finding format schema.",
		"generator_id":      "The identifier of the system that generated the finding.",
		"record_state":      "The current state of the finding record.",
		"title":             "A short human-readable title for the finding.",
		"workflow":          "Information about the finding's workflow status.",
		"severity":          "The severity level of the finding.",
		"updated_at":        "The timestamp when the finding was last updated.",
		"workflow_state":    "The current state of the finding in the workflow.",
		"aws_account_id":    "The AWS account ID associated with the finding.",
		"processed_at":      "The timestamp when the finding was processed.",

		// Tailpipe-specific metadata fields
		"tp_akas":      "The list of AWS ARNs associated with the finding.",
		"tp_index":     "The AWS account ID where the finding was generated.",
		"tp_timestamp": "The timestamp when the finding was generated.",
		"tp_date":      "The date when the finding was generated, truncated to day.",
	}
}
