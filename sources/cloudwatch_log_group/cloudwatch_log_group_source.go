// Package cloudwatch provides functionality to collect logs from AWS CloudWatch Log Groups
//
// This package enables the collection of log events from AWS CloudWatch log groups, supporting incremental collection, filtering, and batching for efficient processing.
package cloudwatch_log_group

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const (
	// AwsCloudwatchLogGroupSourceIdentifier is the unique identifier for the CloudWatch log group source
	AwsCloudwatchLogGroupSourceIdentifier = "aws_cloudwatch_log_group"
)

// AwsCloudWatchLogGroupSource is responsible for collecting events from log streams within a CloudWatch log group.
// It implements the RowSource interface and manages collection state to support incremental and efficient log collection.
type AwsCloudWatchLogGroupSource struct {
	// Embeds the base RowSourceImpl with CloudWatch-specific config and AWS connection.
	row_source.RowSourceImpl[*AwsCloudWatchLogGroupSourceConfig, *config.AwsConnection]

	// client is the AWS CloudWatch Logs client used for API calls.
	client *cloudwatchlogs.Client
	// errorList accumulates errors encountered during collection for reporting.
	errorList []error
	// state tracks progress and supports incremental collection across log streams.
	//state *CloudWatchLogGroupCollectionState
}

// Init sets up the CloudWatch log group source with the provided parameters and options.
// It initializes the collection state, AWS client, and validates the configuration.
// If a specific start time is provided, it clears the previous collection state to force recollection.
func (s *AwsCloudWatchLogGroupSource) Init(ctx context.Context, params *row_source.RowSourceParams, opts ...row_source.RowSourceOption) error {
	// Set up the collection state constructor
	s.NewCollectionStateFunc = func() collection_state.CollectionState {
		return NewCloudWatchLogGroupCollectionState()
	}

	// NOTE: set the granularity to be 1 minute
	// (we actually set a func on our base RowSourceImpl to get the granularity
	// this is to avoid an initialisation ordering issue when setting artifact source granularity)
	s.RowSourceImpl.GetGranularityFunc = s.getGranularity

	// Initialize the base implementation
	if err := s.RowSourceImpl.Init(ctx, params, opts...); err != nil {
		return err
	}

	// Initialize AWS CloudWatch client
	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}

	s.client = client
	s.errorList = []error{}

	return nil
}

// getGranularity returns the granularity for this source type, which is set to 1 minute.
func (s *AwsCloudWatchLogGroupSource) getGranularity() time.Duration {
	return time.Minute
}

// Identifier returns the unique identifier for this source type, used in the plugin system.
func (s *AwsCloudWatchLogGroupSource) Identifier() string {
	return AwsCloudwatchLogGroupSourceIdentifier
}

// matchesAnyPattern returns true if the target string matches any of the provided patterns (supports wildcards).
func matchesAnyPattern(target string, patterns []string) bool {
	for _, pattern := range patterns {
		match, err := filepath.Match(pattern, target)
		if err != nil {
			slog.Error("error matching pattern", "pattern", pattern, "error", err)
			continue
		}
		if match {
			return true
		}
	}
	return false
}

// Collect retrieves log events from CloudWatch log streams within the specified time range.
//
// This function is responsible for collecting log events from all relevant log streams in the configured CloudWatch log group.
// The process includes:
//  1. Retrieving all log streams that match the configuration (optionally filtered by name/pattern).
//  2. Batching log streams to efficiently query events in groups (up to 100 at a time).
//  3. For each batch, querying CloudWatch Logs for events within the desired time window.
//  4. Sorting and processing each event, skipping already-collected events based on collection state.
//  5. Enriching and forwarding each new event for downstream processing.
//  6. Updating the collection state to support incremental collection and avoid duplicates.
//  7. Aggregating and returning any errors encountered during the process.
//
// Returns an error if any step fails, or if errors are encountered during log collection.
func (s *AwsCloudWatchLogGroupSource) Collect(ctx context.Context) error {
	// Get all log streams matching the prefix in the specified log group
	logStreamCollection, err := s.getLogStreamsToCollect(ctx)
	if err != nil {
		return fmt.Errorf("failed to collect log streams, %w", err)
	}

	slog.Debug("Total log stream collected based on '--from' flag",
		"count", len(logStreamCollection),
		"log_group", s.Config.LogGroupName)

	// Filter out the log streams that are not in the list of log stream names
	if len(s.Config.LogStreamNames) > 0 {
		filteredLogStreamCollection := []cwTypes.LogStream{}
		logStreamNames := s.Config.LogStreamNames

		for _, ls := range logStreamCollection {
			if ls.LogStreamName == nil {
				s.errorList = append(s.errorList, fmt.Errorf("skipping stream with nil name in log group %s", s.Config.LogGroupName))
				continue
			}

			if matchesAnyPattern(*ls.LogStreamName, logStreamNames) {
				filteredLogStreamCollection = append(filteredLogStreamCollection, ls)
			}
		}

		// Use the filtered collection
		logStreamCollection = filteredLogStreamCollection
	}

	slog.Info("Starting collection", "total_streams", len(logStreamCollection))

	batchLogStream := [][]string{}

	streamNames := []string{}

	for _, stream := range logStreamCollection {
		streamNames = append(streamNames, *stream.LogStreamName)
	}

	// Loop through the array with step 2 to create sub-arrays of size 100
	for i := 0; i < len(streamNames); i += 100 {
		// Append a sub-array of size 100 (or less if the last one is incomplete)
		end := i + 100
		if end > len(streamNames) {
			end = len(streamNames)
		}
		batchLogStream = append(batchLogStream, streamNames[i:end])
	}

	batchCount := 0
	for _, batch := range batchLogStream {
		batchCount++
		slog.Info("Processing batch log streams",
			"batch", batchCount,
			"log_group", s.Config.LogGroupName)
		// Convert time range to milliseconds for CloudWatch API
		startTimeMillis := s.FromTime.UnixMilli()
		endTimeMillis := time.Now().UnixMilli()

		input := &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:   &s.Config.LogGroupName,
			LogStreamNames: batch,
			StartTime:      aws.Int64(startTimeMillis),
			EndTime:        aws.Int64(endTimeMillis),
		}

		events, err := s.filterLogEvents(ctx, input)
		if err != nil {
			s.errorList = append(s.errorList, fmt.Errorf("failed to filter log events for stream %s: %w", batch, err))
			continue
		}

		events = sortFilteredLogEvents(events)

		// Process each event in the page
		for _, event := range events {
			if event.Message == nil || *event.Message == "" {
				s.errorList = append(s.errorList, fmt.Errorf("empty message in stream %s at timestamp %d", *event.LogStreamName, *event.Timestamp))
				continue
			}

			slog.Info("Processing stream", "stream", *event.LogStreamName)
			// Set up source enrichment fields for the current stream
			sourceEnrichmentFields := &schema.SourceEnrichment{
				CommonFields: schema.CommonFields{
					TpSourceType:     AwsCloudwatchLogGroupSourceIdentifier,
					TpSourceName:     &s.Config.LogGroupName,
					TpSourceLocation: event.LogStreamName,
				},
			}

			timestamp := time.UnixMilli(*event.Timestamp)
			// Skip already collected events based on state
			if !s.CollectionState.ShouldCollect(*event.LogStreamName, timestamp) {
				slog.Debug("Skipping already collected event",
					"stream", *event.LogStreamName,
					"timestamp", timestamp.Format(time.RFC3339))
				continue
			}

			row := &types.RowData{
				Data:             event,
				SourceEnrichment: sourceEnrichmentFields,
			}

			// Update collection state with the processed event
			if err := s.CollectionState.OnCollected(*event.LogStreamName, timestamp); err != nil {
				s.errorList = append(s.errorList, fmt.Errorf("failed to update collection state for stream %s: %w", *event.LogStreamName, err))
				continue
			}

			// Send the row for processing
			if err := s.OnRow(ctx, row); err != nil {
				s.errorList = append(s.errorList, fmt.Errorf("error processing row in stream %s: %w", *event.LogStreamName, err))
				continue
			}
		}
	}

	// Return collected errors if any
	if len(s.errorList) > 0 {
		return fmt.Errorf("encountered %d errors during log collection: %v", len(s.errorList), s.errorList)
	}

	return nil
}

// filterLogEvents retrieves all log events for the given input, handling pagination.
// Returns a slice of FilteredLogEvent and any error encountered.
func (s *AwsCloudWatchLogGroupSource) filterLogEvents(ctx context.Context, input *cloudwatchlogs.FilterLogEventsInput) ([]cwTypes.FilteredLogEvent, error) {

	allEvents := []cwTypes.FilteredLogEvent{}

	paginator := cloudwatchlogs.NewFilterLogEventsPaginator(s.client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		allEvents = append(allEvents, output.Events...)
	}

	return allEvents, nil
}

// sortFilteredLogEvents sorts a slice of FilteredLogEvent by LogStreamName and Timestamp.
// This ensures events are processed in a consistent order.
func sortFilteredLogEvents(events []cwTypes.FilteredLogEvent) []cwTypes.FilteredLogEvent {
	// Sorting the events by LogStreamName and then by Timestamp
	sort.Slice(events, func(i, j int) bool {
		// First, sort by LogStreamName, then by Timestamp
		if *events[i].LogStreamName != *events[j].LogStreamName {
			return *events[i].LogStreamName < *events[j].LogStreamName
		}
		return *events[i].Timestamp < *events[j].Timestamp
	})

	return events
}

// getLogStreamsToCollect retrieves all log streams in a log group that match the specified prefix.
// It paginates through the DescribeLogStreams API and stops when streams are older than the configured start time.
// Returns a sorted slice of log streams from oldest to newest.
func (s *AwsCloudWatchLogGroupSource) getLogStreamsToCollect(ctx context.Context) ([]cwTypes.LogStream, error) {
	var logStreams []cwTypes.LogStream
	var nextToken *string

	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &s.Config.LogGroupName,
		NextToken:    nextToken,
		OrderBy:      cwTypes.OrderByLastEventTime,
		Descending:   aws.Bool(true),
	}

	paginator := cloudwatchlogs.NewDescribeLogStreamsPaginator(s.client, input, func(o *cloudwatchlogs.DescribeLogStreamsPaginatorOptions) {
		o.Limit = int32(50)
		o.StopOnDuplicateToken = true
	})

	// Flag to indicate whether to stop pagination
	stopPagination := false

	// Collect the log streams first
	for paginator.HasMorePages() && !stopPagination {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to describe log streams, %w", err)
		}

		for _, ls := range output.LogStreams {
			// Skip if LastEventTimestamp is nil
			if ls.LastEventTimestamp == nil {
				continue
			}

			// Convert LastEventTimestamp once and reuse
			lastEventTime := time.UnixMilli(*ls.LastEventTimestamp)

			// If LastEventTimestamp is before FromTime, break the loop
			if lastEventTime.Before(s.FromTime) {
				stopPagination = true
				break
			}

			// Add log stream to the collection
			logStreams = append(logStreams, ls)
		}

		// Check if we need to stop pagination
		if stopPagination {
			slog.Debug("Stopping pagination as lastEventTime is before FromTime",
				"log_group", s.Config.LogGroupName)
			break
		}
	}

	// If no log streams were collected, log the information and return
	if len(logStreams) == 0 {
		slog.Info("No log streams found to collect", "logGroupName", s.Config.LogGroupName)
		return nil, nil
	}

	// Sort the log streams from oldest to newest based on LastEventTimestamp
	sort.Slice(logStreams, func(i, j int) bool {
		// Ensure nil values come last
		if logStreams[i].LastEventTimestamp == nil || logStreams[j].LastEventTimestamp == nil {
			return logStreams[i].LastEventTimestamp == nil
		}
		// Sort by LastEventTimestamp
		return *logStreams[i].LastEventTimestamp < *logStreams[j].LastEventTimestamp
	})

	return logStreams, nil
}

// getClient initializes and returns an AWS CloudWatch Logs client for the configured region.
// Returns an error if the client cannot be created.
func (s *AwsCloudWatchLogGroupSource) getClient(ctx context.Context) (*cloudwatchlogs.Client, error) {
	region := s.Config.Region

	cfg, err := s.Connection.GetClientConfiguration(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("failed to get client configuration, %w", err)
	}

	client := cloudwatchlogs.NewFromConfig(*cfg)
	return client, nil
}
