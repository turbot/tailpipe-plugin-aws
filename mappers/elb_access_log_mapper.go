package mappers

import (
	"context"
	"fmt"

	"github.com/satyrius/gonx"
)

const elbLogFormat = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason" $conn_trace_id`
const elbLogFormatNoConnTrace = `$type $timestamp $elb $client $target $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn "$trace_id" "$domain_name" "$chosen_cert_arn" $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason" "$target_list" "$target_status_list" "$classification" "$classification_reason"`

type ElbAccessLogMapper struct {
	fullParser   *gonx.Parser
	noConnParser *gonx.Parser
}

func NewElbAccessLogMapper() *ElbAccessLogMapper {
	return &ElbAccessLogMapper{
		fullParser:   gonx.NewParser(elbLogFormat),
		noConnParser: gonx.NewParser(elbLogFormatNoConnTrace),
	}
}

func (c *ElbAccessLogMapper) Identifier() string {
	return "elb_access_log_mapper"
}

func (c *ElbAccessLogMapper) Map(ctx context.Context, a any) ([]any, error) {
	var out []any
	var parsed *gonx.Entry
	var err error

	// validate input type is string
	input, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %T", a)
	}

	// parse log line
	parsed, err = c.fullParser.ParseString(input)
	if err != nil {
		parsed, err = c.noConnParser.ParseString(input)
		if err != nil {
			return nil, fmt.Errorf("error parsing log line: %w", err)
		}
	}

	fields := make(map[string]string)

	fields = parsed.Fields()
	out = append(out, fields)

	return out, nil
}
