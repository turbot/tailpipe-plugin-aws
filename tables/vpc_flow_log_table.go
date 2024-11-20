package tables

import (
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/mappers"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
	"github.com/turbot/tailpipe-plugin-sdk/constants"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

const VpcFlowLogTableIdentifier = "aws_vpc_flow_log"

// register the table from the package init function
func init() {
	table.RegisterTable[*rows.VpcFlowLog, *VpcFlowLogTable]()
}

// VpcFlowLogTable - table for VPC Flow Logs
type VpcFlowLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[*rows.VpcFlowLog, *VpcFlowLogTableConfig, *artifact_source.AwsConnection]
}

func (c *VpcFlowLogTable) initMapper() func() table.Mapper[*rows.VpcFlowLog] {
	f := func() table.Mapper[*rows.VpcFlowLog] {
		return mappers.NewVpcFlowLogMapper(c.Config.Fields)
	}
	return f
}

func (c *VpcFlowLogTable) SupportedSources() []*table.SourceMetadata[*rows.VpcFlowLog] {
	return []*table.SourceMetadata[*rows.VpcFlowLog]{
		{
			// any artifact source
			SourceName: constants.ArtifactSourceIdentifier,
			MapperFunc: c.initMapper(),
			Options: []row_source.RowSourceOption{
				artifact_source.WithRowPerLine(),
			},
		},
	}
}

// Identifier implements table.Table
func (c *VpcFlowLogTable) Identifier() string {
	return VpcFlowLogTableIdentifier
}

// EnrichRow implements table.Table
func (c *VpcFlowLogTable) EnrichRow(row *rows.VpcFlowLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.VpcFlowLog, error) {
	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		row.CommonFields = *sourceEnrichmentFields
	}

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
