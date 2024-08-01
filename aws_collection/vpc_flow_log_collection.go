package aws_collection

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
)

// VPCFlowLogLogCollection - collection for VPC Flow Logs
type VPCFlowLogLogCollection struct {
	// all collections must embed collection.CollectionBase
	collection.CollectionBase[VpcFlowLogCollectionConfig]
}

func NewVPCFlowLogLogCollection() collection.Collection {
	return &VPCFlowLogLogCollection{}
}

// Identifier implements collection.Collection
func (c *VPCFlowLogLogCollection) Identifier() string {
	return "aws_vpc_flow_log"
}

//
//// Init implements collection.Collection
//func (c *VPCFlowLogLogCollection) Init(ctx context.Context, collectionConfig, sourceConfig *hcl.Data) error {
//	}
//	//// TEMP - this will actually parse (or the base will)
//	//// unmarshal the config
//	//config := &VpcFlowLogCollectionConfig{
//	//	Paths: []string{"/Users/kai/tailpipe_data/flowlog"},
//	//	// use defaults when using cloudwatch
//	//	//		//Fields: []string{"timestamp",
//	//	//		//	"version",
//	//	//		//	"account-id",
//	//	//		//	"interface-id",
//	//	//		//	"srcaddr",
//	//	//		//	"dstaddr",
//	//	//		//	"srcport",
//	//	//		//	"dstport",
//	//	//		//	"protocol",
//	//	//		//	"packets",
//	//	//		//	"bytes",
//	//	//		//	"start",
//	//	//		//	"end",
//	//	//		//	"action",
//	//	//		//	"log-status",
//	//	//		//},
//	//}
//	//
//
//	// init the config = this will set default fields if needed
//	if err := c.Config.Init(); err != nil {
//		return fmt.Errorf("error initializing config: %w", err)
//	}
//	// todo validate config
//
//}

// GetRowSchema implements collection.Collection
// return an instance of the row struct
func (c *VPCFlowLogLogCollection) GetRowSchema() any {
	return aws_types.AwsVpcFlowLog{}
}

// EnrichRow implements collection.Collection
func (c *VPCFlowLogLogCollection) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
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

//
//// use the config to configure the ArtifactSource
//func (c *VPCFlowLogLogCollection) getSource(ctx context.Context, config *VpcFlowLogCollectionConfig) (plugin.RowSource, error) {
//	// TODO populate from config
//	//srcConfig := &artifact.AwsS3BucketSourceConfig{
//	//	SourceConfigBase: artifact.SourceConfigBase{
//	//		// TODO #config where do we configure this
//	//		TmpDir: os.TempDir(),
//	//	},
//	//	Bucket:       "silverwater-flowlog-s3-bucket",
//	//	Extensions:   []string{".gz"},
//	//	AccessKey:    "",
//	//	SecretKey:    "",
//	//	SessionToken: "",
//	//}
//	//
//	//artifactSource, err := artifact.NewAwsS3BucketSource(sourceConfig)
//	//if err != nil {
//	//	return nil, fmt.Errorf("error creating s3 bucket source: %w", err)
//	//}
//
//	//artifactSource := artifact.NewFileSystemSource(&artifact.FileSystemSourceConfig{
//
//	//	Paths:      config.Paths,
//	//	Extensions: []string{".gz"},
//	//})
//	srcConfig := &artifact.AwsCloudWatchSourceConfig{
//		SourceConfigBase: artifact.SourceConfigBase{
//			// TODO #config where do we configure this
//			TmpDir: os.TempDir(),
//		},
//		AccessKey:       "ASIARNKUQPUT2DKH2CJH",
//		SecretKey:       "o3W4OjtkAs418H4h/rJGcorIHQVWw8qVGCFXTW0Y",
//		SessionToken:    "IQoJb3JpZ2luX2VjEOf//////////wEaCXVzLWVhc3QtMiJIMEYCIQCMQtQOXk2UvNDtrKpGVv/KpXph86EJ8LVnGK+RZM5epQIhAO3oruAIomKUK8ljVNcQrBL7CtkkHJsSboBX9FrF/pwUKo8DCMD//////////wEQAxoMMDk3MzUwODc2NDU1Igxp82/qL0ja9Ki1M90q4wLDWC+4stlMrQYQV2hGIZXekIRClllYiwyzModW5cx1lVDnp1s1Twnm2eDcA707NzMX5v5k6p+A2svLx6JUg2ymGOBHB067IyanCXNCzGK/+FG+Ec2L419r2Jt+sxspvNthcQhkVFAPIkdFBCr0cdcIc60IrVF5TCHhsf94bQHRE/NXGwwZeeUZiXCaz6QQDsr56+p/bGHR1QTxLfwZWSNs/tI1IuTRaiBk0LOfACgWAYWq5VWwDQbqoYCIBOU+rpH/YNNmsUgQQbIS9VKbu2vuAoUeXuOiST2J5YRriX+potowmnICtkEFuIAfpuYz6uA4xwT3fOtQNjs3dEfKjus7ZrFv0oO3PtUHWGC3Nsn0y2/jkwYlasTekfKNd0Wz/cnKwEjtNWWFw1+MXNbocGc9lBulmN4StMOTUI/jQdlm3tXJ861CD+4Qm23aulw9IFZODdL7dHcHx1gIPAsNt1qeXEgOMIvLibUGOqUBBoIM4Qt+VO0jO/1Li4WEfPB/Pd3+D12VaPwC7503YtNEaZOg9SmWeXht/G9Fq6B+bUsNrA0ADAFLotr4QZZsUST2BJU5lf3vAIbi9tH4WUMzY4kK7huqax3TICN7kuRO/kxbWvXyeqTPOayi47xgEOc5YiRhyuyjZfwxpZdJVxvOG+AOoqbhphRe9DvbVhUylrKDj0FpH93zuDOYUl0a7PCBM1ns",
//		LogGroupName:    "/victor/vpc/flowlog",
//		LogStreamPrefix: utils.ToStringPointer("eni"),
//		StartTime:       time.Now().Add(-time.Hour * 24),
//		EndTime:         time.Now(),
//	}
//
//	artifactSource, err := artifact.NewAwsCloudWatchSource(ctx, srcConfig)
//
//	// create empty paging data to pass to source
//	// TODO maybe source creates for itself??
//	pagingData, err := c.GetPagingDataSchema()
//	if err != nil {
//		return nil, fmt.Errorf("error creating paging data: %w", err)
//	}
//
//	source, err := row_source.NewArtifactRowSource(
//		artifactSource,
//		pagingData,
//		// we expect a log row per line of log data
//		row_source.WithRowPerLine(),
//	)
//	if err != nil {
//		return nil, fmt.Errorf("error creating artifact row source: %w", err)
//	}
//
//	return source, nil
//}
