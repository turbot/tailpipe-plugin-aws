package aws_collection

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/artifact"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"time"
	//"github.com/turbot/tailpipe-plugin-sdk/collection"
	//sdkconfig "github.com/turbot/tailpipe-plugin-sdk/config"
	//"github.com/turbot/tailpipe-plugin-sdk/source"
)

type FlowlogLogCollection struct {
	// all collections must embed collection.Base
	// this add observer and enrich functions
	collection.Base

	// the collection config
	Config *FlowLogCollectionConfig
}

func NewFlowlogLogCollection() plugin.Collection {
	l := &FlowlogLogCollection{}
	// TODO avoid need for plugin implementation to do this
	// Init sets us as the Enricher property on Base
	l.Base.Init(l)

	return l
}

// GetRowStruct implements Collection
// return an instance of the row struct
func (c *FlowlogLogCollection) GetRowStruct() any {
	return aws_types.FlowLog{}
}

// Init implements Collection
func (c *FlowlogLogCollection) Init(config any) error {
	// TEMP - this will actually parse (or the base will)

	// todo - parse config
	c.Config = config.(*FlowLogCollectionConfig)
	// todo validate config

	// todo create source from config
	source, err := c.getSource(c.Config)
	if err != nil {
		return err
	}
	return c.AddSource(source)
}

func (c *FlowlogLogCollection) getSource(config *FlowLogCollectionConfig) (plugin.RowSource, error) {

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

	artifactSource := artifact.NewFileSystemSource(&artifact.FileSystemSourceConfig{
		Paths:      config.Paths,
		Extensions: []string{".gz"},
	})

	//artifactSource, err := artifact.NewAwsCloudWatchSource(&artifact.AwsCloudWatchSourceConfig{
	//	AccessKey:    "",
	//	SecretKey:    "",
	//	SessionToken: "",
	//	LogGroupName: "/victor/vpc/flowlog",
	//}, )

	source, err := row_source.NewArtifactRowSource(
		artifactSource,
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
