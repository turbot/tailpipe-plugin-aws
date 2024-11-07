# ALB vs ELB Log Table Implementation Comparison

## Structural Differences

### Table Implementation

```go
// ALB Implementation
type AlbAccessLogTable struct {
    table.TableImpl[*rows.AlbAccessLog, *AlbAccessLogTableConfig, *config.AwsConnection]
}

// ELB Implementation
type ElbAccessLogTable struct {
    table.TableImpl[*rows.ElbAccessLog, *ElbAccessLogTableConfig, *config.AwsConnection]
}
```

### Configuration

```go
// ALB Configuration is more extensive
type AlbAccessLogTableConfig struct {
    Remain hcl.Body `hcl:",remain" json:"-"`
    LogFormat *string `json:"log_format,omitempty" hcl:"log_format,optional"`
    FilePattern *string `json:"file_pattern,omitempty" hcl:"file_pattern,optional"`
    Timezone *string `json:"timezone,omitempty" hcl:"timezone,optional"`
}

// ELB Configuration is minimal
type ElbAccessLogTableConfig struct {
}
```

## Data Processing Differences

### Parsing Approach

#### ALB
- Uses custom mapper with detailed string parsing
- Handles quoted strings explicitly
- More flexible parsing logic for complex fields

```go
func (m *AlbLogMapper) Map(_ context.Context, data any) ([]*rows.AlbAccessLog, error) {
    // Custom string parsing with quote handling
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
    // ... field processing
}
```

#### ELB
- Uses generic delimited line mapper
- Relies on predefined format string
- Simpler but less flexible parsing

```go
const elbLogFormat = `$type $timestamp $elb $client $target $request_processing_time...`

func (c *ElbAccessLogTable) initMapper() {
    c.Mapper = table.NewDelimitedLineMapper(rows.NewElbAccessLog, elbLogFormat, elbLogFormatNoConnTrace)
}
```

## Field Handling Differences

### ALB-Specific Fields
```go
type AlbAccessLog struct {
    // ALB-specific fields
    AlbName               string    `json:"alb_name"`
    MatchedRulePriority  int       `json:"matched_rule_priority"`
    RequestCreationTime  time.Time `json:"request_creation_time"`
    Classification       *string   `json:"classification,omitempty"`
    ClassificationReason *string   `json:"classification_reason,omitempty"`
}
```

### ELB-Specific Fields
```go
type ElbAccessLog struct {
    // ELB-specific fields
    Elb                    string    `json:"elb"`
    ConnTraceID            string    `json:"conn_trace_id"`
}
```

## Error Handling

### ALB
- More detailed error handling with specific error messages
- Validation for complex field combinations
```go
if len(fields) < 27 { // Minimum required fields
    return nil, fmt.Errorf("invalid number of fields in log entry: got %d, want at least 27", len(fields))
}
```

### ELB
- Simpler error handling relying on delimiter mapper
- Less validation of field relationships

## Enrichment Implementation

### ALB
```go
func (t *AlbAccessLogTable) EnrichRow(row *rows.AlbAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.AlbAccessLog, error) {
    // ALB-specific enrichment
    if t.Connection != nil && t.Connection.DefaultRegion != nil {
        row.TpIndex = *t.Connection.DefaultRegion
    } else {
        row.TpIndex = "unknown"
    }
    // ... other enrichment
}
```

### ELB
```go
func (c *ElbAccessLogTable) EnrichRow(row *rows.ElbAccessLog, sourceEnrichmentFields *enrichment.CommonFields) (*rows.ElbAccessLog, error) {
    // ELB-specific enrichment
    row.TpPartition = c.Identifier()
    if row.TpIndex == "" {
        row.TpIndex = c.Identifier()
    }
    // ... other enrichment
}
```

## Key Differences Summary

1. **Configuration Flexibility**
   - ALB: More configurable with options for log format, file patterns, and timezone
   - ELB: Minimal configuration, relies on fixed formats

2. **Parsing Strategy**
   - ALB: Custom parser with quote handling and flexible field extraction
   - ELB: Uses generic delimiter-based parsing

3. **Field Support**
   - ALB: Supports newer fields like rule priority and classification
   - ELB: Focuses on classic ELB fields

4. **Error Handling**
   - ALB: More comprehensive error checking and validation
   - ELB: Basic error handling through delimiter mapper

5. **Enrichment**
   - ALB: Region-aware enrichment
   - ELB: Basic identifier-based enrichment

6. **Performance Implications**
   - ALB: More CPU-intensive due to custom parsing
   - ELB: Potentially faster due to simpler delimiter-based parsing