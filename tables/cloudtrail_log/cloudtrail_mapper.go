package cloudtrail_log

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

type CloudTrailMapper struct{}

func (m CloudTrailMapper) Identifier() string {
	return "cloudtrail_log_mapper"
}

func (m CloudTrailMapper) Map(ctx context.Context, data any, options ...mappers.MapOption[CloudTrailLog]) (CloudTrailLog, error) {
	line, ok := data.([]byte)
	if !ok {
		return CloudTrailLog{}, fmt.Errorf("expected []byte, got %T", data)
	}

	var record struct {
		Records []CloudTrailLog `json:"Records"`
	}

	if err := json.Unmarshal(line, &record); err != nil {
		// Try unmarshaling as a single record (for CloudWatch logs)
		var singleRecord CloudTrailLog
		if err := json.Unmarshal(line, &singleRecord); err != nil {
			return CloudTrailLog{}, fmt.Errorf("failed to parse CloudTrail log: %v", err)
		}
		return singleRecord, nil
	}

	if len(record.Records) == 0 {
		return CloudTrailLog{}, fmt.Errorf("no CloudTrail records found in log")
	}

	// Return the first record (we process one record at a time)
	return record.Records[0], nil
}
