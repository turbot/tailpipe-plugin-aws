package tables

// TODO this table need updating to use a RowStruct which embeds CommonFields
//
//
//type LambdaLogTable struct {
//	// all tables must embed table.TableImpl
//	table.TableImpl[string, *LambdaLogTableConfig, *config.AwsConnection]
//}
//
//func NewLambdaLogTable() table.Table {
//	return &LambdaLogTable{}
//}
//
//func (c *LambdaLogTable) Identifier() string {
//	return "aws_lambda_log"
//}
//
//func (c *LambdaLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
//	return []row_source.RowSourceOption{
//		artifact_source.WithRowPerLine(),
//	}
//}
//
//func (c *LambdaLogTable) GetRowSchema() types.Validatable {
//	return &rows.AwsLambdaLog{}
//}
//
//func (c *LambdaLogTable) GetConfigSchema() parse.Config {
//	return &LambdaLogTableConfig{}
//}
//
//// TODO K why does this not fail to compile?
//func (c *LambdaLogTable) EnrichRow(rawRow string, sourceEnrichmentFields *enrichment.CommonFields) (any, error) {
//	var row rows.AwsLambdaLog
//	if sourceEnrichmentFields != nil {
//		row.CommonFields = *sourceEnrichmentFields
//
//		ts := sourceEnrichmentFields.TpTimestamp
//		row.Timestamp = &ts
//	}
//
//	// remove trailing newline
//	rawRow = strings.TrimSuffix(rawRow, "\n")
//	fields := strings.Fields(rawRow)
//
//	switch fields[0] {
//	case "START", "END":
//		row.LogType = &fields[0]
//		row.RequestID = &fields[2]
//	case "REPORT":
//		row.LogType = &fields[0]
//		row.RequestID = &fields[2]
//		duration, err := strconv.ParseFloat(fields[4], 64)
//		if err != nil {
//			return nil, fmt.Errorf("error parsing duration: %w", err)
//		}
//		row.Duration = &duration
//		billed, err := strconv.ParseFloat(fields[8], 64)
//		if err != nil {
//			return nil, fmt.Errorf("error parsing billed duration: %w", err)
//		}
//		row.BilledDuration = &billed
//		mem, err := strconv.Atoi(fields[12])
//		if err != nil {
//			return nil, fmt.Errorf("error parsing memory size: %w", err)
//		}
//		row.MemorySize = &mem
//		maxMem, err := strconv.Atoi(fields[17])
//		if err != nil {
//			return nil, fmt.Errorf("error parsing max memory used: %w", err)
//		}
//		row.MaxMemoryUsed = &maxMem
//	default:
//		t := "LOG"
//		row.LogType = &t
//		// TODO: #enrich should we overwrite the timestamp with that in the log entry?
//		row.RequestID = &fields[1]
//		row.LogLevel = &fields[2]
//		strip := fmt.Sprintf("%s%s", strings.Join(fields[:3], "\t"), "\t")
//		stripped := strings.TrimPrefix(rawRow, strip)
//		row.Message = &stripped
//	}
//
//	// Record standardization
//	row.TpID = xid.New().String()
//	row.TpIngestTimestamp = time.Now()
//
//	// Hive fields
//	row.TpPartition = c.Identifier()
//	if row.TpIndex == "" {
//		row.TpIndex = c.Identifier() // TODO: #refactor figure out how to get connection (account ID?)
//	}
//	// convert to date in format yy-mm-dd
//	row.TpDate = row.Timestamp.Format("2006-01-02")
//
//	return row, nil
//}
