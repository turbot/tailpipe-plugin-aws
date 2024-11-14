package tables

//
//// register the table from the package init function
//func init() {
//	table.RegisterTable(NewVPCFlowLogLogTable)
//}
//
//// VPCFlowLogLogTable - table for VPC Flow Logs
//type VPCFlowLogLogTable struct {
//	// all tables must embed table.TableImpl
//	table.TableImpl[*rows.VpcFlowLog, *VpcFlowLogTableConfig, *config.AwsConnection]
//}
//
//func NewVPCFlowLogLogTable() table.Table[*rows.VpcFlowLog] {
//	return &VPCFlowLogLogTable{}
//}
//
//func (c *VPCFlowLogLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
//	// call base init
//	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
//		return err
//	}
//
//	c.initMapper()
//	return nil
//}
//
//func (c *VPCFlowLogLogTable) initMapper() {
//	// TODO switch on source
//	c.Mapper = mappers.NewVpcFlowlogMapper(c.Config.Fields)
//}
//
//// Identifier implements table.Table
//func (c *VPCFlowLogLogTable) Identifier() string {
//	return "aws_vpc_flow_log"
//}
//
//// GetRowSchema implements table.Table
//// return an instance of the row struct
//func (c *VPCFlowLogLogTable) GetRowSchema() types.RowStruct {
//	return rows.VpcFlowLog{}
//}
//
//func (c *VPCFlowLogLogTable) GetConfigSchema() parse.Config {
//	return &VpcFlowLogTableConfig{}
//}
//
//// EnrichRow implements table.Table
//func (c *VPCFlowLogLogTable) EnrichRow(row *rows.VpcFlowLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.VpcFlowLog, error) {
//	// initialize the enrichment fields to any fields provided by the source
//	if sourceEnrichmentFields != nil {
//		row.CommonFields = *sourceEnrichmentFields
//	}
//
//	// Record standardization
//	row.TpID = xid.New().String()
//
//	// TODO is source type actually the source, i.e compressed file source etc>
//	// should these all be filled in by the source???
//	row.TpSourceType = c.Identifier()
//	//row.TpSourceName = ???
//	//row.TpSourceLocation = ???
//	row.TpIngestTimestamp = time.Now()
//
//	// Hive fields
//	// TODO - should be based on the definition in HCL
//	row.TpPartition = "default"
//	if row.AccountID != nil {
//		row.TpIndex = *row.AccountID
//	}
//
//	// populate the year, month, day from start time
//	if row.Timestamp != nil {
//		// convert to date in format yy-mm-dd
//		row.TpDate = row.Timestamp.In(time.UTC).Truncate(24 * time.Hour)
//		row.TpTimestamp = *row.Timestamp
//	} else if row.Start != nil {
//		// convert to date in format yy-mm-dd
//		// TODO is Start unix millis?? if so why do we convert it for TpTimestamp
//		row.TpDate = time.UnixMilli(*row.Start).Truncate(24 * time.Hour)
//
//		//convert from unis seconds to milliseconds
//		row.TpTimestamp = time.Unix(0, int64(*row.Start*1000)*int64(time.Millisecond))
//	}
//
//	//row.TpAkas = ???
//	//row.TpTags = ???
//	//row.TpDomains = ???
//	//row.TpEmails = ???
//	//row.TpUsernames = ???
//
//	// ips
//	if row.SrcAddr != nil {
//		row.TpSourceIP = row.SrcAddr
//		row.TpIps = append(row.TpIps, *row.SrcAddr)
//	}
//	if row.DstAddr != nil {
//		row.TpDestinationIP = row.DstAddr
//		row.TpIps = append(row.TpIps, *row.DstAddr)
//	}
//
//	return row, nil
//}
