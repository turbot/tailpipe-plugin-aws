package aws_table

import (
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/turbot/tailpipe-plugin-aws/aws_types"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"github.com/turbot/tailpipe-plugin-sdk/helpers"
	"github.com/turbot/tailpipe-plugin-sdk/parse"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// VPCFlowLogLogTable - table for VPC Flow Logs
type VPCFlowLogLogTable struct {
	// all tables must embed table.TableBase
	table.TableBase[*VpcFlowLogTableConfig]
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
	return aws_types.AwsVpcFlowLog{}
}

func (c *VPCFlowLogLogTable) GetConfigSchema() parse.Config {
	return &VpcFlowLogTableConfig{}
}

// EnrichRow implements table.Table
func (c *VPCFlowLogLogTable) EnrichRow(row any, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
	// row must be a string
	rowString, ok := row.(string)
	if !ok {
		return nil, fmt.Errorf("invalid row type %T, expected string", row)
	}
	record, err := aws_types.FlowLogFromString(rowString, c.Config.Fields)

	if err != nil {
		return nil, fmt.Errorf("error parsing row: %s", err)
	}

	// initialize the enrichment fields to any fields provided by the source
	if sourceEnrichmentFields != nil {
		record.CommonFields = *sourceEnrichmentFields
	}

	// Record standardization
	record.TpID = xid.New().String()

	// TODO is source type actually the source, i.e compressed file source etc>
	// should these all be filled in by the source???
	record.TpSourceType = c.Identifier()
	//record.TpSourceName = ???
	//record.TpSourceLocation = ???
	record.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))

	// Hive fields
	// TODO - should be based on the definition in HCL
	record.TpPartition = "default"
	if record.AccountID != nil {
		record.TpIndex = *record.AccountID
	}

	// populate the year, month, day from start time
	if record.Timestamp != nil {
		// convert to date in format yy-mm-dd
		record.TpDate = record.Timestamp.In(time.UTC).Format("2006-01-02")
		record.TpTimestamp = helpers.UnixMillis(record.Timestamp.UnixNano() / int64(time.Millisecond))

	} else if record.Start != nil {
		// convert to date in format yy-mm-dd
		// TODO is Start unix millis?? if so why do we convert it for TpTimestamp
		record.TpDate = time.UnixMilli(*record.Start).Format("2006-01-02")

		//convert from unis seconds to milliseconds
		record.TpTimestamp = helpers.UnixMillis(*record.Start * 1000)
	}

	//record.TpAkas = ???
	//record.TpTags = ???
	//record.TpDomains = ???
	//record.TpEmails = ???
	//record.TpUsernames = ???

	// ips
	if record.SrcAddr != nil {
		record.TpSourceIP = record.SrcAddr
		record.TpIps = append(record.TpIps, *record.SrcAddr)
	}
	if record.DstAddr != nil {
		record.TpDestinationIP = record.DstAddr
		record.TpIps = append(record.TpIps, *record.DstAddr)
	}

	return record, nil

}
