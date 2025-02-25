package vpc_flow_log

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const VpcFlowLogTableIdentifier = "aws_vpc_flow_log"

// VpcFlowLogTable - table for VPC Flow Logs
type VpcFlowLogTable struct{}

func (c *VpcFlowLogTable) GetSourceMetadata() []*table.SourceMetadata[*VpcFlowLog] {

	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		// Handle the case when Partition logs by time selected to:
		// "Every 24 hours(default)" - %{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}
		// "Every 1 hour(60 minutes)" - %{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:hour}
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz"),

		// The exported logs from AWS CloudWatch Log Group will follow a path format like:
		// `s3://test-vpc-flow-log-from-log-group/exportedlogs/453054b5-4448-4504-b2b4-15b6aa675a59/eni-0416a1c81c87ab9c9-all/000000.gz`.
		// Additionally, the logs do not include a header value, and it's unclear how to handle this scenario effectively.
	}

	return []*table.SourceMetadata[*VpcFlowLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithArtifactExtractor(NewVPCFlowLogExtractor()),
			},
		},
	}
}

// Identifier implements table.Table
func (c *VpcFlowLogTable) Identifier() string {
	return VpcFlowLogTableIdentifier
}

// EnrichRow implements table.Table
func (c *VpcFlowLogTable) EnrichRow(row *VpcFlowLog, sourceEnrichmentFields schema.SourceEnrichment) (*VpcFlowLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	row.CommonFields = sourceEnrichmentFields.CommonFields

	// Record standardization
	row.TpID = xid.New().String()
	row.TpIngestTimestamp = time.Now()

	// ips
	if row.SrcAddr != nil {
		row.TpSourceIP = row.SrcAddr
		row.TpIps = append(row.TpIps, *row.SrcAddr)
	}
	if row.PktSrcAddr != nil {
		row.TpIps = append(row.TpIps, *row.PktSrcAddr)
	}
	if row.DstAddr != nil {
		row.TpDestinationIP = row.DstAddr
		row.TpIps = append(row.TpIps, *row.DstAddr)
	}
	if row.PktDstAddr != nil {
		row.TpIps = append(row.TpIps, *row.PktDstAddr)
	}

	// TODO: Is it correct?
	if row.AccountID != nil {
		row.TpIndex = *row.AccountID
	} else {
		row.TpIndex = "default"
	}

	// TpSourceLocation
	if row.SublocationID != nil {
		row.TpSourceLocation = row.SublocationID
	}

	// TpAkas
	if row.ECSClusterARN != nil {
		row.TpAkas = append(row.TpAkas, *row.ECSClusterARN)
	}
	if row.ECSContainerInstanceARN != nil {
		row.TpAkas = append(row.TpAkas, *row.ECSContainerInstanceARN)
	}
	if row.ECSTaskARN != nil {
		row.TpAkas = append(row.TpAkas, *row.ECSTaskARN)
	}
	if row.ECSTaskDefinitionARN != nil {
		row.TpAkas = append(row.TpAkas, *row.ECSTaskDefinitionARN)
	}

	// TODO: How to handle if the log don't have timestamp value
	// populate the year, month, day from start time
	if row.Timestamp != nil {
		// convert to date in format yy-mm-dd
		row.TpDate = row.Timestamp.In(time.UTC).Truncate(24 * time.Hour)
		row.TpTimestamp = *row.Timestamp
	} else if row.Start != nil {
		// convert to date in format yy-mm-dd
		// TODO is Start unix millis?? if so why do we convert it for TpTimestamp
		row.TpDate = row.Start.Truncate(24 * time.Hour)

		//convert from unis seconds to milliseconds
		row.TpTimestamp = *row.Start
	} else {
		// TODO: is it correct to fallback to the current time
		t := time.Now()
		row.TpDate = time.UnixMilli(t.UnixMilli()).Truncate(24 * time.Hour)
		row.TpTimestamp = t
	}

	return row, nil
}
