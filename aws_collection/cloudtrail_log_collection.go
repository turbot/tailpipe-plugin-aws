package aws_collection

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-aws/util"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"strings"
	"time"
	//"github.com/turbot/tailpipe-plugin-sdk/collection"
	//sdkconfig "github.com/turbot/tailpipe-plugin-sdk/config"
	//"github.com/turbot/tailpipe-plugin-sdk/source"
)

type CloudTrailLogCollection struct {
	// all collections must embed collection.Base
	// this add observer and enrich functions
	collection.Base

	// the collection config
	Config CloudTrailLogCollectionConfig
}

func NewCloudTrailLogCollection(config CloudTrailLogCollectionConfig, source plugin.Source) *CloudTrailLogCollection {
	l := &CloudTrailLogCollection{
		Config: config,
	}
	// Init sets the Source property on Base and adds us as an observer to it
	// It also sets us as the Enricher property on Base
	l.Base.Init(source, l)

	return l
}

func (c CloudTrailLogCollection) Identifier() string {
	return "aws_cloudtrail_log"
}

func (c CloudTrailLogCollection) EnrichRow(row any, sourceEnrichmentFields map[string]any) (any, error) {
	// row must be an AWSCloudTrail
	record, ok := row.(aws_types.AWSCloudTrail)
	if !ok {
		return nil, fmt.Errorf("invalid row type %T, expected AWSCloudTrail", row)
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

	// now add an enrichment fields provided by the source
	if sourceEnrichmentFields != nil {
		if sourceLocation, ok := sourceEnrichmentFields["sourceLocation"]; ok {
			// verify it is a string
			sourceStr, ok := sourceLocation.(string)
			if ok {
				return nil, fmt.Errorf("sourceLocation is not a string")
			}
			record.TpSourceLocation = &sourceStr
		}
	}
	return row, nil

}
