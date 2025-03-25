package cloudwatch

import (
	"encoding/json"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
)

// CloudWatchCollectionState tracks collection state for multiple log streams
type CloudWatchCollectionState struct {
	// Map of log stream name to its time range collection state
	LogStreams map[string]*collection_state.TimeRangeCollectionStateImpl `json:"log_streams"`
	config     *AwsCloudWatchSourceConfig
}

// Init implements collection_state.CollectionState
func (s *CloudWatchCollectionState) Init(config *AwsCloudWatchSourceConfig, id string) error {
	s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
	s.config = config
	return nil
}

// IsEmpty implements collection_state.CollectionState
func (s *CloudWatchCollectionState) IsEmpty() bool {
	return len(s.LogStreams) == 0
}

// NewCloudWatchCollectionState creates a new CloudWatchCollectionState
func NewCloudWatchCollectionState() collection_state.CollectionState[*AwsCloudWatchSourceConfig] {
	return &CloudWatchCollectionState{
		LogStreams: make(map[string]*collection_state.TimeRangeCollectionStateImpl),
	}
}

// OnCollected updates the collection state for a specific log stream
func (s *CloudWatchCollectionState) OnCollected(logStreamName string, timestamp time.Time) error {
	// Get or create time range state for this log stream
	timeRangeState, exists := s.LogStreams[logStreamName]
	if !exists {
		timeRangeState = collection_state.NewTimeRangeCollectionStateImpl(collection_state.CollectionOrderChronological)
		s.LogStreams[logStreamName] = timeRangeState
	}

	return timeRangeState.OnCollected(logStreamName, timestamp)
}

// GetFromTime returns the earliest from time across all log streams
func (s *CloudWatchCollectionState) GetFromTime() time.Time {
	var earliestTime time.Time
	isFirst := true

	for _, state := range s.LogStreams {
		if isFirst {
			earliestTime = state.GetStartTime()
			isFirst = false
			continue
		}

		if state.GetStartTime().Before(earliestTime) {
			earliestTime = state.GetStartTime()
		}
	}

	return earliestTime
}

// GetToTime returns the latest to time across all log streams
func (s *CloudWatchCollectionState) GetToTime() time.Time {
	var latestTime time.Time
	isFirst := true

	for _, state := range s.LogStreams {
		if isFirst {
			latestTime = state.GetEndTime()
			isFirst = false
			continue
		}

		if state.GetEndTime().After(latestTime) {
			latestTime = state.GetEndTime()
		}
	}

	return latestTime
}

// GetEndTime implements collection_state.CollectionState
func (s *CloudWatchCollectionState) GetEndTime() time.Time {
	return s.GetToTime()
}

// GetGranularity implements collection_state.CollectionState
func (s *CloudWatchCollectionState) GetGranularity() time.Duration {
	return time.Minute
}

// SetGranularity implements collection_state.CollectionState
func (s *CloudWatchCollectionState) SetGranularity(granularity time.Duration) {
	// No-op since we use a fixed granularity
}

// SetEndTime implements collection_state.CollectionState
func (s *CloudWatchCollectionState) SetEndTime(endTime time.Time) {
	// No-op since we track end time per stream
}

// ShouldCollect implements collection_state.CollectionState
func (s *CloudWatchCollectionState) ShouldCollect(id string, timestamp time.Time) bool {
	return true
}

// Save implements collection_state.CollectionState
func (s *CloudWatchCollectionState) Save() error {
	_, err := json.Marshal(s)
	return err
}

// MarshalJSON implements json.Marshaler
func (s *CloudWatchCollectionState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.LogStreams)
}

// UnmarshalJSON implements json.Unmarshaler
func (s *CloudWatchCollectionState) UnmarshalJSON(data []byte) error {
	if s.LogStreams == nil {
		s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
	}
	return json.Unmarshal(data, &s.LogStreams)
}

// GetFromTimeForStream returns the from time for a specific log stream
func (s *CloudWatchCollectionState) GetFromTimeForStream(logStreamName string) time.Time {
	state, exists := s.LogStreams[logStreamName]
	if !exists {
		return time.Time{}
	}
	return state.GetStartTime()
}

// GetToTimeForStream returns the to time for a specific log stream
func (s *CloudWatchCollectionState) GetToTimeForStream(logStreamName string) time.Time {
	state, exists := s.LogStreams[logStreamName]
	if !exists {
		return time.Time{}
	}
	return state.GetEndTime()
}

// Clear implements collection_state.CollectionState
func (s *CloudWatchCollectionState) Clear() {
	s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
}

// GetConfig implements collection_state.CollectionState
func (s *CloudWatchCollectionState) GetConfig() *AwsCloudWatchSourceConfig {
	return s.config
}

// SetConfig implements collection_state.CollectionState
func (s *CloudWatchCollectionState) SetConfig(config *AwsCloudWatchSourceConfig) {
	s.config = config
}

// GetStartTime implements collection_state.CollectionState
func (s *CloudWatchCollectionState) GetStartTime() time.Time {
	return s.GetFromTime()
}
