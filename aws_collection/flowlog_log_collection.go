package aws_collection

import (
	"context"
	"fmt"
	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/artifact"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/paging"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"time"
	//"github.com/turbot/tailpipe-plugin-sdk/collection"
	//sdkconfig "github.com/turbot/tailpipe-plugin-sdk/config"
	//"github.com/turbot/tailpipe-plugin-sdk/source"
)

type FlowlogLogCollection struct {
	// all collections must embed collection.Base
	collection.Base

	// the collection config
	Config *FlowLogCollectionConfig
}

func NewFlowlogLogCollection() plugin.Collection {
	l := &FlowlogLogCollection{}
	return l
}

// GetRowStruct implements Collection
// return an instance of the row struct
func (c *FlowlogLogCollection) GetRowStruct() any {
	return aws_types.FlowLog{}
}

// Init implements Collection
func (c *FlowlogLogCollection) Init(ctx context.Context, configData []byte) error {
	// TEMP - this will actually parse (or the base will)
	// unmarshal the config
	config := &FlowLogCollectionConfig{
		Paths: []string{"/Users/kai/tailpipe_data/flowlog"},
		// use defaults when using cloudwatch
		//		//Fields: []string{"timestamp",
		//		//	"version",
		//		//	"account-id",
		//		//	"interface-id",
		//		//	"srcaddr",
		//		//	"dstaddr",
		//		//	"srcport",
		//		//	"dstport",
		//		//	"protocol",
		//		//	"packets",
		//		//	"bytes",
		//		//	"start",
		//		//	"end",
		//		//	"action",
		//		//	"log-status",
		//		//},
	}

	//err := json.Unmarshal(configData, config)
	//if err != nil {
	//	return fmt.Errorf("error unmarshalling config: %w", err)
	//}

	// todo - parse config as hcl
	c.Config = config
	// todo validate config

	// init the config = this will set default fields if needed
	if err := c.Config.Init(); err != nil {
		return fmt.Errorf("error initializing config: %w", err)
	}
	// todo validate config

	// todo create source from config
	source, err := c.getSource(ctx, c.Config)
	if err != nil {
		return err
	}
	return c.AddSource(source)
}

func (c *FlowlogLogCollection) getSource(ctx context.Context, config *FlowLogCollectionConfig) (plugin.RowSource, error) {

	// TODO populate from config
	//sourceConfig := &artifact.AwsS3BucketSourceConfig{
	//	Bucket:       "silverwater-flowlog-s3-bucket",
	//	Extensions:   []string{".gz"},
	//	AccessKey:    "",
	//	SecretKey:    "",
	//	SessionToken: "",
	//}
	//
	//artifactSource, err := artifact.NewAwsS3BucketSource(sourceConfig)
	//if err != nil {
	//	return nil, fmt.Errorf("error creating s3 bucket source: %w", err)
	//}

	//artifactSource := artifact.NewFileSystemSource(&artifact.FileSystemSourceConfig{
	//	Paths:      config.Paths,
	//	Extensions: []string{".gz"},
	//})

	artifactSource, err := artifact.NewAwsCloudWatchSource(ctx, &artifact.AwsCloudWatchSourceConfig{
		AccessKey:       "ASIARNKUQPUTWCYZF7H3",
		SecretKey:       "mKasF5qHQ0ejKr5xAC0xh4Iz6shFeNvr0pJGuVmW",
		SessionToken:    "IQoJb3JpZ2luX2VjEFQaCXVzLWVhc3QtMiJIMEYCIQCQZSo4wDMbOkp3tDJIT7N825klJrP5G0MNyUCrMsVLKgIhALQ07bWzfoNgUxHcQpxn15w2dwZUK5awABAUmKcXSnXgKoYDCC0QAxoMMDk3MzUwODc2NDU1Igz/AKtTka19aFUqgngq4wIsRk/Mu/shGA9zsHCjd8lQAHuazr9+6LKIdGstvfsk4BaYnb/WGAHcEOMhKfkpqmqolHB4OTCZWKj5fcQKwx+UfwvtqZ7DktbIJDlYcUm0Rc7Gyn5W+FvcxCWQLUCwEs3uCVH7sz3zoszJCI8OjRkrhVXMKUsRS5wZ8PYaVheiMCN7319UQcj7v8x3SLr2ex4aF3xXmHgAJYV0t6Y4SiHKYZXL3t/JfP5hvG1/DNHbryfYIaRybyl3ulpXZ19jU77Ki1yjSQXLlD7sbR2STtx1hEUjmMM0z0dOfIrTwVbi7eaf4i0NKGG1FxsM7p2oqXSSjya6INzfZlJixw+KETo8L2Y0CwVycxqrYNxDjJX+wlBBVXFSKxXtljQl3+bCmXNAmixnsdkdDzhDsCCVeshz3RR+omFzkEFstCoSjH1XtQyfUrMBAYe0zOGklo7Pi1voXJp8v4B4iZqhAhz+2MUr9TqPMPSv6bQGOqUBpHVwsRrK1+a9JYetSw6c5A43jli1PtaxoAdGQtMCfSWq7SQf2EywiZJA9kwlILOumObODQZQwuBXPMZ+6dwuaeQIlUsUPW4DU/F2whyTTBfgYwyOw/wZTnIOOdCaR04JcUZIla6ty+cBduHhN0w8TO5XdCSqzNH9rMsDYQgcZD+o+UDVOqhIj6Cg0cZrxzhrC2AEKI8e/CzuP51gl3HUonvRSwuQ",
		LogGroupName:    "/victor/vpc/flowlog",
		LogStreamPrefix: utils.ToStringPointer("eni-000b"),
		StartTime:       time.Now().Add(-time.Hour * 24),
		EndTime:         time.Now(),
	})

	// create empty paging data to pass to source
	// TODO maybe source creates for itself??
	pagingData, err := c.NewPagingData()
	if err != nil {
		return nil, fmt.Errorf("error creating paging data: %w", err)
	}

	source, err := row_source.NewArtifactRowSource(
		artifactSource,
		pagingData,
		// we expect a log row per line of log data
		row_source.WithRowPerLine(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating artifact row source: %w", err)
	}

	return source, nil
}

// Identifier implements Collection
func (c *FlowlogLogCollection) Identifier() string {
	return "aws_flow_log"
}

// EnrichRow implements RowEnricher
func (c *FlowlogLogCollection) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// row must be a string
	rowString, ok := row.(string)
	if !ok {
		return nil, fmt.Errorf("invalid row type %T, expected string", row)
	}
	record, err := aws_types.FlowLogFromString(rowString, c.Config.Fields)

	if err != nil {
		return nil, fmt.Errorf("error parsing row: %s", err)
	}

	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		record.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	record.TpID = xid.New().String()

	// TODO is source type actually the source, i.e compressed file source etc>
	// should these all be filled in by the source???
	record.TpSourceType = c.Identifier()
	//record.TpSourceName = ???
	//record.TpSourceLocation = ???
	record.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))

	// Hive fields
	// TODO - should be based on the definition in HCL
	record.TpCollection = "default"
	if record.AccountID != nil {
		record.TpConnection = *record.AccountID
	}

	// populate the year, month, day from start time
	if record.Timestamp != nil {
		// TODO factor into function
		record.TpYear = int32(record.Timestamp.In(time.UTC).Year())
		record.TpMonth = int32(record.Timestamp.In(time.UTC).Month())
		record.TpDay = int32(record.Timestamp.In(time.UTC).Day())

		record.TpTimestamp = helpers.UnixMillis(record.Timestamp.UnixNano() / int64(time.Millisecond))

	} else if record.Start != nil {
		record.TpYear = int32(time.Unix(int64(*record.Start)/1000, 0).In(time.UTC).Year())
		record.TpMonth = int32(time.Unix(int64(*record.Start)/1000, 0).In(time.UTC).Month())
		record.TpDay = int32(time.Unix(int64(*record.Start)/1000, 0).In(time.UTC).Day())

		//convert from unis seconds to milliseconds
		record.TpTimestamp = helpers.UnixMillis(*record.Start * 1000)
	}

	//record.TpAkas = ???
	//record.TpTags = ???
	//record.TpDomains = ???
	//record.TpEmails = ???
	//record.TpUsernames = ???

	// ips
	if record.SrcAddr != nil {
		record.TpSourceIP = record.SrcAddr
		record.TpIps = append(record.TpIps, *record.SrcAddr)
	}
	if record.DstAddr != nil {
		record.TpDestinationIP = record.DstAddr
		record.TpIps = append(record.TpIps, *record.DstAddr)
	}

	return record, nil

}

func (c *FlowlogLogCollection) NewPagingData() (paging.Data, error) {
	// TODO use config to determine the type of paging data to return
	// hard coded to cloudwatch for now
	return paging.NewCloudwatch(), nil
}
