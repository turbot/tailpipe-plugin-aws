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

// CloudWatchLogGroupCollectionState tracks collection state for multiple log streams within a CloudWatch log group.
// It maintains a map of log stream names to their individual time range collection states,
// allowing for incremental collection and resumption of collection from the last processed event.
type CloudWatchLogGroupCollectionState struct {
	// Map of log stream name to its time range collection state
	LogStreams map[string]*collection_state.TimeRangeCollectionStateImpl `json:"log_streams"`
	// Configuration for the CloudWatch source
	config *AwsCloudWatchLogGroupSourceConfig
	// Path to the serialized collection state JSON file
	jsonPath string
	// Time when the collection state was last modified
	LastModifiedTime time.Time `json:"last_modified_time,omitempty"`
	// Time when the collection state was last saved to disk
	lastSaveTime time.Time
}

// NewCloudWatchLogGroupCollectionState creates a new CloudWatchCollectionState instance.
// It initializes an empty map for log streams and sets the initial modification time.
func NewCloudWatchLogGroupCollectionState() collection_state.CollectionState[*AwsCloudWatchLogGroupSourceConfig] {
	return &CloudWatchLogGroupCollectionState{
		LogStreams:       make(map[string]*collection_state.TimeRangeCollectionStateImpl),
		LastModifiedTime: time.Now(),
	}
}

// Init initializes the collection state with the provided configuration and state file path.
// If a state file exists at the given path, it loads and deserializes the state.
// If no file exists or the state is empty, it initializes a new empty state.
func (s *CloudWatchLogGroupCollectionState) Init(config *AwsCloudWatchLogGroupSourceConfig, path string) error {
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

	// Initialize or reinitialize the maps if nil
	if s.LogStreams == nil {
		s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
	}

	return nil
}

// IsEmpty returns true if no log streams have been collected yet
func (s *CloudWatchLogGroupCollectionState) IsEmpty() bool {
	return len(s.LogStreams) == 0
}

// Save persists the current collection state to disk if it has been modified since the last save.
// The state is serialized as JSON and written to the configured file path.
// It creates any necessary directories and updates the last save time on success.
func (s *CloudWatchLogGroupCollectionState) Save() error {
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
func (s *CloudWatchLogGroupCollectionState) GetConfig() *AwsCloudWatchLogGroupSourceConfig {
	return s.config
}

// SetConfig updates the CloudWatch source configuration
func (s *CloudWatchLogGroupCollectionState) SetConfig(config *AwsCloudWatchLogGroupSourceConfig) {
	s.config = config
}

// GetStartTime returns an empty time.Time since global start time is not tracked
// Individual log stream start times are tracked separately
func (s *CloudWatchLogGroupCollectionState) GetStartTime() time.Time {
	return time.Time{} // We don't track start time in this implementation
}

// GetEndTime returns an empty time.Time since global end time is not tracked
// Individual log stream end times are tracked separately
func (s *CloudWatchLogGroupCollectionState) GetEndTime() time.Time {
	return time.Time{} // We don't track end time in this implementation
}

// SetEndTime is a no-op since global end time is not tracked
func (s *CloudWatchLogGroupCollectionState) SetEndTime(endTime time.Time) {
	// No-op in this implementation
}

// GetGranularity returns the time granularity for collection state tracking
// Uses a fixed granularity of one minute
func (s *CloudWatchLogGroupCollectionState) GetGranularity() time.Duration {
	return time.Minute
}

// SetGranularity is a no-op since we use a fixed granularity
func (s *CloudWatchLogGroupCollectionState) SetGranularity(granularity time.Duration) {
	// No-op since we use a fixed granularity
}

// Clear resets the collection state by removing all log stream states
func (s *CloudWatchLogGroupCollectionState) Clear() {
	s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionStateImpl)
}

// OnCollected updates the collection state for a specific log stream when an event is processed.
// It creates a new time range state for the stream if it doesn't exist,
// and updates the last modified time to trigger a state save.
func (s *CloudWatchLogGroupCollectionState) OnCollected(logStreamName string, timestamp time.Time) error {
	// Update last modified time
	s.LastModifiedTime = time.Now()

	// Get or create time range state for this log stream
	timeRangeState, exists := s.LogStreams[logStreamName]
	if !exists {
		timeRangeState = collection_state.NewTimeRangeCollectionStateImpl(collection_state.CollectionOrderChronological)
		timeRangeState.SetGranularity(s.GetGranularity())
		s.LogStreams[logStreamName] = timeRangeState
		// s.ProcessedEventIds = append(s.ProcessedEventIds, eventId)
	}

	// Call OnCollected on the time range state
	if err := timeRangeState.OnCollected(logStreamName, timestamp); err != nil {
		return fmt.Errorf("failed to update time range state for stream %s: %w", logStreamName, err)
	}

	return nil
}

// GetFromTime returns the earliest start time across all log streams.
// This represents the earliest point in time from which we have collected events.
func (s *CloudWatchLogGroupCollectionState) GetFromTime() time.Time {
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
func (s *CloudWatchLogGroupCollectionState) GetToTime() time.Time {
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
func (s *CloudWatchLogGroupCollectionState) ShouldCollect(logStreamName string, timestamp time.Time) bool {
	if state, exists := s.LogStreams[logStreamName]; exists {
		//  Granularity defines a small timing window (buffer zone) added after the end time of the last collected event.
		// It helps decide whether a new event should be considered already collected or needs to be collected.
		// It compensates for small delays or natural jitter between closely occurring events.
		//
		// Based on analysis of CloudWatch log streams:
		// - Smallest gap between events is ~4-35ms
		// - Use 6ms for stricter separation

		// When running the command `tailpipe collect <table_name>.<partition>` multiple times with a larger granularity,
		// duplicate events may be collected if those events were already collected in a previous run.
		state.Granularity = 6 * time.Millisecond
		return state.ShouldCollect(logStreamName, timestamp)
	}
	return true
}

// GetStartTimeForStream returns the start time for a specific log stream.
// If the stream doesn't exist in the state, returns zero time.
func (s *CloudWatchLogGroupCollectionState) GetStartTimeForStream(logStreamName string) time.Time {
	if state, exists := s.LogStreams[logStreamName]; exists {
		return state.GetStartTime()
	}
	return time.Time{}
}

// GetEndTimeForStream returns the end time for a specific log stream.
// If the stream doesn't exist in the state, returns zero time.
func (s *CloudWatchLogGroupCollectionState) GetEndTimeForStream(logStreamName string) time.Time {
	if state, exists := s.LogStreams[logStreamName]; exists {
		return state.GetEndTime()
	}
	return time.Time{}
}
