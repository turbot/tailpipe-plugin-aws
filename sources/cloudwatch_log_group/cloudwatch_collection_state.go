// Package cloudwatch provides functionality to collect logs from AWS CloudWatch Log Groups
package cloudwatch_log_group

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
)

// CloudWatchCollectionState tracks collection state for multiple log streams within a CloudWatch log group.
// It maintains a map of log stream names to their individual time range collection states,
// allowing for incremental collection and resumption of collection from the last processed event.
type CloudWatchCollectionState struct {
	// Map of log stream name to its time range collection state
	LogStreams map[string]*collection_state.TimeRangeCollectionStateImpl `json:"log_streams"`
	// Configuration for the CloudWatch source
	config *AwsCloudWatchSourceConfig
	// Path to the serialized collection state JSON file
	jsonPath string
	// Time when the collection state was last modified
	LastModifiedTime time.Time `json:"last_modified_time,omitempty"`
	// Time when the collection state was last saved to disk
	lastSaveTime time.Time
}

// NewCloudWatchCollectionState creates a new CloudWatchCollectionState instance.
// It initializes an empty map for log streams and sets the initial modification time.
func NewCloudWatchCollectionState() collection_state.CollectionState[*AwsCloudWatchSourceConfig] {
	return &CloudWatchCollectionState{
		LogStreams:       make(map[string]*collection_state.TimeRangeCollectionStateImpl),
		LastModifiedTime: time.Now(),
	}
}

// Init initializes the collection state with the provided configuration and state file path.
// If a state file exists at the given path, it loads and deserializes the state.
// If no file exists or the state is empty, it initializes a new empty state.
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

	// Initialize or reinitialize the map if nil
	if s.LogStreams == nil {
		s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
	}

	return nil
}

// IsEmpty returns true if no log streams have been collected yet
func (s *CloudWatchCollectionState) IsEmpty() bool {
	return len(s.LogStreams) == 0
}

// Save persists the current collection state to disk if it has been modified since the last save.
// The state is serialized as JSON and written to the configured file path.
// It creates any necessary directories and updates the last save time on success.
func (s *CloudWatchCollectionState) Save() error {
	// Skip save if no modifications since last save
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

// GetConfig returns the current CloudWatch source configuration
func (s *CloudWatchCollectionState) GetConfig() *AwsCloudWatchSourceConfig {
	return s.config
}

// SetConfig updates the CloudWatch source configuration
func (s *CloudWatchCollectionState) SetConfig(config *AwsCloudWatchSourceConfig) {
	s.config = config
}

// GetStartTime returns an empty time.Time since global start time is not tracked
// Individual log stream start times are tracked separately
func (s *CloudWatchCollectionState) GetStartTime() time.Time {
	return time.Time{} // We don't track start time in this implementation
}

// GetEndTime returns an empty time.Time since global end time is not tracked
// Individual log stream end times are tracked separately
func (s *CloudWatchCollectionState) GetEndTime() time.Time {
	return time.Time{} // We don't track end time in this implementation
}

// SetEndTime is a no-op since global end time is not tracked
func (s *CloudWatchCollectionState) SetEndTime(endTime time.Time) {
	// No-op in this implementation
}

// GetGranularity returns the time granularity for collection state tracking
// Uses a fixed granularity of one minute
func (s *CloudWatchCollectionState) GetGranularity() time.Duration {
	return time.Minute
}

// SetGranularity is a no-op since we use a fixed granularity
func (s *CloudWatchCollectionState) SetGranularity(granularity time.Duration) {
	// No-op since we use a fixed granularity
}

// Clear resets the collection state by removing all log stream states
func (s *CloudWatchCollectionState) Clear() {
	s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
}

// OnCollected updates the collection state for a specific log stream when an event is processed.
// It creates a new time range state for the stream if it doesn't exist,
// and updates the last modified time to trigger a state save.
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

// GetFromTime returns the earliest start time across all log streams.
// This represents the earliest point in time from which we have collected events.
func (s *CloudWatchCollectionState) GetFromTime() time.Time {
	var earliestTime time.Time

	for _, state := range s.LogStreams {
		startTime := state.GetStartTime()

		if earliestTime.IsZero() || startTime.Before(earliestTime) {
			earliestTime = startTime
		}
	}

	return earliestTime
}

// GetToTime returns the latest end time across all log streams.
// This represents the most recent point in time up to which we have collected events.
func (s *CloudWatchCollectionState) GetToTime() time.Time {
	var latestTime time.Time

	for _, state := range s.LogStreams {
		endTime := state.GetEndTime()

		if latestTime.IsZero() || endTime.After(latestTime) {
			latestTime = endTime
		}
	}

	return latestTime
}

// ShouldCollect determines whether an event with the given timestamp should be collected
// for the specified log stream based on its time range state.
func (s *CloudWatchCollectionState) ShouldCollect(id string, timestamp time.Time) bool {
	if state, exists := s.LogStreams[id]; exists {
		return state.ShouldCollect(id, timestamp)
	}
	return true
}

// GetStartTimeForStream returns the start time for a specific log stream.
// If the stream doesn't exist in the state, returns zero time.
func (s *CloudWatchCollectionState) GetStartTimeForStream(logStreamName string) time.Time {
	if state, exists := s.LogStreams[logStreamName]; exists {
		return state.GetStartTime()
	}
	return time.Time{}
}

// GetEndTimeForStream returns the end time for a specific log stream.
// If the stream doesn't exist in the state, returns zero time.
func (s *CloudWatchCollectionState) GetEndTimeForStream(logStreamName string) time.Time {
	if state, exists := s.LogStreams[logStreamName]; exists {
		return state.GetEndTime()
	}
	return time.Time{}
}
