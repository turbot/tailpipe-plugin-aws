package waf_traffic_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

// Define a nested struct for httpRequest headers
type Header struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Define a nested struct for the httpRequest field
type HttpRequest struct {
	ClientIp    *string   `json:"clientIp,omitempty"`
	Country     *string   `json:"country,omitempty"`
	Headers     *[]Header `json:"headers,omitempty" parquet:"type=JSON"`
	Uri         *string   `json:"uri,omitempty"`
	Args        *string   `json:"args,omitempty"`
	HttpVersion *string   `json:"httpVersion,omitempty"`
	HttpMethod  *string   `json:"httpMethod,omitempty"`
	RequestId   *string   `json:"requestId,omitempty"`
}

type TerminatingRuleMatch struct {
	ConditionType    *string  `json:"conditionType"`
	SensitivityLevel *string  `json:"sensitivityLevel"`
	Location         *string  `json:"location"`
	MatchedData      []string `json:"matchedData"`
}

type CustomValue struct {
	Key   *string `json:"key"`
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

// RateBasedRule represents the main JSON structure
type RateBasedRule struct {
	RateBasedRuleID     *interface{}  `json:"rateBasedRuleId"` // Assuming ID could be a string or number
	RateBasedRuleName   *string       `json:"rateBasedRuleName"`
	LimitKey            *string       `json:"limitKey"`
	MaxRateAllowed      *int          `json:"maxRateAllowed"`
	EvaluationWindowSec *int32        `json:"evaluationWindowSec"`
	CustomValues        []CustomValue `json:"customValues" parquet:"type=JSON"`
}

type RuleMatchDetail struct {
	ConditionType    *string  `json:"conditionType"`
	SensitivityLevel *string  `json:"sensitivityLevel"`
	Location         *string  `json:"location"`
	MatchedData      []string `json:"matchedData"`
}

// Rule represents a rule entry
type Rule struct {
	RuleID           *string           `json:"ruleId"`
	Action           *string           `json:"action"`
	RuleMatchDetails []RuleMatchDetail `json:"ruleMatchDetails,omitempty"`
	CaptchaResponse  CaptchaResponse   `json:"captchaResponse,omitempty"`
}

type CaptchaResponse struct {
	ResponseCode   *int    `json:"responseCode,omitempty"`
	SolveTimestamp *int64  `json:"solveTimestamp,omitempty"`
	FailureReason  *string `json:"failureReason,omitempty"`
}

// RuleGroup represents the main JSON structure
type RuleGroup struct {
	RuleGroupID                 string `json:"ruleGroupId"`
	TerminatingRule             *Rule  `json:"terminatingRule,omitempty"` // Can be null
	NonTerminatingMatchingRules []Rule `json:"nonTerminatingMatchingRules"`
	ExcludedRules               []Rule `json:"excludedRules,omitempty"` // Can be null
}

type Labels struct {
	Name *string `json:"labels,omitempty"`
}

// WafTrafficLog struct with fields aligned to the provided JSON
type WafTrafficLog struct {
	schema.CommonFields

	Action                      *string                `json:"action"`
	CaptchaResponse             *CaptchaResponse       `json:"captchaResponse,omitempty" parquet:"name=captcha_response"`
	FormatVersion               *int32                 `json:"formatVersion" parquet:"name=format_version"`
	HttpRequest                 *HttpRequest           `json:"httpRequest,omitempty" parquet:"name=http_request"`
	HttpSourceId                *string                `json:"httpSourceId,omitempty" parquet:"name=http_source_id"`
	HttpSourceName              *string                `json:"httpSourceName,omitempty" parquet:"name=http_source_name"`
	Labels                      []Labels               `json:"labels,omitempty" parquet:"type=JSON"`
	NonTerminatingMatchingRules []Rule                 `json:"nonTerminatingMatchingRules,omitempty" parquet:"name=non_terminating_matching_rules, type=JSON"`
	RateBasedRuleList           []RateBasedRule        `json:"rateBasedRuleList,omitempty" parquet:"name=rate_based_rule_list, type=JSON"`
	RequestHeadersInserted      []Header               `json:"requestHeadersInserted,omitempty" parquet:"name=request_headers_inserted, type=JSON"`
	RuleGroupList               []RuleGroup            `json:"ruleGroupList,omitempty" parquet:"name=rule_group_list, type=JSON"`
	TerminatingRuleId           *string                `json:"terminatingRuleId,omitempty" parquet:"name=terminating_rule_id"`
	TerminatingRuleMatchDetails []TerminatingRuleMatch `json:"terminatingRuleMatchDetails,omitempty" parquet:"name=terminating_rule_match_details, type=JSON"`
	TerminatingRuleType         *string                `json:"terminatingRuleType,omitempty" parquet:"name=terminating_rule_type"`
	Timestamp                   *time.Time             `json:"timestamp"`
	WebAclId                    *string                `json:"webAclId" parquet:"name=web_acl_id"`
}

func (c *WafTrafficLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"action":                         "The terminating action that AWS WAF applied to the request. This indicates either allow, block, CAPTCHA, or challenge. The CAPTCHA and Challenge actions are terminating when the web request doesn't contain a valid token.",
		"captcha_response":               "The CAPTCHA action status for the request, populated when a CAPTCHA action is applied to the request. This field is populated for any CAPTCHA action, whether terminating or non-terminating. If a request has the CAPTCHA action applied multiple times, this field is populated from the last time the action was applied.",
		"format_version":                 "The format version for the log.",
		"http_request":                   "The metadata about the request.",
		"http_source_id":                 "The ID of the associated resource.",
		"http_source_name":               "The source of the request. Possible values: CF for Amazon CloudFront, APIGW for Amazon API Gateway, ALB for Application Load Balancer, APPSYNC for AWS AppSync, COGNITOIDP for Amazon Cognito, APPRUNNER for App Runner, and VERIFIED_ACCESS for Verified Access.",
		"labels":                         "The labels on the web request. These labels were applied by rules that were used to evaluate the request. AWS WAF logs the first 100 labels.",
		"non_terminating_matching_rules": "The list of non-terminating rules that matched the request. Each item in the list contains action, ruleId, and ruleMatchDetails.",
		"rate_based_rule_list":           "The list of rate-based rules that acted on the request. For information about rate-based rules, see Using rate-based rule statements in AWS WAF.",
		"request_headers_inserted":       "The list of headers inserted for custom request handling.",
		"rule_group_list":                "The list of rule groups that acted on this request, with match information.",
		"terminating_rule_id":            "The ID of the rule that terminated the request. If nothing terminates the request, the value is Default_Action.",
		"terminating_rule_match_details": "Detailed information about the terminating rule that matched the request. A terminating rule has an action that ends the inspection process against a web request. Possible actions for a terminating rule include Allow, Block, CAPTCHA, and Challenge. During the inspection of a web request, at the first rule that matches the request and that has a terminating action, AWS WAF stops the inspection and applies the action. The web request might contain other threats, in addition to the one that's reported in the log for the matching terminating rule.",
		"terminating_rule_type":          "The type of rule that terminated the request. Possible values: RATE_BASED, REGULAR, GROUP, and MANAGED_RULE_GROUP.",
		"timestamp":                      "The date and time when the request was made, in ISO 8601 format.",
		"web_acl_id":                     "The GUID of the web ACL.",

		// Override table specific tp_* column descriptions
		"tp_akas":      "List of ARNs (Amazon Resource Names) associated with the event, if applicable.",
		"tp_index":     "The AWS account ID that processed or received the request.",
		"tp_ips":       "IP addresses related to the request, including the source (client) IP and any intermediary addresses.",
		"tp_timestamp": "The timestamp when the request was made, formatted in ISO 8601 (UTC).",
	}
}
