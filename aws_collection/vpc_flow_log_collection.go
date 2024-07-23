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

// VPCFlowLogLogCollection - collection for VPC Flow Logs
type VPCFlowLogLogCollection struct {
	// all collections must embed collection.Base
	collection.Base

	// the collection config
	Config *VpcFlowLogCollectionConfig
}

func NewVPCFlowLogLogCollection() plugin.Collection {
	l := &VPCFlowLogLogCollection{}
	return l
}

// Identifier implements plugin.Collection
func (c *VPCFlowLogLogCollection) Identifier() string {
	return "aws_vpc_flow_log"
}

// Init implements plugin.Collection
func (c *VPCFlowLogLogCollection) Init(ctx context.Context, configData []byte) error {
	// TEMP - this will actually parse (or the base will)
	// unmarshal the config
	config := &VpcFlowLogCollectionConfig{
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

// GetRowSchema implements plugin.Collection
// return an instance of the row struct
func (c *VPCFlowLogLogCollection) GetRowSchema() any {
	return aws_types.AwsVpcFlowLog{}
}

// GetPagingDataSchema implements plugin.Collection
func (c *VPCFlowLogLogCollection) GetPagingDataSchema() (paging.Data, error) {
	// TODO use config to determine the type of paging data to return
	// hard coded to cloudwatch for now
	return paging.NewCloudwatch(), nil
}

// GetConfigSchema implements plugin.Collection
func (c *VPCFlowLogLogCollection) GetConfigSchema() any {
	return VpcFlowLogCollectionConfig{}
}

// EnrichRow implements plugin.Collection
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

// use the config to configure the Source
func (c *VPCFlowLogLogCollection) getSource(ctx context.Context, config *VpcFlowLogCollectionConfig) (plugin.RowSource, error) {

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
		AccessKey:       "ASIARNKUQPUT75NBF3WU",
		SecretKey:       "R1cB+xJrGo12btp9K3abxtg4QuTr880zOWRAYRbD",
		SessionToken:    "IQoJb3JpZ2luX2VjELH//////////wEaCXVzLWVhc3QtMiJHMEUCIQDmzhz7qZpGdGBsCxxM9EX7aDOpZI5cYkyWYnsBKtRcbwIgKKXrltoeHUdC0fudFXeC29umNGngk8sClVGvBHwGj3UqjwMIiv//////////ARADGgwwOTczNTA4NzY0NTUiDLsl05hNSFLze2TTCyrjAhz+vNi8yxuOP+SMj+ZvFA8EfkVWUOR//HABPXAU6sk8mP99T4xkxok/r/ib7wuIBL9CFhTZrJESHKwHGnvNJZRWbzBNQ586Pq+hVBudrORW5upTgic2/CPuH1qXBqmDT8Jnkhol2nPb+waxsWwHV5pGL5Cn2aVKeCyknS1DNdTOFnkOxMW19SZteC2xEnSzLqAfPFQe+5FzaYvCu/29PtOUuZYbHM/h5RGar7kmlVDpTWKaaPXDNWHxMRxQR6hgxZ3nZEA7D44ok5tLiCwNME1dXMSYVjiEXRyCZEZqOfHATZH9/uQpLEVpdXtNU27YaugmjucI4dmeMyjaX+wGyd361kGfoY9lQceu5inRS86BB6uPFi8S1w3lfbdqW/ZPpIdgjJ0fP1I9VVmLuTdGSGrq0OdoColOPjOusBV7pmdTv/UvqRhLRtTj975MlF/kVnRrnTB6YIIR6g9VDjkJ/ARX/tsw6+f9tAY6pgE+JL4P4PKd0nEr8jnMFfDvb85TYOege8Dagrr2fU9VkGIW8a1HtqAU+MNBW9x9kl12raTtJBj9agne6BgPAseqCr/x1aV2s5QtUqtL0drWG4oh+RpU+DUgDrc93ZrKw6+H7KNRpBNGK6egZKsjmomWfT+ZLihjk0ZPe+dRHDv49KnjO8xgLfDCKFFAB2UOQLHsdhx4hYpPPcIHe/onYrI0PCilboia",
		LogGroupName:    "/victor/vpc/flowlog",
		LogStreamPrefix: utils.ToStringPointer("eni"),
		StartTime:       time.Now().Add(-time.Hour * 24),
		EndTime:         time.Now(),
	})

	// create empty paging data to pass to source
	// TODO maybe source creates for itself??
	pagingData, err := c.GetPagingDataSchema()
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
