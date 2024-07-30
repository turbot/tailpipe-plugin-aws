package aws_source

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/artifact"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
)

// CloudwatchMapper is an [plugin.Mapper] implementation
// that receives JSON data and returns AWSCloudTrail records
type CloudwatchMapper struct {
}

// NewCloudwatchMapper creates a new CloudwatchMapper
func NewCloudwatchMapper() artifact.Mapper {
	return &CloudwatchMapper{}
}

func (c *CloudwatchMapper) Identifier() string {
	return artifact.AWSCloudwatchMapperIdentifier
}

// Map unmarshalls JSON into an AWSCloudTrailBatch object and extracts AWSCloudTrail records from it
func (c *CloudwatchMapper) Map(_ context.Context, a *artifact.ArtifactData) ([]*artifact.ArtifactData, error) {
	// the expected input type is a JSON string deserializable to a map with keys "IngestionTime", "Timestamp" and "Message"
	jsonString, ok := a.Data.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %T", a)
	}

	var cloudwatchEntry map[string]any
	err := json.Unmarshal([]byte(jsonString), &cloudwatchEntry)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}
	row := cloudwatchEntry["Message"].(string)
	metadata := a.Metadata.Clone()
	metadata.TpIngestTimestamp = helpers.UnixMillis(cloudwatchEntry["IngestionTime"].(float64))
	metadata.TpTimestamp = helpers.UnixMillis(cloudwatchEntry["Timestamp"].(float64))

	d := artifact.NewData(row, artifact.WithMetadata(metadata))

	return []*artifact.ArtifactData{d}, nil
}
