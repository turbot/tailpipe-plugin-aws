package sources

import (
	"encoding/json"
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
	"log/slog"
	"os"
	"sync"
	"time"
)

// AwsCloudwatchCollectionState contains collection state data for the AwsCloudwatchSource artifact source
// This contains the latest timestamp fetched for each log stream in a SINGLE log group
type AwsCloudwatchCollectionState struct {
	// The timestamp of the last collected log for each log stream
	// expressed as the number of milliseconds after Jan 1, 1970 00:00:00 UTC.
	LogStreamTimestamps map[string]time.Time `json:"timestamps"`

	// path to the serialised collection state JSON
	jsonPath         string
	lastModifiedTime time.Time
	lastSaveTime     time.Time

	mut *sync.RWMutex
}

func NewAwsCloudwatchCollectionState() collection_state.CollectionState[*AwsCloudWatchSourceConfig] {
	return &AwsCloudwatchCollectionState{
		LogStreamTimestamps: make(map[string]time.Time),
		mut:                 &sync.RWMutex{},
	}
}

func (s *AwsCloudwatchCollectionState) Init(_ *AwsCloudWatchSourceConfig, path string) error {
	s.jsonPath = path

	// if there is a file at the path, load it
	if _, err := os.Stat(path); err == nil {
		// TODO #err should we just warn and delete/rename the file
		// read the file
		jsonBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read collection state file '%s': %w", path, err)
		}
		err = json.Unmarshal(jsonBytes, s)
		if err != nil {
			return fmt.Errorf("failed to unmarshal collection state file '%s': %w", path, err)
		}
	}
	return nil
}

func (s *AwsCloudwatchCollectionState) Save() error {
	s.mut.Lock()
	defer s.mut.Unlock()

	// if the last save time is after the last modified time, then we have nothing to do
	if s.lastSaveTime.After(s.lastModifiedTime) {
		slog.Info("collection state has not been modified since last save")
		// nothing to do
		return nil
	}
	slog.Info("Saving collection state", "lastModifiedTime", s.lastModifiedTime, "lastSaveTime", s.lastSaveTime)

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	// ensure the target file path is valid
	if s.jsonPath == "" {
		return fmt.Errorf("collection state path is not set")
	}

	// if we are empty, delete the file
	if s.IsEmpty() {
		err := os.Remove(s.jsonPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete collection state file: %w", err)
		}
		return nil
	}

	// write the JSON data to the file, overwriting any existing data
	err = os.WriteFile(s.jsonPath, jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write collection state to file: %w", err)
	}

	// update the last save time
	s.lastSaveTime = time.Now()

	return nil
}

func (s *AwsCloudwatchCollectionState) ShouldCollect(id string, timestamp time.Time) bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	// get timestamp for this log stream
	logStreamTime := s.LogStreamTimestamps[id]

	// if this timestamp is after the latest timestamp we have stored, we SHOULD collect
	return timestamp.After(logStreamTime)
}

func (s *AwsCloudwatchCollectionState) OnCollected(id string, timestamp time.Time) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	// store modified time to ensure we save the state
	s.lastModifiedTime = time.Now()

	logStreamTime := s.LogStreamTimestamps[id]
	if timestamp.After(logStreamTime) {
		s.LogStreamTimestamps[id] = timestamp
	}
	return nil
}

// GetStartTime returns the latest start time of all the log streams
func (s *AwsCloudwatchCollectionState) GetStartTime() time.Time {
	// find the earliest end time of all the log streams
	var startTime time.Time
	for _, timestamp := range s.LogStreamTimestamps {
		if startTime.IsZero() || timestamp.After(startTime) {
			startTime = timestamp
		}
	}

	return startTime
}

// GetEndTime returns the earliest end time of all the log streams
// (not currently used as the Cloudwatch source retrieves the end time for the specific log stream	)
func (s *AwsCloudwatchCollectionState) GetEndTime() time.Time {
	// find the earliest end time of all the log streams
	var endTime time.Time
	for _, timestamp := range s.LogStreamTimestamps {
		if endTime.IsZero() || timestamp.Before(endTime) {
			endTime = timestamp
		}
	}

	return endTime
}

func (s *AwsCloudwatchCollectionState) SetEndTime(newEndTime time.Time) {
	for id := range s.LogStreamTimestamps {
		if s.LogStreamTimestamps[id].After(newEndTime) {
			s.LogStreamTimestamps[id] = newEndTime
		}
	}
}

func (s *AwsCloudwatchCollectionState) Clear() {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.LogStreamTimestamps = make(map[string]time.Time)
}

func (s *AwsCloudwatchCollectionState) IsEmpty() bool {
	return len(s.LogStreamTimestamps) == 0
}

func (s *AwsCloudwatchCollectionState) GetGranularity() time.Duration {
	// this will never be called - we do not use granularity for this source
	// however - return the 'correct' value to avoid confusion
	return time.Nanosecond
}

func (s *AwsCloudwatchCollectionState) SetGranularity(time.Duration) {
	// do nothing - this should not be called
}
