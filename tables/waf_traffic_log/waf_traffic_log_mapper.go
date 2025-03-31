package waf_traffic_log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	helper "github.com/turbot/tailpipe-plugin-aws/tables"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type WafMapper struct {
}

func (c *WafMapper) Identifier() string {
	return "aws_waf_traffic_log_mapper"
}

func (c *WafMapper) Map(_ context.Context, a any, _ ...table.MapOption[*WafTrafficLog]) (*WafTrafficLog, error) {
	var jsonBytes []byte

	switch v := a.(type) {
	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	case *string:
		jsonBytes = []byte(*v)
	default:
		return nil, fmt.Errorf("expected byte[] or string, got %T", a)
	}

	// decode JSON into WafTrafficLog
	var log WafTrafficLog
	err := unmarshalWafTrafficLog(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w; partial log: %+v", err, log)
	}

	return &log, nil
}

func unmarshalWafTrafficLog(data []byte, log *WafTrafficLog) error {
	var temp struct {
		Action                      *string                   `json:"action"`
		CaptchaResponse             *CaptchaResponse          `json:"captchaResponse"`
		FormatVersion               *int32                    `json:"formatVersion"`
		HttpRequest                 *HttpRequest              `json:"httpRequest,omitempty"`
		HttpSourceId                *string                   `json:"httpSourceId,omitempty"`
		HttpSourceName              *string                   `json:"httpSourceName,omitempty"`
		Labels                      []Labels                  `json:"labels,omitempty"`
		NonTerminatingMatchingRules []NonTerminatingRuleMatch `json:"nonTerminatingMatchingRules,omitempty"`
		RateBasedRuleList           []RateBasedRule           `json:"rateBasedRuleList,omitempty"`
		RequestHeadersInserted      []Header                  `json:"requestHeadersInserted,omitempty"`
		RuleGroupList               []RuleGroup               `json:"ruleGroupList,omitempty"`
		TerminatingRuleId           *string                   `json:"terminatingRuleId,omitempty"`
		TerminatingRuleMatchDetails []TerminatingRuleMatch    `json:"terminatingRuleMatchDetails,omitempty"`
		TerminatingRuleType         *string                   `json:"terminatingRuleType,omitempty"`
		Timestamp                   *int64                    `json:"timestamp"`
		WebAclId                    *string                   `json:"webAclId"`
	}

	// Unmarshal JSON into temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Assign values from temp struct
	log.FormatVersion = temp.FormatVersion
	log.WebAclId = temp.WebAclId
	log.TerminatingRuleMatchDetails = temp.TerminatingRuleMatchDetails
	log.TerminatingRuleId = temp.TerminatingRuleId
	log.TerminatingRuleType = temp.TerminatingRuleType
	log.Action = temp.Action
	log.RuleGroupList = temp.RuleGroupList
	log.RateBasedRuleList = temp.RateBasedRuleList
	log.NonTerminatingMatchingRules = temp.NonTerminatingMatchingRules
	log.HttpRequest = temp.HttpRequest
	log.RequestHeadersInserted = temp.RequestHeadersInserted
	log.Labels = temp.Labels
	log.CaptchaResponse = temp.CaptchaResponse

	// For a rule that triggered on SQLi detection(terminating/non-terminating) will not have HttpSourceName and HttpSourceId.
	log.HttpSourceName = helper.NilIfDash(temp.HttpSourceName)
	log.HttpSourceId = helper.NilIfDash(temp.HttpSourceId)

	// Convert timestamp (if exists) to *time.Time
	if temp.Timestamp != nil {
		parsedTime := time.UnixMilli(*temp.Timestamp).UTC()
		log.Timestamp = &parsedTime
	}

	return nil
}
