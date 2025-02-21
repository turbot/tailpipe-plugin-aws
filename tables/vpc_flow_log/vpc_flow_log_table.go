package vpc_flow_log

import (
	"log/slog"
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
type VpcFlowLogTable struct {}

func (c *VpcFlowLogTable) GetSourceMetadata() []*table.SourceMetadata[*VpcFlowLog] {

	defaultS3ArtifactConfig := &artifact_source_config.ArtifactSourceConfigImpl{
		// TODO: Handle the case when Partition logs by time selected to: 
		// "Every 24 hours(default)" - %{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}
		// "Every 1 hour(60 minutes)" - %{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:hour}
		FileLayout: utils.ToStringPointer("AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:hour}/%{DATA}.log.gz"),

		// s3://delete-me-35/AWSLogs/xxxxxxxxxxxx/vpcflowlogs/us-east-1/2025/02/12/10/xxxxxxxxxxxx_vpcflowlogs_us-east-1_fl-05f57cebeaa68cbdb_20250212T1050Z_884a1fb8.log.gz

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

	// TODO: How to handle if the log don't have timestamp value 
	// populate the year, month, day from start time
	if row.Timestamp != nil {
		// convert to date in format yy-mm-dd
		slog.Error("Timestamp taken form the field ===>>", "rows.Timestamp", "Ok")
		row.TpDate = row.Timestamp.In(time.UTC).Truncate(24 * time.Hour)
		row.TpTimestamp = *row.Timestamp
	} else if row.Start != nil {
		// convert to date in format yy-mm-dd
		// TODO is Start unix millis?? if so why do we convert it for TpTimestamp
		slog.Error("Timestamp taken form the field ===>>", "rows.Start", "Ok")
		row.TpDate = time.UnixMilli(*row.Start).Truncate(24 * time.Hour)

		//convert from unis seconds to milliseconds
		row.TpTimestamp = time.Unix(0, int64(*row.Start*1000)*int64(time.Millisecond))
	} else {
		// TODO: is it correct to fallback to the current time
		slog.Error("Timestamp taken form the field ===>>", "rows.NOT-AVAILABLE", "Ok")
		t := time.Now()
		row.TpDate = time.UnixMilli(t.UnixMilli()).Truncate(24 * time.Hour)
		row.TpTimestamp = t
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