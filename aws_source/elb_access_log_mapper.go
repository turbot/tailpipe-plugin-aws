package aws_source

import (
	"context"
	"fmt"

	"github.com/satyrius/gonx"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const elbLogFormat = `$type $timestamp $elb $client_ip:$client_port $target_ip:$target_port $request_processing_time $target_processing_time $response_processing_time $elb_status_code $target_status_code $received_bytes $sent_bytes "$request" "$user_agent" $ssl_cipher $ssl_protocol $target_group_arn $trace_id $domain_name $chosen_cert_arn $matched_rule_priority $request_creation_time "$actions_executed" "$redirect_url" "$error_reason"`

type ELBAccessLogMapper struct {
}

func NewELBAccessLogMapper() *ELBAccessLogMapper {
	return &ELBAccessLogMapper{}
}

func (c *ELBAccessLogMapper) Identifier() string {
	return "elb_access_log_mapper"
}

func (c *ELBAccessLogMapper) Map(ctx context.Context, a *types.RowData) ([]*types.RowData, error) {
	var out []*types.RowData

	// validate input type is string
	input, ok := a.Data.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %T", a.Data)
	}
	inputMetadata := a.Metadata

	// parse log line
	parser := gonx.NewParser(elbLogFormat)
	parsed, err := parser.ParseString(input)
	if err != nil {
		return nil, fmt.Errorf("error parsing log line: %w", err)
	}

	fields := make(map[string]string)

	fields = parsed.Fields()
	out = append(out, types.NewData(fields, types.WithMetadata(inputMetadata)))

	return out, nil
}
