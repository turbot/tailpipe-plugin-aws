package aws_collection

import (
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-aws/util"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_mapper"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
)

// CloudTrailLogCollection - collection for CloudTrail logs
type CloudTrailLogCollection struct {
	// all collections must embed collection.CollectionBase
	collection.CollectionBase[CloudTrailLogCollectionConfig]

	// the collection config
	Config *CloudTrailLogCollectionConfig
}

func NewCloudTrailLogCollection() collection.Collection {
	return &CloudTrailLogCollection{}
}

func (c *CloudTrailLogCollection) SupportedSources() []string {
	// TODO #source do we need to to specify the type  or artifact source supported?
	return []string{
		// TODO #factory provider a shortcut for all artifact sources
		artifact_source.AwsS3BucketSourceIdentifier,
		artifact_source.FileSystemSourceIdentifier,
		artifact_source.GcpStorageBucketSourceIdentifier,
		artifact_source.AWSCloudwatchSourceIdentifier,
	}
}

// Identifier implements plugin.Collection
func (c *CloudTrailLogCollection) Identifier() string {
	return "aws_cloudtrail_log"
}

// GetSourceOptions returns any options which should be passed to the given source type
func (c *CloudTrailLogCollection) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
	switch sourceType {
	// TODO #factory provider a shortcut for all artifact sources
	// if source is an artifact source, use the cloudtrail mapper
	case artifact_source.AwsS3BucketSourceIdentifier,
		artifact_source.FileSystemSourceIdentifier,
		artifact_source.GcpStorageBucketSourceIdentifier,
		artifact_source.AWSCloudwatchSourceIdentifier:
		return []row_source.RowSourceOption{artifact_source.WithMapper(artifact_mapper.NewCloudwatchMapper())}
	}
	return nil
}

// GetRowSchema implements plugin.Collection
func (c *CloudTrailLogCollection) GetRowSchema() any {
	return aws_types.AWSCloudTrail{}
}

// EnrichRow implements plugin.Collection
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

//
//// use the config to configure the ArtifactSource
//func (c *CloudTrailLogCollection) getSource(configData *hcl.Data) (plugin.RowSource, error) {
//switch configData.Type {
//
//}
//	var cfg CloudTrailLogCollectionConfig
//	err := hcl.ParseConfig(configData, &cfg)
//	if err != nil {
//		return nil, fmt.Errorf("error parsing config: %w", err)
//	}
//
//
//	sourceConfig := &artifact.FileSystemSourceConfig{Paths: hcl.Paths, Extensions: []string{".gz"}}
//
//	artifactSource := artifact.NewFileSystemSource(sourceConfig)
//	artifactMapper := aws_source.NewCloudtrailMapper()
//
//	// create empty paging data to pass to source
//	// TODO maybe source creates for itself??
//	pagingData, err := c.GetPagingDataSchema()
//	if err != nil {
//		return nil, fmt.Errorf("error creating paging data: %w", err)
//	}
//
//	source, err := row_source.NewArtifactRowSource(artifactSource, pagingData, )
//
//	if err != nil {
//		return nil, fmt.Errorf("error creating artifact row source: %w", err)
//	}
//
//	return source, nil
//}
