package cloudtrail_log

import (
	"context"
	"encoding/json"
	"fmt"

	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

type CloudTrailMapper struct {
}

func (m *CloudTrailMapper) Identifier() string {
	return "cloudtrail_mapper"
}

func (m *CloudTrailMapper) Map(_ context.Context, a any, _ ...mappers.MapOption[*CloudTrailLog]) (*CloudTrailLog, error) {
	var log CloudTrailLog
	var jsonBytes []byte
	var err error

	switch v := a.(type) {
	case *CloudTrailLog:
		return v, nil
	case CloudTrailLog:
		return &v, nil
	case []byte:
		jsonBytes = v
	case *string:
		jsonBytes, err = m.decodeString(*v)
		if err != nil {
			return nil, fmt.Errorf("error decoding string: %w", err)
		}
	case string:
		jsonBytes, err = m.decodeString(v)
		if err != nil {
			return nil, fmt.Errorf("error decoding string: %w", err)
		}
	case cwTypes.FilteredLogEvent:
		jsonBytes = []byte(*v.Message)
	default:
		return nil, fmt.Errorf("expected byte[], string or rows.CloudTailLog got %T", a)
	}

	err = json.Unmarshal(jsonBytes, &log)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return &log, nil
}

func (m *CloudTrailMapper) decodeString(input string) ([]byte, error) {
	inputBytes := []byte(input)

	// Attempt Direct Json Unmarshalling
	var result map[string]interface{}
	err := json.Unmarshal(inputBytes, &result)
	if err == nil {
		return inputBytes, nil
	}

	// Attempt Unquoting
	var unescaped string
	err = json.Unmarshal([]byte(`"`+input+`"`), &unescaped) // Wrap the input in quotes to mimic a JSON string literal
	if err != nil {
		return nil, fmt.Errorf("failed to unescape JSON string: %w", err)
	}

	// Decode the unescaped string
	unescapedBytes := []byte(unescaped)
	err = json.Unmarshal(unescapedBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode unescaped JSON: %w", err)
	}

	return unescapedBytes, nil
}
