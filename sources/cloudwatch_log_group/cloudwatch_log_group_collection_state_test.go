package cloudwatch_log_group

import (
	"encoding/json"
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
	"testing"
	"time"
)

func TestCloudWatchLogGroupCollectionState_MigrateFromLegacyState(t *testing.T) {
	tests := []struct {
		name     string
		legacy   *CloudWatchLogGroupCollectionStateLegacy
		expected *CloudWatchLogGroupCollectionState
	}{
		{
			name: "migrate two trunks with different orders",
			legacy: buildCloudWatchLogGroupCollectionStateLegacy(map[string]*collection_state.TimeRangeCollectionStateLegacy{
				"/trunk1": buildTimeRangeCollectionStateLegacy("2023-10-01 00:00:00", "2023-12-01 01:00:00", time.Hour*24, collection_state.CollectionOrderChronological, "object1", "object2"),
				"/trunk2": buildTimeRangeCollectionStateLegacy("2023-11-01 00:00:00", "2023-11-30 00:00:00", time.Hour*24, collection_state.CollectionOrderReverse, "object3"),
			}, timeString("2023-12-01 12:00:00")),
			expected: buildCloudWatchLogGroupCollectionState(map[string]*collection_state.TimeRangeCollectionState{
				"/trunk1": buildTimeRangeCollectionState(collection_state.CollectionOrderChronological, time.Hour*24,
					buildTimeRangeState("2023-10-01 00:00:00", "2023-11-30 01:00:00", time.Hour*24, collection_state.CollectionOrderChronological, "object1", "object2"),
				),
				"/trunk2": buildTimeRangeCollectionState(collection_state.CollectionOrderReverse, time.Hour*24,
					buildTimeRangeState("2023-11-30 00:00:00", "2023-11-01 00:00:00", time.Hour*24, collection_state.CollectionOrderReverse, "object3"),
				),
			}, time.Hour*24),
		},
		{
			name:     "empty trunks",
			legacy:   buildCloudWatchLogGroupCollectionStateLegacy(map[string]*collection_state.TimeRangeCollectionStateLegacy{}, timeString("2023-12-01 12:00:00")),
			expected: buildCloudWatchLogGroupCollectionState(map[string]*collection_state.TimeRangeCollectionState{}, 0),
		},
		{
			name: "nil trunk",
			legacy: buildCloudWatchLogGroupCollectionStateLegacy(map[string]*collection_state.TimeRangeCollectionStateLegacy{
				"/trunk1": nil,
			}, timeString("2023-12-01 12:00:00")),
			expected: buildCloudWatchLogGroupCollectionState(map[string]*collection_state.TimeRangeCollectionState{}, 0),
		},
		{
			name: "trunk with no end objects",
			legacy: buildCloudWatchLogGroupCollectionStateLegacy(map[string]*collection_state.TimeRangeCollectionStateLegacy{
				"/trunk1": buildTimeRangeCollectionStateLegacy("2023-10-01 00:00:00", "2023-12-01 01:00:00", time.Hour*24, collection_state.CollectionOrderChronological),
			}, timeString("2023-12-01 12:00:00")),
			expected: buildCloudWatchLogGroupCollectionState(map[string]*collection_state.TimeRangeCollectionState{
				"/trunk1": buildTimeRangeCollectionState(collection_state.CollectionOrderChronological, time.Hour*24,
					buildTimeRangeState("2023-10-01 00:00:00", "2023-11-30 01:00:00", time.Hour*24, collection_state.CollectionOrderChronological),
				),
			}, time.Hour*24),
		},
		{
			name: "different granularity",
			legacy: buildCloudWatchLogGroupCollectionStateLegacy(map[string]*collection_state.TimeRangeCollectionStateLegacy{
				"/trunk1": buildTimeRangeCollectionStateLegacy("2023-10-01 00:00:00", "2023-12-01 01:00:00", time.Hour*24, collection_state.CollectionOrderChronological, "object1"),
				"/trunk2": buildTimeRangeCollectionStateLegacy("2023-11-01 00:00:00", "2023-11-30 00:00:00", time.Hour, collection_state.CollectionOrderReverse, "object2"),
			}, timeString("2023-12-01 12:00:00")),
			expected: buildCloudWatchLogGroupCollectionState(map[string]*collection_state.TimeRangeCollectionState{
				"/trunk1": buildTimeRangeCollectionState(collection_state.CollectionOrderChronological, time.Hour*24,
					buildTimeRangeState("2023-10-01 00:00:00", "2023-11-30 01:00:00", time.Hour*24, collection_state.CollectionOrderChronological, "object1"),
				),
				"/trunk2": buildTimeRangeCollectionState(collection_state.CollectionOrderReverse, time.Hour,
					buildTimeRangeState("2023-11-30 00:00:00", "2023-11-01 00:00:00", time.Hour, collection_state.CollectionOrderReverse, "object2"),
				),
			}, time.Hour*24),
		},
		{
			name: "reverse order with multiple end objects",
			legacy: buildCloudWatchLogGroupCollectionStateLegacy(map[string]*collection_state.TimeRangeCollectionStateLegacy{
				"/trunk1": buildTimeRangeCollectionStateLegacy("2023-10-01 00:00:00", "2023-12-01 01:00:00", time.Hour*24, collection_state.CollectionOrderReverse, "object1", "object2"),
			}, timeString("2023-12-01 12:00:00")),
			expected: buildCloudWatchLogGroupCollectionState(map[string]*collection_state.TimeRangeCollectionState{
				"/trunk1": buildTimeRangeCollectionState(collection_state.CollectionOrderReverse, time.Hour*24,
					buildTimeRangeState("2023-12-01 01:00:00", "2023-10-01 00:00:00", time.Hour*24, collection_state.CollectionOrderReverse, "object1", "object2"),
				),
			}, time.Hour*24),
		},
		{
			name: "single trunk, single object",
			legacy: buildCloudWatchLogGroupCollectionStateLegacy(map[string]*collection_state.TimeRangeCollectionStateLegacy{
				"/trunk1": buildTimeRangeCollectionStateLegacy("2023-10-01 00:00:00", "2023-10-02 00:00:00", time.Hour*24, collection_state.CollectionOrderChronological, "object1"),
			}, timeString("2023-12-01 12:00:00")),
			expected: buildCloudWatchLogGroupCollectionState(map[string]*collection_state.TimeRangeCollectionState{
				"/trunk1": buildTimeRangeCollectionState(collection_state.CollectionOrderChronological, time.Hour*24,
					buildTimeRangeState("2023-10-01 00:00:00", "2023-10-01 00:00:00", time.Hour*24, collection_state.CollectionOrderChronological, "object1"),
				),
			}, time.Hour*24),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			legacyBytes, err := json.Marshal(tt.legacy)
			if err != nil {
				t.Fatalf("Failed to marshal legacy state: %v", err)
			}

			newState := NewCloudWatchLogGroupCollectionState().(*CloudWatchLogGroupCollectionState)
			err = newState.MigrateFromLegacyState(legacyBytes)
			if err != nil {
				t.Fatalf("Failed to migrate legacy state: %v", err)
			}

			if equal, msg := cloudWatchLogGroupCollectionStateEquals(newState, tt.expected); !equal {
				t.Errorf("state after migration: %s", msg)
			}
		})
	}
}

// buildCloudWatchLogGroupCollectionStateLegacy constructs a legacy artifact collection state for tests
func buildCloudWatchLogGroupCollectionStateLegacy(trunks map[string]*collection_state.TimeRangeCollectionStateLegacy, lastModifiedTime time.Time) *CloudWatchLogGroupCollectionStateLegacy {
	return &CloudWatchLogGroupCollectionStateLegacy{
		LogStreams:       trunks,
		LastModifiedTime: lastModifiedTime,
	}
}

// buildTimeRangeCollectionStateLegacy constructs a legacy time range collection state for tests
func buildTimeRangeCollectionStateLegacy(fromStr, toStr string, granularity time.Duration, order collection_state.CollectionOrder, endObjects ...string) *collection_state.TimeRangeCollectionStateLegacy {
	from, err := time.Parse("2006-01-02 15:04:05", fromStr)
	if err != nil {
		panic(err)
	}
	to, err := time.Parse("2006-01-02 15:04:05", toStr)
	if err != nil {
		panic(err)
	}
	endObjectsMap := make(map[string]struct{})
	for _, obj := range endObjects {
		endObjectsMap[obj] = struct{}{}
	}
	return &collection_state.TimeRangeCollectionStateLegacy{
		FirstEntryTime:  from,
		LastEntryTime:   to,
		EndTime:         to.Add(-granularity),
		EndObjects:      endObjectsMap,
		Granularity:     granularity,
		CollectionOrder: order,
	}
}

func cloudWatchLogGroupCollectionStateEquals(got, want *CloudWatchLogGroupCollectionState) (bool, string) {
	if len(got.LogStreams) != len(want.LogStreams) {
		return false, fmt.Sprintf("trunk count = %v, want %v", len(got.LogStreams), len(want.LogStreams))
	}
	for k, expectedStream := range want.LogStreams {
		actualStream, ok := got.LogStreams[k]
		if !ok {
			return false, fmt.Sprintf("missing trunk %v", k)
		}
		if equal, msg := actualStream.Compare(expectedStream); !equal {
			return false, fmt.Sprintf("trunk %v: %s", k, msg)
		}
	}
	//if got.GetGranularity() != want.GetGranularity() {
	//	return false, fmt.Sprintf("granularity = %v, want %v", got.GetGranularity(), want.GetGranularity())
	//}
	return true, ""
}
func buildTimeRangeCollectionState(order collection_state.CollectionOrder, granularity time.Duration, ranges ...*collection_state.TimeRangeObjectState) *collection_state.TimeRangeCollectionState {
	return &collection_state.TimeRangeCollectionState{
		TimeRanges:  ranges,
		Granularity: granularity,
		Order:       order,
	}
}

func buildCloudWatchLogGroupCollectionState(logStreams map[string]*collection_state.TimeRangeCollectionState, granularity time.Duration) *CloudWatchLogGroupCollectionState {
	return &CloudWatchLogGroupCollectionState{
		LogStreams: logStreams,
	}
}

func timeString(timeStr string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		panic(err)
	}
	return t
}

func buildTimeRangeState(fromStr, toStr string, granularity time.Duration, order collection_state.CollectionOrder, endObjects ...string) *collection_state.TimeRangeObjectState {
	from, err := time.Parse("2006-01-02 15:04:05", fromStr)
	if err != nil {
		panic(err)
	}
	to, err := time.Parse("2006-01-02 15:04:05", toStr)
	if err != nil {
		panic(err)
	}
	endObjectsMap := make(map[string]struct{})
	for _, obj := range endObjects {
		endObjectsMap[obj] = struct{}{}
	}
	return &collection_state.TimeRangeObjectState{
		TimeRange: collection_state.DirectionalTimeRange{
			From:            from,
			To:              to,
			CollectionOrder: order,
		},
		EndObjects:  endObjectsMap,
		Granularity: granularity,
	}
}
