package aws_collection

import (
	"context"
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/artifact"
	"github.com/turbot/tailpipe-plugin-sdk/paging"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_source"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-aws/util"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

type CloudTrailLogCollection struct {
	// all collections must embed collection.Base
	collection.Base

	// the collection config
	Config *CloudTrailLogCollectionConfig
}

func NewCloudTrailLogCollection() plugin.Collection {
	return &CloudTrailLogCollection{}
}

// Identifier implements Collection
func (c *CloudTrailLogCollection) Identifier() string {
	return "aws_cloudtrail_log"
}

// GetRowStruct implements Collection
// return an instance of the row struct
func (c *CloudTrailLogCollection) GetRowStruct() any {
	return aws_types.AWSCloudTrail{}
}

// Init implements Collection
func (c *CloudTrailLogCollection) Init(ctx context.Context, configData []byte) error {
	// TEMP - this will actually parse (or the base will)
	// unmarshal the config
	config := &CloudTrailLogCollectionConfig{
		Paths: []string{"/Users/kai/tailpipe_data/flaws_cloudtrail_logs"},
	}

	//err := json.Unmarshal(configData, config)
	//if err != nil {
	//	return fmt.Errorf("error unmarshalling config: %w", err)
	//}
	// todo - parse config as hcl
	c.Config = config
	// todo validate config

	// todo create source from config
	source, err := c.getSource(c.Config)
	if err != nil {
		return err
	}
	return c.AddSource(source)
}

func (c *CloudTrailLogCollection) getSource(config *CloudTrailLogCollectionConfig) (plugin.RowSource, error) {
	sourceConfig := &artifact.FileSystemSourceConfig{Paths: config.Paths, Extensions: []string{".gz"}}

	artifactSource := artifact.NewFileSystemSource(sourceConfig)
	artifactMapper := aws_source.NewCloudtrailMapper()

	// create empty paging data to pass to source
	// TODO maybe source creates for itself??
	pagingData, err := c.NewPagingData()
	if err != nil {
		return nil, fmt.Errorf("error creating paging data: %w", err)
	}

	source, err := row_source.NewArtifactRowSource(artifactSource, pagingData, row_source.WithMapper(artifactMapper))

	if err != nil {
		return nil, fmt.Errorf("error creating artifact row source: %w", err)
	}

	return source, nil
}

// EnrichRow implements RowEnricher
func (c *CloudTrailLogCollection) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// row must be an AWSCloudTrail
	record, ok := row.(aws_types.AWSCloudTrail)
	if !ok {
		return nil, fmt.Errorf("invalid row type %T, expected AWSCloudTrail", row)
	}

	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		record.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	record.TpID = xid.New().String()
	record.TpSourceType = "aws_cloudtrail_log"
	record.TpTimestamp = record.EventTime
	record.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))
	if record.SourceIPAddress != nil {
		record.TpSourceIP = record.SourceIPAddress
		record.TpIps = append(record.TpIps, *record.SourceIPAddress)
	}
	for _, resource := range record.Resources {
		if resource.ARN != nil {
			newAkas := util.AwsAkasFromArn(*resource.ARN)
			record.TpAkas = append(record.TpAkas, newAkas...)
		}
	}
	// If it's an AKIA, then record that as an identity. Do not record ASIA*
	// keys etc.
	if record.UserIdentity.AccessKeyId != nil {
		if strings.HasPrefix(*record.UserIdentity.AccessKeyId, "AKIA") {
			record.TpUsernames = append(record.TpUsernames, *record.UserIdentity.AccessKeyId)
		}
	}
	if record.UserIdentity.UserName != nil {
		record.TpUsernames = append(record.TpUsernames, *record.UserIdentity.UserName)
	}

	// Hive fields
	record.TpCollection = "default" // TODO - should be based on the definition in HCL
	record.TpConnection = record.RecipientAccountId
	record.TpYear = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Year())
	record.TpMonth = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Month())
	record.TpDay = int32(time.Unix(int64(record.EventTime)/1000, 0).In(time.UTC).Day())

	return record, nil
}

// TODO use config to determine correct paging data
// TODO move to source???
func (c *CloudTrailLogCollection) NewPagingData() (paging.Data, error) {
	// TODO use config to determine the type of paging data to return
	// hard coded to cloudwatch for now
	return paging.NewCloudwatch(), nil
}
