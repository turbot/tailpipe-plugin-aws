package cloudwatch

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
)

// CloudWatchCollectionState tracks collection state for multiple log streams
type CloudWatchCollectionState struct {
	// Map of log stream name to its time range collection state
	LogStreams map[string]*collection_state.TimeRangeCollectionStateImpl `json:"log_streams"`
	config     *AwsCloudWatchSourceConfig
	// Path to the serialized collection state JSON
	jsonPath string
	// Time of last modification
	LastModifiedTime time.Time `json:"last_modified_time,omitempty"`
	// Time of last save
	lastSaveTime time.Time
}

// NewCloudWatchCollectionState creates a new CloudWatchCollectionState
func NewCloudWatchCollectionState() collection_state.CollectionState[*AwsCloudWatchSourceConfig] {
	return &CloudWatchCollectionState{
		LogStreams:       make(map[string]*collection_state.TimeRangeCollectionStateImpl),
		LastModifiedTime: time.Now(),
	}
}

// Init implements collection_state.CollectionState
func (s *CloudWatchCollectionState) Init(config *AwsCloudWatchSourceConfig, path string) error {
	s.jsonPath = path
	s.config = config

	// If there is a file at the path, load it
	if _, err := os.Stat(path); err == nil {
		jsonBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read collection state file: %w", err)
		}

		if err := json.Unmarshal(jsonBytes, s); err != nil {
			return fmt.Errorf("failed to unmarshal collection state: %w", err)
		}
	}

	// Initialize or reinitialize the map
	if s.LogStreams == nil {
		s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
	}

	return nil
}

// IsEmpty implements collection_state.CollectionState
func (s *CloudWatchCollectionState) IsEmpty() bool {
	return len(s.LogStreams) == 0
}

// Save implements collection_state.CollectionState
func (s *CloudWatchCollectionState) Save() error {
	// If the last save time is after the last modified time, then we have nothing to do
	if s.lastSaveTime.After(s.LastModifiedTime) {
		return nil
	}

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal collection state: %w", err)
	}

	// Ensure the target file path is valid
	if s.jsonPath == "" {
		return fmt.Errorf("collection state path is not set")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(s.jsonPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the JSON data to the file, overwriting any existing data
	if err := os.WriteFile(s.jsonPath, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to write collection state to file: %w", err)
	}

	// Update the last save time
	s.lastSaveTime = time.Now()
	slog.Debug("Successfully saved collection state", "path", s.jsonPath, "size", len(jsonBytes))

	return nil
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
	return time.Time{} // We don't track start time in this implementation
}

// GetEndTime implements collection_state.CollectionState
func (s *CloudWatchCollectionState) GetEndTime() time.Time {
	return time.Time{} // We don't track end time in this implementation
}

// SetEndTime implements collection_state.CollectionState
func (s *CloudWatchCollectionState) SetEndTime(endTime time.Time) {
	// No-op in this implementation
}

// GetGranularity implements collection_state.CollectionState
func (s *CloudWatchCollectionState) GetGranularity() time.Duration {
	return time.Minute
}

// SetGranularity implements collection_state.CollectionState
func (s *CloudWatchCollectionState) SetGranularity(granularity time.Duration) {
	// No-op since we use a fixed granularity
}

// Clear implements collection_state.CollectionState
func (s *CloudWatchCollectionState) Clear() {
	s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
}

// OnCollected updates the collection state for a specific log stream
func (s *CloudWatchCollectionState) OnCollected(logStreamName string, timestamp time.Time) error {
	// Update last modified time
	s.LastModifiedTime = time.Now()

	// Get or create time range state for this log stream
	timeRangeState, exists := s.LogStreams[logStreamName]
	if !exists {
		timeRangeState = collection_state.NewTimeRangeCollectionStateImpl(collection_state.CollectionOrderChronological)
		timeRangeState.SetGranularity(s.GetGranularity())
		s.LogStreams[logStreamName] = timeRangeState
	}

	// Call OnCollected on the time range state
	if err := timeRangeState.OnCollected(logStreamName, timestamp); err != nil {
		return fmt.Errorf("failed to update time range state for stream %s: %w", logStreamName, err)
	}

	return nil
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

// ShouldCollect implements collection_state.CollectionState
func (s *CloudWatchCollectionState) ShouldCollect(id string, timestamp time.Time) bool {
	state, exists := s.LogStreams[id]
	if !exists {
		return true
	}
	return state.ShouldCollect(id, timestamp)
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
