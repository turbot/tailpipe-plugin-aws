package mappers

import (
	"context"
	"fmt"

	"github.com/satyrius/gonx"
)

const s3ServerAccessLogFormat = `$bucket_owner $bucket [$timestamp] $remote_ip $requester $request_id $operation $key "$request_uri" $http_status $error_code $bytes_sent $object_size $total_time $turn_around_time "$referer" "$user_agent" $version_id $host_id $signature_version $cipher_suite $authentication_type $host_header $tls_version $access_point_arn $acl_required`
const s3ServerAccessLogFormatReduced = `$bucket_owner $bucket [$timestamp] $remote_ip $requester $request_id $operation $key "$request_uri" $http_status $error_code $bytes_sent $object_size $total_time $turn_around_time "$referer" "$user_agent" $version_id $host_id $signature_version $cipher_suite $authentication_type $host_header $tls_version $access_point_arn`

type S3ServerAccessLogMapper struct {
	fullParser    *gonx.Parser
	reducedParser *gonx.Parser
}

func NewS3ServerAccessLogMapper() *S3ServerAccessLogMapper {
	return &S3ServerAccessLogMapper{
		fullParser:    gonx.NewParser(s3ServerAccessLogFormat),
		reducedParser: gonx.NewParser(s3ServerAccessLogFormatReduced),
	}
}

func (c *S3ServerAccessLogMapper) Identifier() string {
	return "s3_server_access_log_mapper"
}

func (c *S3ServerAccessLogMapper) Map(ctx context.Context, a any) ([]any, error) {
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
		parsed, err = c.reducedParser.ParseString(input)
		if err != nil {
			return nil, fmt.Errorf("error parsing log line: %w", err)
		}
	}

	fields := make(map[string]string)
	fields = parsed.Fields()
	out = append(out, fields)

	return out, nil
}
