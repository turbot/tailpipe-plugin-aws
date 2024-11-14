package tables

//
//// register the table from the package init function
//func init() {
//	table.RegisterTable(NewWafTrafficLogTable)
//}
//
//// WafTrafficLogTable - table for Waf traffic logs
//type WafTrafficLogTable struct {
//	// all tables must embed table.TableImpl
//	table.TableImpl[*rows.WafTrafficLog, *WafTrafficLogTableConfig, *config.AwsConnection]
//}
//
//func NewWafTrafficLogTable() table.Table[*rows.WafTrafficLog] {
//	return &WafTrafficLogTable{}
//}
//
//func (c *WafTrafficLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
//	// call base init
//	if err := c.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
//		return err
//	}
//
//	c.initMapper()
//	return nil
//}
//
//func (c *WafTrafficLogTable) initMapper() {
//	// TODO switch on source
//
//	// if the source is an artifact source, we need a mapper
//	c.Mapper = mappers.NewWafMapper()
//}
//
//// Identifier implements table.Table
//func (c *WafTrafficLogTable) Identifier() string {
//	return "aws_waf_traffic_log"
//}
//
//// GetSourceOptions returns any options which should be passed to the given source type
//func (c *WafTrafficLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
//	var opts []row_source.RowSourceOption
//
//	switch sourceType {
//	case artifact_source.AwsS3BucketSourceIdentifier:
//		// the default file layout for Waf traffic logs in S3
//		defaultArtifactConfig := &artifact_source_config.ArtifactSourceConfigBase{
//			// TODO #config finalise default Waf traffic file layout
//			FileLayout: utils.ToStringPointer("waf-access-log-sample_(?P<year>\\d{4})(?P<month>\\d{2})(?P<day>\\d{2})(?P<hour>\\d{2})(?P<minute>\\d{2})(?P<second>\\d{2}).gz"),
//		}
//		opts = append(opts, artifact_source.WithDefaultArtifactSourceConfig(defaultArtifactConfig), artifact_source.WithRowPerLine())
//	}
//
//	return opts
//}
//
//// GetRowSchema implements table.Table
//func (c *WafTrafficLogTable) GetRowSchema() types.RowStruct {
//	return rows.WafTrafficLog{}
//}
//
//func (c *WafTrafficLogTable) GetConfigSchema() parse.Config {
//	return &WafTrafficLogTableConfig{}
//}
//
//// EnrichRow implements table.Table
//func (c *WafTrafficLogTable) EnrichRow(row *rows.WafTrafficLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.WafTrafficLog, error) { // initialize the enrichment fields to any fields provided by the source
//	if sourceEnrichmentFields != nil {
//		row.CommonFields = *sourceEnrichmentFields
//	}
//
//	// Record standardization
//	row.TpID = xid.New().String()
//	row.TpSourceType = "aws_waf_traffic_log"
//	row.TpTimestamp = *row.Timestamp
//	row.TpIngestTimestamp = time.Now()
//
//	// Hive fields
//	// Check if row.HttpSourceId is not nil before dereferencing
//	if row.HttpSourceId != nil {
//		row.TpIndex = strings.ReplaceAll(*row.HttpSourceId, "/", `_`)
//	} else {
//		// Handle the case where HttpSourceId is nil, if needed
//		row.TpIndex = "" // or assign a default value if desired
//	}
//	// convert to date in format yy-mm-dd
//	row.TpDate = row.Timestamp.Truncate(24 * time.Hour)
//
//	return row, nil
//}
