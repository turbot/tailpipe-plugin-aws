package cloudwatch

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	typehelpers "github.com/turbot/go-kit/types"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

const (
	AwsCloudwatchSourceIdentifier = "aws_cloudwatch_log_group"
	defaultCloudwatchRegion         = "us-east-1"
)

// register the source from the package init function
func init() {
	row_source.RegisterRowSource[*AwsCloudWatchSource]()
}

// AwsCloudWatchSource is responsible for collection of events from log streams within a log group in AWS CloudWatch
type AwsCloudWatchSource struct {
	row_source.RowSourceImpl[*AwsCloudWatchSourceConfig, *config.AwsConnection]

	client *cloudwatchlogs.Client
}

func (s *AwsCloudWatchSource) Init(ctx context.Context, params *row_source.RowSourceParams, opts ...row_source.RowSourceOption) error {
	slog.Info("Initializing AwsCloudwatchSource")
	// set the collection state ctor
	
	s.NewCollectionStateFunc = NewAwsCloudwatchCollectionState

	// call base init to set config/connection
	err := s.RowSourceImpl.Init(ctx, params, opts...)
	if err != nil {
		return fmt.Errorf("failed to init %s source, %w", AwsCloudwatchSourceIdentifier, err)
	}

	// initialize client
	client, err := s.getClient(ctx)
	if err != nil {
		return err
	}
	s.client = client

	return nil
}

func (s *AwsCloudWatchSource) Identifier() string {
	return AwsCloudwatchSourceIdentifier
}

func (s *AwsCloudWatchSource) Collect(ctx context.Context) error {
	collectionState := s.CollectionState.(*AwsCloudwatchCollectionState)

	// obtain log streams which have active events in the time range
	logStreamCollection, err := s.collectLogStreams(ctx, s.Config.LogGroupName, s.Config.LogStreamPrefix, collectionState)
	if err != nil {
		return fmt.Errorf("failed to collect log streams, %w", err)
	}

	// collect events from each log stream
	for _, ls := range logStreamCollection {
		sourceEnrichmentFields := &schema.SourceEnrichment{
			CommonFields: schema.CommonFields{
				TpSourceType:     AwsCloudwatchSourceIdentifier,
				TpSourceName:     &s.Config.LogGroupName,
				TpSourceLocation: ls.LogStream.LogStreamName,
			},
		}

		if ls.StartTime >= ls.EndTime {
			slog.Warn("log stream %s has invalid time range, skipping", ls.LogStream.LogStreamName)
			continue
		}

		// To ensure smoother execution, we have set the value to 7000, even though the maximum allowable limit is 10000.
		// https://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_GetLogEvents.html#API_GetLogEvents_RequestSyntax
		var pageSize int32 = 7000
		var nextToken *string

		input := &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  &s.Config.LogGroupName,
			LogStreamName: ls.LogStream.LogStreamName,
			StartFromHead: aws.Bool(true),
			StartTime:     &ls.StartTime,
			EndTime:       &ls.EndTime,
			Limit:         &pageSize,
		}

		paginator := cloudwatchlogs.NewGetLogEventsPaginator(s.client, input)
		for paginator.HasMorePages() {
			var output *cloudwatchlogs.GetLogEventsOutput
			output, err = paginator.NextPage(ctx)
			if err != nil {
				return fmt.Errorf("failed to get log events, %w", err)
			}

			for _, event := range output.Events {
				if event.Message == nil || *event.Message == "" {
					continue
				}
				row := &types.RowData{
					Data:             event.Message,
					SourceEnrichment: sourceEnrichmentFields,
				}

				// update collection state
				// build time from unix mill
				unixMillis := *event.Timestamp
				timestamp := time.Unix(0, unixMillis*int64(time.Millisecond))

				err := collectionState.OnCollected(*ls.LogStream.LogStreamName, timestamp)
				if err != nil {
					return fmt.Errorf("failed to update collection state, %w", err)
				}

				if err := s.OnRow(ctx, row); err != nil {
					return fmt.Errorf("error processing row: %w", err)
				}
			}

			if nextToken != nil && output.NextForwardToken != nil && *nextToken == *output.NextForwardToken {
				break
			}

			nextToken = output.NextForwardToken
		}
	}

	return nil
}

func (s *AwsCloudWatchSource) collectLogStreams(ctx context.Context, logGroupName string, logStreamPrefix *string, collectionState *AwsCloudwatchCollectionState) ([]logStreamsToCollect, error) {
	var logStreams []logStreamsToCollect
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

		for _, logStream := range output.LogStreams {
			streamName := typehelpers.StringValue(logStream.LogStreamName)
			start, end := s.getTimeRange(streamName, collectionState)
			if s.logStreamHasEntriesInTimeRange(logStream, start, end) {
				logStreams = append(logStreams, logStreamsToCollect{LogStream: logStream, StartTime: start, EndTime: end})
			}
		}

		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return logStreams, nil
}

func (s *AwsCloudWatchSource) logStreamHasEntriesInTimeRange(logStream cwtypes.LogStream, startTime, endTime int64) bool {
	if logStream.LastIngestionTime == nil || logStream.FirstEventTimestamp == nil {
		return false
	}
	if startTime >= endTime {
		return false
	}
	return *logStream.LastIngestionTime > startTime && *logStream.FirstEventTimestamp < endTime
}

// use the collection state data (if present) and the configured time range to determine the start and end time
func (s *AwsCloudWatchSource) getTimeRange(logStream string, collectionState *AwsCloudwatchCollectionState) (int64, int64) {
	startTime := s.FromTime.UnixMilli()
	endTime := time.Now().UnixMilli()

	if collectionState != nil {
		// set start time from collection state data if present
		if prevTimestamp, ok := collectionState.LogStreamTimestamps[logStream]; ok {
			return prevTimestamp.Add(time.Microsecond).UnixMilli(), endTime
		}
	}
	return startTime, endTime
}

func (s *AwsCloudWatchSource) getClient(ctx context.Context) (*cloudwatchlogs.Client, error) {
	tempRegion := defaultCloudwatchRegion
	if s.Config.Region != nil {
		tempRegion = *s.Config.Region
	}

	cfg, err := s.Connection.GetClientConfiguration(ctx, &tempRegion)
	if err != nil {
		return nil, fmt.Errorf("failed to get client configuration, %w", err)
	}

	client := cloudwatchlogs.NewFromConfig(*cfg)
	return client, nil
}

type logStreamsToCollect struct {
	LogStream cwtypes.LogStream
	StartTime int64
	EndTime   int64
}