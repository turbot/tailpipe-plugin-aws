package sources

import (
	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
)

// AwsCloudwatchCollectionState contains collection state data for the AwsCloudwatchSource artifact source
// This contains the latest timestamp fetched for each log stream in a SINGLE log group
type AwsCloudwatchCollectionState struct {
	collection_state.CollectionStateImpl[*AwsCloudWatchSourceConfig]
	// The timestamp of the last collected log for each log stream
	// expressed as the number of milliseconds after Jan 1, 1970 00:00:00 UTC.
	Timestamps map[string]int64 `json:"timestamps"`
}

func NewAwsCloudwatchCollectionState() collection_state.CollectionState[*AwsCloudWatchSourceConfig] {
	// TODO handle storing path/loading/saving state
	return &AwsCloudwatchCollectionState{
		Timestamps: make(map[string]int64),
	}
}

func (s *AwsCloudwatchCollectionState) Init(config *AwsCloudWatchSourceConfig, path string) error {
	return s.CollectionStateImpl.Init(config, path)
}

// Upsert adds new/updates an existing log stream  with its current timestamp
func (s *AwsCloudwatchCollectionState) Upsert(name string, time int64) {
	s.Mut.Lock()
	defer s.Mut.Unlock()

	if s.Timestamps == nil {
		s.Timestamps = make(map[string]int64)
	}
	if time == 0 {
		return
	}

	currentTime := s.Timestamps[name]
	if time > currentTime {
		s.Timestamps[name] = time
	}
}

func (s *AwsCloudwatchCollectionState) IsEmpty() bool {
	return len(s.Timestamps) == 0
}
