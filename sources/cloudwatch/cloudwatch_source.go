package cloudwatch

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const (
	AwsCloudwatchSourceIdentifier = "aws_cloudwatch_log_group"
	defaultCloudwatchRegion       = "us-east-1"
)

// AwsCloudWatchSource is responsible for collection of events from log streams within a log group in AWS CloudWatch
type AwsCloudWatchSource struct {
	row_source.RowSourceImpl[*AwsCloudWatchSourceConfig, *config.AwsConnection]

	client    *cloudwatchlogs.Client
	errorList []error
	state     *CloudWatchCollectionState
}

func (s *AwsCloudWatchSource) Init(ctx context.Context, params *row_source.RowSourceParams, opts ...row_source.RowSourceOption) error {
	// set the collection state ctor
	s.NewCollectionStateFunc = func() collection_state.CollectionState[*AwsCloudWatchSourceConfig] {
		return NewCloudWatchCollectionState()
	}

	if err := s.RowSourceImpl.Init(ctx, params, opts...); err != nil {
		return err
	}

	// initialize client
	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}

	s.client = client
	s.errorList = []error{}

	// Get the collection state from the base implementation
	state, ok := s.CollectionState.(*CloudWatchCollectionState)
	if !ok {
		return fmt.Errorf("invalid collection state type: expected *CloudWatchCollectionState")
	}
	s.state = state

	return nil
}

func (s *AwsCloudWatchSource) Identifier() string {
	return AwsCloudwatchSourceIdentifier
}

func (s *AwsCloudWatchSource) Collect(ctx context.Context) error {
	// obtain log streams which have active events in the time range
	logStreamCollection, err := s.getLogStreamsToCollect(ctx, s.Config.LogGroupName, s.Config.LogStreamPrefix)
	if err != nil {
		return fmt.Errorf("failed to collect log streams, %w", err)
	}

	slog.Info("Starting collection", "total_streams", len(logStreamCollection))

	// collect events from each log stream
	for _, ls := range logStreamCollection {
		if ls.LogStreamName == nil {
			continue
		}

		slog.Info("Processing stream", "stream", *ls.LogStreamName)
		sourceEnrichmentFields := &schema.SourceEnrichment{
			CommonFields: schema.CommonFields{
				TpSourceType:     AwsCloudwatchSourceIdentifier,
				TpSourceName:     &s.Config.LogGroupName,
				TpSourceLocation: ls.LogStreamName,
			},
		}

		// Convert Unix seconds to milliseconds for CloudWatch API
		startTimeMillis := s.FromTime.UnixMilli()
		endTimeMillis := time.Now().UnixMilli()

		input := &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  &s.Config.LogGroupName,
			LogStreamName: ls.LogStreamName,
			StartFromHead: aws.Bool(true),
			StartTime:     aws.Int64(startTimeMillis),
			EndTime:       aws.Int64(endTimeMillis),
			Limit:         aws.Int32(10000),
		}

		var (
			nextToken   *string
			totalEvents int
		)

		paginator := cloudwatchlogs.NewGetLogEventsPaginator(s.client, input)
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				return fmt.Errorf("failed to get log events for stream %s: %w", *ls.LogStreamName, err)
			}

			// Break if no events in this page
			if len(output.Events) == 0 {
				slog.Debug("No events in page", "stream", *ls.LogStreamName)
				break
			}

			for _, event := range output.Events {
				if event.Message == nil || *event.Message == "" {
					continue
				}

				timestamp := time.UnixMilli(*event.Timestamp)

				// Skip this event if we've already collected it
				if !s.state.ShouldCollect(*ls.LogStreamName, timestamp) {
					slog.Debug("Skipping already collected event",
						"stream", *ls.LogStreamName,
						"timestamp", timestamp.Format(time.RFC3339))
					continue
				}

				row := &types.RowData{
					Data:             event.Message,
					SourceEnrichment: sourceEnrichmentFields,
				}

				// update collection state
				if err := s.state.OnCollected(*ls.LogStreamName, timestamp); err != nil {
					return fmt.Errorf("failed to update collection state: %w", err)
				}

				if err := s.OnRow(ctx, row); err != nil {
					return fmt.Errorf("error processing row: %w", err)
				}

				totalEvents++
			}

			// Break if token hasn't changed (to avoid infinite loop)
			if nextToken != nil && *output.NextForwardToken == *nextToken {
				slog.Debug("Stopping: NextForwardToken hasn't changed",
					"stream", *ls.LogStreamName,
					"total_events", totalEvents)
				break
			}
			nextToken = output.NextForwardToken
		}
	}

	return nil
}

// Get log steams for the specified log group
func (s *AwsCloudWatchSource) getLogStreamsToCollect(ctx context.Context, logGroupName string, logStreamPrefix *string) ([]cwtypes.LogStream, error) {
	var logStreams []cwtypes.LogStream
	var nextToken *string

	for {
		input := &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        &logGroupName,
			LogStreamNamePrefix: logStreamPrefix,
			NextToken:           nextToken,
		}

		output, err := s.client.DescribeLogStreams(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to describe log streams, %w", err)
		}

		// If no log streams found, break
		if len(output.LogStreams) == 0 {
			slog.Debug("No log streams found", "logGroupName", logGroupName, "logStreamPrefix", logStreamPrefix)
			break
		}

		logStreams = append(logStreams, output.LogStreams...)

		if output.NextToken == nil {
			slog.Debug("No more log streams to fetch", "logGroupName", logGroupName, "logStreamPrefix", logStreamPrefix, "total_streams", len(logStreams))
			break
		}
		nextToken = output.NextToken
	}

	if len(logStreams) == 0 {
		slog.Info("No log streams found to collect", "logGroupName", logGroupName, "logStreamPrefix", logStreamPrefix)
	}

	return logStreams, nil
}

func (s *AwsCloudWatchSource) getClient(ctx context.Context) (*cloudwatchlogs.Client, error) {
	tempRegion := defaultCloudwatchRegion
	if s.Config != nil && s.Config.Region != nil {
		tempRegion = *s.Config.Region
	}

	cfg, err := s.Connection.GetClientConfiguration(ctx, &tempRegion)
	if err != nil {
		return nil, fmt.Errorf("failed to get client configuration, %w", err)
	}

	client := cloudwatchlogs.NewFromConfig(*cfg)
	return client, nil
}
