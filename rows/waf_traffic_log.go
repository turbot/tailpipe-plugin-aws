package rows

import (
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"time"
)

// Define a nested struct for httpRequest headers
type Header struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Define a nested struct for the httpRequest field
type HttpRequest struct {
	ClientIp    *string  `json:"clientIp,omitempty"`
	Country     *string  `json:"country,omitempty"`
	Headers     []Header `json:"headers,omitempty"`
	Uri         *string  `json:"uri,omitempty"`
	Args        *string  `json:"args,omitempty"`
	HttpVersion *string  `json:"httpVersion,omitempty"`
	HttpMethod  *string  `json:"httpMethod,omitempty"`
	RequestId   *string  `json:"requestId,omitempty"`
}

// VpcFlowLog struct with fields aligned to the provided JSON
type WafTrafficLog struct {
	enrichment.CommonFields

	Timestamp          *time.Time `json:"timestamp,omitempty"`
	FormatVersion      *int32     `json:"format_version,omitempty"`
	WebAclId           *string    `json:"web_acl_id,omitempty"`
	TerminatingRuleId  *string    `json:"terminating_rule_id,omitempty"`
	TerminatingRuleType *string    `json:"terminating_rule_type,omitempty"`
	Action             *string    `json:"action,omitempty"`
	HttpSourceName     *string    `json:"httpSourceName,omitempty"`
	HttpSourceId       *string    `json:"httpSourceId,omitempty"`
	RuleGroupList      []string   `json:"ruleGroupList,omitempty"`
	RateBasedRuleList  []string   `json:"rateBasedRuleList,omitempty"`
	NonTerminatingMatchingRules []string `json:"nonTerminatingMatchingRules,omitempty"`
	HttpRequest        HttpRequest `json:"httpRequest,omitempty"`
}
