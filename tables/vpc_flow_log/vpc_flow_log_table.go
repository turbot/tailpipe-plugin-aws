package vpc_flow_log

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source_config"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const VpcFlowLogTableIdentifier = "aws_vpc_flow_log"

func init() {
	// 1. row struct
	// 2. input format struct
	// 3. table implementation
	table.RegisterTableFormat[*VpcFlowLog, *VpcFlowLogTableFormat, *VpcFlowLogTable]()
}

// VpcFlowLogTable - table for VPC Flow Logs
type VpcFlowLogTable struct {
	table.TableWithFormatImpl[*VpcFlowLogTableFormat]
}

func (c *VpcFlowLogTable) GetSourceMetadata() []*table.SourceMetadata[*VpcFlowLog] {
	fields := DefaultFlowLogFields
	// if Format was provided in config, it will have been populated by the factory
	if c.Format != nil {
		fields = c.Format.Fields
	}

	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		// FileLayout: utils.ToStringPointer("AWSLogs/632902152528/vpcflowlogs/us-east-1/2025/02/12/10/%{DATA}.log.gz"),
		// FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:hour}/%{DATA}.log.gz"),

		// s3://delete-me-35/AWSLogs/xxxxxxxxxxxx/vpcflowlogs/us-east-1/2025/02/12/10/xxxxxxxxxxxx_vpcflowlogs_us-east-1_fl-05f57cebeaa68cbdb_20250212T1050Z_884a1fb8.log.gz

	}

	return []*table.SourceMetadata[*VpcFlowLog]{
		{
			// S3 artifact source
			SourceName: s3_bucket.AwsS3BucketSourceIdentifier,
			Mapper:     NewVpcFlowLogMapper(fields),
			Options: []row_source.RowSourceOption{
				artifact_source.WithDefaultArtifactSourceConfig(defaultS3ArtifactConfig),
				artifact_source.WithRowPerLine(),
			},
		},
		// {
		// 	// any artifact source
		// 	SourceName: constants.ArtifactSourceIdentifier,
		// 	Mapper:     mappers.NewVpcFlowLogMapper(fields),
		// 	Options: []row_source.RowSourceOption{
		// 		artifact_source.WithRowPerLine(),
		// 	},
		// },
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

	// TODO is source type actually the source, i.e compressed file source etc>
	// should these all be filled in by the source???
	row.TpSourceType = c.Identifier()
	//row.TpSourceName = ???
	//row.TpSourceLocation = ???
	row.TpIngestTimestamp = time.Now()

	// Hive fields
	if row.AccountID != nil {
		row.TpIndex = *row.AccountID
	} else {
		row.TpIndex = "default"
	}

	// populate the year, month, day from start time
	if row.Timestamp != nil {
		// convert to date in format yy-mm-dd
		row.TpDate = row.Timestamp.In(time.UTC).Truncate(24 * time.Hour)
		row.TpTimestamp = *row.Timestamp
	} else if row.Start != nil {
		// convert to date in format yy-mm-dd
		// TODO is Start unix millis?? if so why do we convert it for TpTimestamp
		row.TpDate = time.UnixMilli(*row.Start).Truncate(24 * time.Hour)

		//convert from unis seconds to milliseconds
		row.TpTimestamp = time.Unix(0, int64(*row.Start*1000)*int64(time.Millisecond))
	}

	//row.TpAkas = ???
	//row.TpTags = ???
	//row.TpDomains = ???
	//row.TpEmails = ???
	//row.TpUsernames = ???

	// ips
	if row.SrcAddr != nil {
		row.TpSourceIP = row.SrcAddr
		row.TpIps = append(row.TpIps, *row.SrcAddr)
	}
	if row.DstAddr != nil {
		row.TpDestinationIP = row.DstAddr
		row.TpIps = append(row.TpIps, *row.DstAddr)
	}

	return row, nil
}
