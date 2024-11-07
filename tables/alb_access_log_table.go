package tables

import (
    "context"
    "fmt"
    "time"
    "strings"
    "strconv"

    "github.com/rs/xid"
    "github.com/turbot/tailpipe-plugin-aws/config"
    "github.com/turbot/tailpipe-plugin-aws/rows"
    "github.com/turbot/tailpipe-plugin-sdk/artifact_source"
    "github.com/turbot/tailpipe-plugin-sdk/enrichment"
    "github.com/turbot/tailpipe-plugin-sdk/helpers"
    "github.com/turbot/tailpipe-plugin-sdk/parse"
    "github.com/turbot/tailpipe-plugin-sdk/row_source"
    "github.com/turbot/tailpipe-plugin-sdk/table"
    "github.com/turbot/tailpipe-plugin-sdk/types"
)

// AlbLogMapper handles the mapping of ALB log lines to structured data
type AlbLogMapper struct{}

func (m *AlbLogMapper) Map(_ context.Context, data any) ([]*rows.AlbAccessLog, error) {
    // Handle string input directly
    lineStr, ok := data.(string)
    if !ok {
        return nil, fmt.Errorf("expected string, got %T", data)
    }

    // Create a new ALB log entry
    logEntry := rows.NewAlbAccessLog()
    
    // Split the line and parse fields
    // Using strings.Split to handle spaces between quoted strings
    var fields []string
    var currentField strings.Builder
    inQuotes := false
    
    for _, char := range lineStr {
        switch char {
        case '"':
            inQuotes = !inQuotes
        case ' ':
            if !inQuotes {
                if currentField.Len() > 0 {
                    fields = append(fields, currentField.String())
                    currentField.Reset()
                }
            } else {
                currentField.WriteRune(char)
            }
        default:
            currentField.WriteRune(char)
        }
    }
    
    // Add the last field
    if currentField.Len() > 0 {
        fields = append(fields, currentField.String())
    }
    
    if len(fields) < 27 { // Minimum required fields
        return nil, fmt.Errorf("invalid number of fields in log entry: got %d, want at least 27", len(fields))
    }

    // Map fields to struct
    logEntry.Type = fields[0]
    
    // Parse timestamp
    timestamp, err := time.Parse(time.RFC3339Nano, fields[1])
    if err != nil {
        return nil, fmt.Errorf("error parsing timestamp: %w", err)
    }
    logEntry.Timestamp = timestamp
    
    // ALB name
    logEntry.AlbName = fields[2]

    // Parse client address
    clientParts := strings.Split(fields[3], ":")
    if len(clientParts) == 2 {
        logEntry.ClientIP = clientParts[0]
        port, err := strconv.Atoi(clientParts[1])
        if err == nil {
            logEntry.ClientPort = port
        }
    }

    // Parse target address
    if fields[4] != "-" {
        targetParts := strings.Split(fields[4], ":")
        if len(targetParts) == 2 {
            logEntry.TargetIP = &targetParts[0]
            port, err := strconv.Atoi(targetParts[1])
            if err == nil {
                logEntry.TargetPort = port
            }
        }
    }

    // Parse processing times
    if rpt, err := strconv.ParseFloat(fields[5], 64); err == nil {
        logEntry.RequestProcessingTime = rpt
    }
    if tpt, err := strconv.ParseFloat(fields[6], 64); err == nil {
        logEntry.TargetProcessingTime = tpt
    }
    if rpt, err := strconv.ParseFloat(fields[7], 64); err == nil {
        logEntry.ResponseProcessingTime = rpt
    }

    // Status codes
    if statusCode, err := strconv.Atoi(fields[8]); err == nil {
        logEntry.AlbStatusCode = &statusCode
    }
    if statusCode, err := strconv.Atoi(fields[9]); err == nil {
        logEntry.TargetStatusCode = &statusCode
    }

    // Bytes
    if received, err := strconv.ParseInt(fields[10], 10, 64); err == nil {
        logEntry.ReceivedBytes = &received
    }
    if sent, err := strconv.ParseInt(fields[11], 10, 64); err == nil {
        logEntry.SentBytes = &sent
    }

    // Remove quotes from quoted fields
    logEntry.Request = strings.Trim(fields[12], "\"")
    logEntry.UserAgent = strings.Trim(fields[13], "\"")
    logEntry.SslCipher = fields[14]
    logEntry.SslProtocol = fields[15]
    logEntry.TargetGroupArn = fields[16]
    logEntry.TraceId = strings.Trim(fields[17], "\"")
    logEntry.DomainName = strings.Trim(fields[18], "\"")
    logEntry.ChosenCertArn = strings.Trim(fields[19], "\"")

    // Parse matched rule priority
    if priority, err := strconv.Atoi(fields[20]); err == nil {
        logEntry.MatchedRulePriority = priority
    }

    // Parse request creation time
    reqCreationTime, err := time.Parse(time.RFC3339Nano, fields[21])
    if err == nil {
        logEntry.RequestCreationTime = reqCreationTime
    }

    // Handle remaining fields
    logEntry.ActionsExecuted = strings.Trim(fields[22], "\"")
    redirectURL := strings.Trim(fields[23], "\"")
    if redirectURL != "-" {
        logEntry.RedirectUrl = &redirectURL
    }
    errorReason := strings.Trim(fields[24], "\"")
    if errorReason != "-" {
        logEntry.ErrorReason = &errorReason
    }
    targetList := strings.Trim(fields[25], "\"")
    if targetList != "-" {
        logEntry.TargetList = &targetList
    }
    targetStatusList := strings.Trim(fields[26], "\"")
    if targetStatusList != "-" {
        logEntry.TargetStatusList = &targetStatusList
    }

    return []*rows.AlbAccessLog{logEntry}, nil
}

func (m *AlbLogMapper) Identifier() string {
    return "alb_log_mapper"
}

type AlbAccessLogTable struct {
    table.TableImpl[*rows.AlbAccessLog, *AlbAccessLogTableConfig, *config.AwsConnection]
}

func NewAlbAccessLogTable() table.Table {
    return &AlbAccessLogTable{}
}

func (t *AlbAccessLogTable) Init(ctx context.Context, connectionSchemaProvider table.ConnectionSchemaProvider, req *types.CollectRequest) error {
    if err := t.TableImpl.Init(ctx, connectionSchemaProvider, req); err != nil {
        return err
    }

    // Set the mapper
    t.Mapper = &AlbLogMapper{}
    return nil
}

func (t *AlbAccessLogTable) Identifier() string {
    return "aws_alb_access_log"
}

func (t *AlbAccessLogTable) GetRowSchema() any {
    return &rows.AlbAccessLog{}
}

func (t *AlbAccessLogTable) GetConfigSchema() parse.Config {
    return &AlbAccessLogTableConfig{}
}

func (t *AlbAccessLogTable) GetSourceOptions(sourceType string) []row_source.RowSourceOption {
    return []row_source.RowSourceOption{
        artifact_source.WithRowPerLine(),
    }
}

func (t *AlbAccessLogTable) EnrichRow(row *rows.AlbAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.AlbAccessLog, error) {
    if sourceEnrichmentFields != nil {
        row.CommonFields = *sourceEnrichmentFields
    }

    // Record standardization
    row.TpID = xid.New().String()
    row.TpSourceType = "aws_alb_access_log"
    row.TpIngestTimestamp = helpers.UnixMillis(time.Now().UnixNano() / int64(time.Millisecond))

    // Set partition
    row.TpPartition = t.Identifier()

    // Extract account ID from target group ARN
    // ARN format: arn:aws:elasticloadbalancing:region:account-id:targetgroup/...
    if row.TargetGroupArn != "" {
        arnParts := strings.Split(row.TargetGroupArn, ":")
        if len(arnParts) >= 5 {
            row.TpIndex = arnParts[4] // account ID is the 5th component (index 4)
        }
    }

    // Fallback to connection's default region if we couldn't get account ID
    if row.TpIndex == "" && t.Connection != nil && t.Connection.DefaultRegion != nil {
        row.TpIndex = *t.Connection.DefaultRegion
    }

    // Set date in yyyy-mm-dd format
    row.TpDate = row.Timestamp.Format("2006-01-02")

    return row, nil
}