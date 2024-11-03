package tables

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// VPCFlowLogLogTable - table for VPC Flow Logs
type VPCFlowLogLogTable struct {
	// all tables must embed table.TableImpl
	table.TableImpl[string, *VpcFlowLogTableConfig, *config.AwsConnection]
}

func NewVPCFlowLogLogTable() table.Table {
	return &VPCFlowLogLogTable{}
}

// Identifier implements table.Table
func (c *VPCFlowLogLogTable) Identifier() string {
	return "aws_vpc_flow_log"
}

// GetRowSchema implements table.Table
// return an instance of the row struct
func (c *VPCFlowLogLogTable) GetRowSchema() any {
	return rows.AwsVpcFlowLog{}
}

func (c *VPCFlowLogLogTable) GetConfigSchema() parse.Config {
	return &VpcFlowLogTableConfig{}
}

// EnrichRow implements table.Table
func (c *VPCFlowLogLogTable) EnrichRow(rawRow string, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	row, err := rows.FlowLogFromString(rawRow, c.Config.Fields)

	if err != nil {
		return nil, fmt.Errorf("error parsing row: %s", err)
	}

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
	row.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))

	// Hive fields
	// TODO - should be based on the definition in HCL
	row.TpPartition = "default"
	if row.AccountID != nil {
		row.TpIndex = *row.AccountID
	}

	// populate the year, month, day from start time
	if row.Timestamp != nil {
		// convert to date in format yy-mm-dd
		row.TpDate = row.Timestamp.In(time.UTC).Format("2006-01-02")
		row.TpTimestamp = helpers.UnixMillis(row.Timestamp.UnixNano() / int64(time.Millisecond))

	} else if row.Start != nil {
		// convert to date in format yy-mm-dd
		// TODO is Start unix millis?? if so why do we convert it for TpTimestamp
		row.TpDate = time.UnixMilli(*row.Start).Format("2006-01-02")

		//convert from unis seconds to milliseconds
		row.TpTimestamp = helpers.UnixMillis(*row.Start * 1000)
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
