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
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz"),
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

	// TpIndex
	if row.InterfaceID != nil {
		row.TpIndex = *row.InterfaceID
	} else if row.SubnetID != nil {
		row.TpIndex = *row.SubnetID
	} else if row.VPCID != nil {
		row.TpIndex = *row.VPCID
	} else {
		row.TpIndex = "default"
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

	// If there is no start/end time, there is no usable log timestamp. In this case, the line will result in an error during collection.
	if row.Start != nil {
		// convert to date in format yyyy-mm-dd
		row.TpDate = row.Start.Truncate(24 * time.Hour)
		row.TpTimestamp = *row.Start
	} else if row.End != nil {
		// convert to date in format yyyy-mm-dd
		row.TpDate = row.End.Truncate(24 * time.Hour)
		row.TpTimestamp = *row.End
	}

	return row, nil
}
