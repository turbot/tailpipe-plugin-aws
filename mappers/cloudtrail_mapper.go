package mappers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type CloudTrailMapper struct {
}

func NewCloudTrailMapper() table.Mapper[*rows.CloudTrailLog] {
	return &CloudTrailMapper{}
}

func (m *CloudTrailMapper) Identifier() string {
	return "cloudtrail_mapper"
}

func (m *CloudTrailMapper) Map(_ context.Context, a any) (*rows.CloudTrailLog, error) {
	var log rows.CloudTrailLog
	var jsonBytes []byte
	var err error

	switch v := a.(type) {
	case *rows.CloudTrailLog:
		return v, nil
	case rows.CloudTrailLog:
		return &v, nil
	case []byte:
		jsonBytes = v
	case string:
		if strings.Contains(v, `\`) {
			v, err = strconv.Unquote(v)
			if err != nil {
				return nil, fmt.Errorf("error unquoting string: %w", err)
			}
		}
		jsonBytes = []byte(v)
	default:
		return nil, fmt.Errorf("expected byte[], string or rows.CloudTailLog got %T", a)
	}

	err = json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return &log, nil
}
