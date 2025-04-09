// Package cloudwatch provides functionality to collect logs from AWS CloudWatch Log Groups
package cloudwatch_log_group

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
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

// AwsCloudWatchLogGroupSource is responsible for collection of events from log streams within a log group in AWS CloudWatch
// It implements the RowSource interface and manages the collection state to support incremental collection
type AwsCloudWatchLogGroupSource struct {
	// Embed the base RowSourceImpl with CloudWatch specific config and AWS connection
	row_source.RowSourceImpl[*AwsCloudWatchLogGroupSourceConfig, *config.AwsConnection]

	// AWS CloudWatch Logs client
	client *cloudwatchlogs.Client
	// List of errors encountered during collection
	errorList []error
	// Collection state to track progress and support incremental collection
	state *CloudWatchLogGroupCollectionState
}

// Init initializes the CloudWatch source with the provided parameters and options
// It sets up the collection state, AWS client, and validates the configuration
func (s *AwsCloudWatchLogGroupSource) Init(ctx context.Context, params *row_source.RowSourceParams, opts ...row_source.RowSourceOption) error {
	// Set up the collection state constructor
	s.NewCollectionStateFunc = func() collection_state.CollectionState[*AwsCloudWatchLogGroupSourceConfig] {
		return NewCloudWatchLogGroupCollectionState()
	}

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

	// Get and validate the collection state from the base implementation
	state, ok := s.CollectionState.(*CloudWatchLogGroupCollectionState)
	if !ok {
		return fmt.Errorf("invalid collection state type: expected *CloudWatchCollectionState")
	}
	s.state = state

	return nil
}

// Identifier returns the unique identifier for this source
func (s *AwsCloudWatchLogGroupSource) Identifier() string {
	return AwsCloudwatchLogGroupSourceIdentifier
}

func matchesAnyPattern(target string, patterns []string) bool {
	for _, pattern := range patterns {
		match, err := filepath.Match(pattern, target)
		if err != nil {
			slog.Error("error matching pattern '%s': %v\n", pattern, err)
			continue
		}
		if match {
			return true
		}
	}
	return false
}

// Collect retrieves log events from CloudWatch log streams within the specified time range
// It handles pagination, maintains collection state, and processes events incrementally
func (s *AwsCloudWatchLogGroupSource) Collect(ctx context.Context) error {
	// Get all log streams matching the prefix in the specified log group
	logStreamCollection, err := s.getLogStreamsToCollect(ctx)
	if err != nil {
		return fmt.Errorf("failed to collect log streams, %w", err)
	}

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

	// Process each log stream
	for _, ls := range logStreamCollection {
		if ls.LogStreamName == nil {
			s.errorList = append(s.errorList, fmt.Errorf("skipping stream with nil name in log group %s", s.Config.LogGroupName))
			continue
		}

		slog.Info("Processing stream", "stream", *ls.LogStreamName)
		// Set up source enrichment fields for the current stream
		sourceEnrichmentFields := &schema.SourceEnrichment{
			CommonFields: schema.CommonFields{
				TpSourceType:     AwsCloudwatchLogGroupSourceIdentifier,
				TpSourceName:     &s.Config.LogGroupName,
				TpSourceLocation: ls.LogStreamName,
			},
		}

		// Convert time range to milliseconds for CloudWatch API
		startTimeMillis := s.FromTime.UnixMilli()
		endTimeMillis := time.Now().UnixMilli()

		// Configure the GetLogEvents API request
		input := &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  &s.Config.LogGroupName,
			LogStreamName: ls.LogStreamName,
			StartFromHead: aws.Bool(true), // Start from oldest events
			StartTime:     aws.Int64(startTimeMillis),
			EndTime:       aws.Int64(endTimeMillis),
			Limit:         aws.Int32(10000), // Maximum allowed by AWS API
		}

		// For incremental collection, start from the last collected event time
		if s.state.GetEndTimeForStream(*ls.LogStreamName).UnixMilli() > startTimeMillis {
			input.StartTime = aws.Int64(s.state.GetEndTimeForStream(*ls.LogStreamName).UnixMilli())
		}

		var (
			nextToken   *string
			totalEvents int
		)

		// Use paginator to handle response pagination automatically
		paginator := cloudwatchlogs.NewGetLogEventsPaginator(s.client, input)
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				s.errorList = append(s.errorList, fmt.Errorf("failed to get log events for stream %s: %w", *ls.LogStreamName, err))
				break // Skip to next stream on error
			}

			// Break if no events in this page
			if len(output.Events) == 0 {
				slog.Debug("No events in page", "stream", *ls.LogStreamName)
				break
			}

			// Process each event in the page
			for _, event := range output.Events {
				if event.Message == nil || *event.Message == "" {
					s.errorList = append(s.errorList, fmt.Errorf("empty or nil message in stream %s at timestamp %d", *ls.LogStreamName, *event.Timestamp))
					continue
				}

				timestamp := time.UnixMilli(*event.Timestamp)

				// Skip already collected events based on state
				if !s.state.ShouldCollect(*ls.LogStreamName, timestamp) {
					slog.Debug("Skipping already collected event",
						"stream", *ls.LogStreamName,
						"timestamp", timestamp.Format(time.RFC3339))
					continue
				}

				// Create row data with the event message and enrichment
				row := &types.RowData{
					Data:             event.Message,
					SourceEnrichment: sourceEnrichmentFields,
				}

				// Update collection state with the processed event
				if err := s.state.OnCollected(*ls.LogStreamName, timestamp); err != nil {
					s.errorList = append(s.errorList, fmt.Errorf("failed to update collection state for stream %s: %w", *ls.LogStreamName, err))
					continue
				}

				// Send the row for processing
				if err := s.OnRow(ctx, row); err != nil {
					s.errorList = append(s.errorList, fmt.Errorf("error processing row in stream %s: %w", *ls.LogStreamName, err))
					continue
				}

				totalEvents++
			}

			// Break if we've received the same token twice (end of stream)
			if nextToken != nil && *output.NextForwardToken == *nextToken {
				slog.Debug("Stopping: NextForwardToken hasn't changed",
					"stream", *ls.LogStreamName,
					"total_events", totalEvents)
				break
			}
			nextToken = output.NextForwardToken
		}
	}

	// Return collected errors if any
	if len(s.errorList) > 0 {
		return fmt.Errorf("encountered %d errors during log collection: %v", len(s.errorList), s.errorList)
	}

	return nil
}

// getLogStreamsToCollect retrieves all log streams in a log group that match the specified prefix
// It handles pagination of the DescribeLogStreams API response
func (s *AwsCloudWatchLogGroupSource) getLogStreamsToCollect(ctx context.Context) ([]cwTypes.LogStream, error) {
	var logStreams []cwTypes.LogStream
	var nextToken *string

	for {
		input := &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: &s.Config.LogGroupName,
			NextToken:    nextToken,
		}

		output, err := s.client.DescribeLogStreams(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to describe log streams, %w", err)
		}

		// Break if no streams found in this page
		if len(output.LogStreams) == 0 {
			slog.Debug("No log streams found", "logGroupName", s.Config.LogGroupName)
			break
		}

		logStreams = append(logStreams, output.LogStreams...)

		// Break if no more pages
		if output.NextToken == nil {
			slog.Debug("No more log streams to fetch", "logGroupName", s.Config.LogGroupName, "total_streams", len(logStreams))
			break
		}
		nextToken = output.NextToken
	}

	if len(logStreams) == 0 {
		slog.Info("No log streams found to collect", "logGroupName", s.Config.LogGroupName)
	}

	return logStreams, nil
}

// getClient initializes and returns an AWS CloudWatch Logs client
// It uses the provided region or falls back to the default region
func (s *AwsCloudWatchLogGroupSource) getClient(ctx context.Context) (*cloudwatchlogs.Client, error) {
	region := s.Config.Region

	cfg, err := s.Connection.GetClientConfiguration(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("failed to get client configuration, %w", err)
	}

	client := cloudwatchlogs.NewFromConfig(*cfg)
	return client, nil
}
