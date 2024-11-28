package rows

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

// Define a nested struct for httpRequest headers
type Header struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Define a nested struct for the httpRequest field
type HttpRequest struct {
	ClientIp    *string   `json:"client_ip,omitempty"`
	Country     *string   `json:"country,omitempty"`
	Headers     *[]Header `json:"headers,omitempty"`
	Uri         *string   `json:"uri,omitempty"`
	Args        *string   `json:"args,omitempty"`
	HttpVersion *string   `json:"http_version,omitempty"`
	HttpMethod  *string   `json:"http_method,omitempty"`
	RequestId   *string   `json:"request_id,omitempty"`
}

// WafTrafficLog struct with fields aligned to the provided JSON
type WafTrafficLog struct {
	enrichment.CommonFields

	Timestamp                   *time.Time   `json:"timestamp,omitempty"`
	FormatVersion               *int32       `json:"format_version,omitempty"`
	WebAclId                    *string      `json:"web_acl_id,omitempty"`
	TerminatingRuleId           *string      `json:"terminating_rule_id,omitempty"`
	TerminatingRuleType         *string      `json:"terminating_rule_type,omitempty"`
	Action                      *string      `json:"action,omitempty"`
	HttpSourceName              *string      `json:"http_source_name,omitempty"`
	HttpSourceId                *string      `json:"http_source_id,omitempty"`
	RuleGroupList               []string     `json:"rule_group_list,omitempty"`
	RateBasedRuleList           []string     `json:"rate_based_rule_list,omitempty"`
	NonTerminatingMatchingRules []string     `json:"non_terminating_matching_rules,omitempty"`
	HttpRequest                 *HttpRequest `json:"http_request,omitempty"`
}
