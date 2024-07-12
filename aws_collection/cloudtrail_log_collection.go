package aws_collection

import (
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/artifact"
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
	// this add observer and enrich functions
	collection.Base

	// the collection config
	Config *CloudTrailLogCollectionConfig
}

func NewCloudTrailLogCollection() plugin.Collection {
	c := &CloudTrailLogCollection{}
	// TODO avoid need for plugin implementation to do this
	// Init sets us as the Enricher property on Base
	c.Base.Init(c)

	return c
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
func (c *CloudTrailLogCollection) Init(config any) error {
	// TEMP - this will actually parse (or the base will)

	// todo - parse config
	c.Config = config.(*CloudTrailLogCollectionConfig)
	// todo validate config

	// todo create source from config
	source, err := c.getSource(c.Config)
	if err != nil {
		return err
	}
	c.AddSource(source)

	return nil
}

func (c *CloudTrailLogCollection) getSource(config *CloudTrailLogCollectionConfig) (plugin.RowSource, error) {
	sourceConfig := &artifact.FileSystemSourceConfig{Paths: config.Paths, Extensions: []string{".gz"}}

	artifactSource := artifact.NewFileSystemSource(sourceConfig)
	artifactLoader := artifact.NewGzipObjectLoader[*aws_types.AWSCloudTrailBatch]()
	artifactMapper := aws_source.NewCloudtrailMapper()

	var source, err = row_source.NewArtifactRowSource(
		artifactSource,
		artifactLoader,
		artifactMapper,
	)

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
	if sourceEnrichmentFields == nil {
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
