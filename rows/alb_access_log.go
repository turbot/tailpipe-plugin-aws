package rows

import (
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/turbot/tailpipe-plugin-sdk/enrichment"
    "github.com/turbot/tailpipe-plugin-sdk/helpers"
)

// AlbAccessLog represents a single ALB access log entry with enrichment fields
type AlbAccessLog struct {
    // Embed required enrichment fields
    enrichment.CommonFields

    // Standard ALB fields
    Type                   string    `json:"type"`
    Timestamp             time.Time `json:"timestamp"`
    AlbName               string    `json:"alb_name"`
    ClientIP              string    `json:"client_ip"`
    ClientPort            int       `json:"client_port"`
    TargetIP              *string   `json:"target_ip,omitempty"`
    TargetPort            int       `json:"target_port"`
    RequestProcessingTime float64   `json:"request_processing_time"`
    TargetProcessingTime  float64   `json:"target_processing_time"`
    ResponseProcessingTime float64   `json:"response_processing_time"`
    AlbStatusCode         *int      `json:"alb_status_code,omitempty"`
    TargetStatusCode      *int      `json:"target_status_code,omitempty"`
    ReceivedBytes         *int64    `json:"received_bytes,omitempty"`
    SentBytes             *int64    `json:"sent_bytes,omitempty"`
    Request               string    `json:"request"`
    UserAgent             string    `json:"user_agent"`
    SslCipher            string    `json:"ssl_cipher"`
    SslProtocol          string    `json:"ssl_protocol"`
    TargetGroupArn       string    `json:"target_group_arn"`
    TraceId              string    `json:"trace_id"`
    DomainName           string    `json:"domain_name"`
    ChosenCertArn        string    `json:"chosen_cert_arn"`
    MatchedRulePriority  int       `json:"matched_rule_priority"`
    RequestCreationTime  time.Time `json:"request_creation_time"`
    ActionsExecuted      string    `json:"actions_executed"`
    RedirectUrl          *string   `json:"redirect_url,omitempty"`
    ErrorReason          *string   `json:"error_reason,omitempty"`
    TargetList           *string   `json:"target_list,omitempty"`
    TargetStatusList     *string   `json:"target_status_list,omitempty"`
    Classification       *string   `json:"classification,omitempty"`
    ClassificationReason *string   `json:"classification_reason,omitempty"`
}

// NewAlbAccessLog creates a new ALB access log entry
func NewAlbAccessLog() *AlbAccessLog {
    return &AlbAccessLog{}
}

// InitialiseFromMap initializes the struct from a map of string values
func (l *AlbAccessLog) InitialiseFromMap(m map[string]string) error {
    for key, value := range m {
        switch key {
        case "timestamp":
            ts, err := time.Parse(time.RFC3339, value)
            if err != nil {
                return fmt.Errorf("error parsing timestamp: %w", err)
            }
            l.Timestamp = ts
            l.TpTimestamp = helpers.UnixMillis(ts.UnixNano() / int64(time.Millisecond))
        case "type":
            l.Type = value
        case "alb":
            l.AlbName = value
        case "client":
            if value != "-" && strings.Contains(value, ":") {
                parts := strings.Split(value, ":")
                if len(parts) != 2 {
                    return fmt.Errorf("invalid client address format: %s", value)
                }
                l.ClientIP = parts[0]
                l.TpSourceIP = &l.ClientIP
                l.TpIps = append(l.TpIps, l.ClientIP)
                
                port, err := strconv.Atoi(parts[1])
                if err != nil {
                    return fmt.Errorf("error parsing client port: %w", err)
                }
                l.ClientPort = port
            }
        case "target":
            if value != "-" && strings.Contains(value, ":") {
                parts := strings.Split(value, ":")
                if len(parts) != 2 {
                    return fmt.Errorf("invalid target address format: %s", value)
                }
                ip := parts[0]
                l.TargetIP = &ip
                l.TpDestinationIP = &ip
                l.TpIps = append(l.TpIps, ip)
                
                port, err := strconv.Atoi(parts[1])
                if err != nil {
                    return fmt.Errorf("error parsing target port: %w", err)
                }
                l.TargetPort = port
            }
        case "request_processing_time":
            rpt, err := strconv.ParseFloat(value, 64)
            if err != nil {
                return fmt.Errorf("error parsing request_processing_time: %w", err)
            }
            l.RequestProcessingTime = rpt
        case "target_processing_time":
            tpt, err := strconv.ParseFloat(value, 64)
            if err != nil {
                return fmt.Errorf("error parsing target_processing_time: %w", err)
            }
            l.TargetProcessingTime = tpt
        case "response_processing_time":
            rpt, err := strconv.ParseFloat(value, 64)
            if err != nil {
                return fmt.Errorf("error parsing response_processing_time: %w", err)
            }
            l.ResponseProcessingTime = rpt
        case "alb_status_code":
            if value != "-" {
                code, err := strconv.Atoi(value)
                if err != nil {
                    return fmt.Errorf("error parsing alb_status_code: %w", err)
                }
                l.AlbStatusCode = &code
            }
        case "target_status_code":
            if value != "-" {
                code, err := strconv.Atoi(value)
                if err != nil {
                    return fmt.Errorf("error parsing target_status_code: %w", err)
                }
                l.TargetStatusCode = &code
            }
        case "received_bytes":
            if value != "-" {
                bytes, err := strconv.ParseInt(value, 10, 64)
                if err != nil {
                    return fmt.Errorf("error parsing received_bytes: %w", err)
                }
                l.ReceivedBytes = &bytes
            }
        case "sent_bytes":
            if value != "-" {
                bytes, err := strconv.ParseInt(value, 10, 64)
                if err != nil {
                    return fmt.Errorf("error parsing sent_bytes: %w", err)
                }
                l.SentBytes = &bytes
            }
        case "request":
            l.Request = value
        case "user_agent":
            l.UserAgent = value
        case "ssl_cipher":
            l.SslCipher = value
        case "ssl_protocol":
            l.SslProtocol = value
        case "target_group_arn":
            l.TargetGroupArn = value
            l.TpAkas = append(l.TpAkas, value)
        case "trace_id":
            l.TraceId = value
        case "domain_name":
            l.DomainName = value
            if value != "-" {
                l.TpDomains = append(l.TpDomains, value)
            }
        case "chosen_cert_arn":
            l.ChosenCertArn = value
        case "matched_rule_priority":
            if value != "-" {
                priority, err := strconv.Atoi(value)
                if err != nil {
                    return fmt.Errorf("error parsing matched_rule_priority: %w", err)
                }
                l.MatchedRulePriority = priority
            }
        case "request_creation_time":
            ts, err := time.Parse(time.RFC3339, value)
            if err != nil {
                return fmt.Errorf("error parsing request_creation_time: %w", err)
            }
            l.RequestCreationTime = ts
        case "actions_executed":
            l.ActionsExecuted = value
        case "redirect_url":
            if value != "-" {
                l.RedirectUrl = &value
            }
        case "error_reason":
            if value != "-" {
                l.ErrorReason = &value
            }
        case "target_list":
            if value != "-" {
                l.TargetList = &value
            }
        case "target_status_list":
            if value != "-" {
                l.TargetStatusList = &value
            }
        case "classification":
            if value != "-" {
                l.Classification = &value
            }
        case "classification_reason":
            if value != "-" {
                l.ClassificationReason = &value
            }
        }
    }
    return nil
}