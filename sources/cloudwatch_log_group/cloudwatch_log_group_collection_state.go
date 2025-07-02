package cloudwatch_log_group

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
)

type CloudWatchLogGroupCollectionStateLegacy struct {
	LogStreams       map[string]*collection_state.TimeRangeCollectionStateLegacy `json:"log_streams,omitempty"`
	LastModifiedTime time.Time                                                   `json:"last_modified_time,omitempty"`
}

// CloudWatchLogGroupCollectionState tracks collection state for multiple log streams within a CloudWatch log group.
// It maintains a map of log stream names to their individual time range collection states,
// allowing for incremental collection and resumption of collection from the last processed event.
type CloudWatchLogGroupCollectionState struct {
	// Map of log stream name to its time range collection state
	LogStreams map[string]*collection_state.TimeRangeCollectionState `json:"log_streams"`
	// Configuration for the CloudWatch source
	//config *AwsCloudWatchLogGroupSourceConfig
	// the time range for the underway collection - populated by Init
	currentDirectionalTimeRange *collection_state.DirectionalTimeRange
	// Granularity defines the time resolution for collection state updates
	Granularity time.Duration `json:"granularity,omitempty"`
}

// NewCloudWatchLogGroupCollectionState creates a new CloudWatchCollectionState instance.
// It initializes an empty map for log streams and sets the initial modification time.
func NewCloudWatchLogGroupCollectionState() collection_state.CollectionState {
	return &CloudWatchLogGroupCollectionState{
		LogStreams: make(map[string]*collection_state.TimeRangeCollectionState),
	}
}

// Init initializes the collection state with the provided configuration and state file path.
// If a state file exists at the given path, it loads and deserializes the state.
// If no file exists or the state is empty, it initializes a new empty state.
func (s *CloudWatchLogGroupCollectionState) Init(timeRange collection_state.DirectionalTimeRange, granularity time.Duration) {
	s.Granularity = granularity

	// Initialize or reinitialize the maps if nil
	if s.LogStreams == nil {
		s.LogStreams = make(map[string]*collection_state.TimeRangeCollectionState)
	}
	// init all log streams with the provided time range and granularity
	for _, state := range s.LogStreams {
		if state != nil {
			state.Init(timeRange, granularity)
		}
	}

	s.currentDirectionalTimeRange = &timeRange

	return
}

// IsEmpty returns true if no log streams have been collected yet
func (s *CloudWatchLogGroupCollectionState) IsEmpty() bool {
	return len(s.LogStreams) == 0
}

// SetConfig updates the CloudWatch source configuration
func (s *CloudWatchLogGroupCollectionState) SetConfig(config *AwsCloudWatchLogGroupSourceConfig) {
	//s.config = config
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

// OnCollected updates the collection state for a specific log stream when an event is processed.
// It creates a new time range state for the stream if it doesn't exist,
// and updates the last modified time to trigger a state save.
func (s *CloudWatchLogGroupCollectionState) OnCollected(logStreamName string, timestamp time.Time) error {
	if s.currentDirectionalTimeRange == nil {
		return fmt.Errorf("currentDirectionalTimeRange is nil - Init must be called before OnCollected")
	}

	// Get or create time range state for this log stream
	timeRangeState, exists := s.LogStreams[logStreamName]
	if !exists {
		timeRangeState = collection_state.NewTimeRangeCollectionState().(*collection_state.TimeRangeCollectionState)
		timeRangeState.Init(*s.currentDirectionalTimeRange, s.Granularity)
		timeRangeState.Order = collection_state.CollectionOrderChronological

		s.LogStreams[logStreamName] = timeRangeState
	}

	// Call OnCollected on the time range state
	if err := timeRangeState.OnCollected(logStreamName, timestamp); err != nil {
		return fmt.Errorf("failed to update time range state for stream %s: %w", logStreamName, err)
	}

	return nil
}

func (s *CloudWatchLogGroupCollectionState) OnCollectionComplete() error {
	for _, logStreamState := range s.LogStreams {
		if logStreamState == nil {
			continue
		}
		// set the end time of the trunk state to the end time of the current collection
		err := logStreamState.OnCollectionComplete()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFromTime returns the earliest start time across all log streams.
// This represents the earliest point in time from which we have collected events.
func (s *CloudWatchLogGroupCollectionState) GetFromTime() time.Time {
	var earliestTime time.Time

	for _, state := range s.LogStreams {
		startTime := state.GetFromTime()

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
		endTime := state.GetToTime()

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
		return state.ShouldCollect(logStreamName, timestamp)
	}
	return true
}

// GetStartTimeForStream returns the start time for a specific log stream.
// If the stream doesn't exist in the state, returns zero time.
func (s *CloudWatchLogGroupCollectionState) GetStartTimeForStream(logStreamName string) time.Time {
	if state, exists := s.LogStreams[logStreamName]; exists {
		return state.GetFromTime()
	}
	return time.Time{}
}

// GetEndTimeForStream returns the end time for a specific log stream.
// If the stream doesn't exist in the state, returns zero time.
func (s *CloudWatchLogGroupCollectionState) GetEndTimeForStream(logStreamName string) time.Time {
	if state, exists := s.LogStreams[logStreamName]; exists {
		return state.GetToTime()
	}
	return time.Time{}
}

func (s *CloudWatchLogGroupCollectionState) Clear(timeRange collection_state.DirectionalTimeRange) {
	for _, state := range s.LogStreams {
		if state != nil {
			state.Clear(timeRange)
		}
	}

}

func (s *CloudWatchLogGroupCollectionState) MigrateFromLegacyState(bytes []byte) error {
	legacyState := &CloudWatchLogGroupCollectionStateLegacy{}
	err := json.Unmarshal(bytes, legacyState)
	if err != nil {
		return fmt.Errorf("failed to unmarshal legacy collection state: %w", err)
	}

	// Convert each trunk state from legacy format to new format
	for trunkPath, legacyTrunkState := range legacyState.LogStreams {
		if legacyTrunkState == nil {
			// Skip nil trunk states
			continue
		}

		// Use the new constructor for legacy trunk states
		s.LogStreams[trunkPath] = collection_state.NewTimeRangeCollectionStateFromLegacy(legacyTrunkState)
	}

	return nil
}

func (s *CloudWatchLogGroupCollectionState) Validate() error {
	var errorList []error
	for _, trunkState := range s.LogStreams {
		if trunkErr := trunkState.Validate(); trunkErr != nil {
			errorList = append(errorList, trunkErr)
		}
	}
	if len(errorList) > 0 {
		return fmt.Errorf("validation failed for artifact collection state: %w", errors.Join(errorList...))
	}
	return nil
}
