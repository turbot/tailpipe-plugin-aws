# AWS ALB Access Log Table - Tailpipe Plugin

The `aws_alb_access_log` table is a powerful component of the Tailpipe AWS plugin that enables SQL-based analysis of AWS Application Load Balancer (ALB) access logs. This table provides a structured way to query and analyze ALB logs, helping you understand traffic patterns, troubleshoot issues, and monitor application performance.

## Features

- **Full ALB Log Field Support**: Captures all ALB access log fields including request details, processing times, SSL information, and more
- **Automatic Type Conversion**: Handles parsing and conversion of fields into appropriate data types (timestamps, integers, etc.)
- **Enriched Fields**: Adds standardized fields like `tp_id`, `tp_timestamp`, and `tp_source_type` for consistent analysis
- **IP Address Analysis**: Extracts and standardizes both client and target IP addresses
- **Integration with AWS S3**: Directly reads ALB logs from S3 buckets where AWS stores them

## Configuration

### Plugin Configuration

First, configure the AWS plugin in your Tailpipe configuration:

```hcl
connection "aws" {
  # Authentication (choose one method)
  profile      = "default"              # Use AWS profile
  access_key   = "YOUR_ACCESS_KEY"      # Or provide access key directly
  secret_key   = "YOUR_SECRET_KEY"      # With secret key
  
  # Optional settings
  default_region = "us-east-1"          # Default AWS region
  regions        = ["us-east-1", "us-west-2"]  # Multiple regions if needed
}
```

### Table Configuration

Configure the ALB access log table source:

```hcl
source "aws_alb_access_log" "production_logs" {
  # S3 bucket where ALB logs are stored
  bucket = "my-alb-logs"
  
  # Optional: Custom log format if not using default
  log_format = "$type $timestamp $elb ..."
  
  # Optional: File pattern for specific log files
  file_pattern = "AWSLogs/*/elasticloadbalancing/*/*"
  
  # Optional: Timezone for parsing timestamps
  timezone = "UTC"
}
```

## Schema

The table provides the following columns:

### Standard Fields
| Column Name | Type | Description |
|------------|------|-------------|
| type | string | The type of request (typically 'http' or 'https') |
| timestamp | timestamp | When the request was received |
| alb_name | string | Name of the ALB |
| client_ip | string | Client's IP address |
| client_port | int | Client's port number |
| target_ip | string | Target instance IP |
| target_port | int | Target instance port |

### Processing Times
| Column Name | Type | Description |
|------------|------|-------------|
| request_processing_time | float64 | Time from request received to sent to target |
| target_processing_time | float64 | Time target took to process request |
| response_processing_time | float64 | Time from receiving target response to sending to client |

### Status and Bytes
| Column Name | Type | Description |
|------------|------|-------------|
| alb_status_code | int | HTTP status code from ALB |
| target_status_code | int | HTTP status code from target |
| received_bytes | int64 | Number of bytes received from client |
| sent_bytes | int64 | Number of bytes sent to client |

### SSL/TLS Information
| Column Name | Type | Description |
|------------|------|-------------|
| ssl_cipher | string | SSL cipher used |
| ssl_protocol | string | SSL protocol version |
| chosen_cert_arn | string | ARN of the certificate used |

### Request Details
| Column Name | Type | Description |
|------------|------|-------------|
| request | string | HTTP request line |
| user_agent | string | Client's user agent |
| domain_name | string | Requested domain name |
| target_group_arn | string | ARN of the target group |

## Example Queries

### 1. Monitor Response Times

Find slow requests (>5s total processing time):
```sql
SELECT 
    timestamp,
    client_ip,
    request,
    request_processing_time + target_processing_time + response_processing_time as total_time,
    alb_status_code
FROM 
    aws_alb_access_log
WHERE 
    (request_processing_time + target_processing_time + response_processing_time) > 5
ORDER BY 
    total_time DESC
LIMIT 10;
```

### 2. Analyze Error Rates

Calculate error rates by target group:
```sql
SELECT 
    target_group_arn,
    COUNT(*) as total_requests,
    SUM(CASE WHEN alb_status_code >= 500 THEN 1 ELSE 0 END) as error_count,
    ROUND(100.0 * SUM(CASE WHEN alb_status_code >= 500 THEN 1 ELSE 0 END) / COUNT(*), 2) as error_rate
FROM 
    aws_alb_access_log
GROUP BY 
    target_group_arn
ORDER BY 
    error_rate DESC;
```

### 3. Traffic Analysis

View requests by hour:
```sql
SELECT 
    DATE_TRUNC('hour', timestamp) as hour,
    COUNT(*) as requests,
    SUM(sent_bytes + received_bytes) / (1024*1024) as total_mb
FROM 
    aws_alb_access_log
GROUP BY 
    hour
ORDER BY 
    hour;
```

## Enrichment Fields

The table automatically adds several enrichment fields:

| Field | Description |
|-------|-------------|
| tp_id | Unique identifier for each log entry |
| tp_timestamp | Normalized timestamp in Unix milliseconds |
| tp_source_type | Always "aws_alb_access_log" |
| tp_source_ip | Client IP address |
| tp_destination_ip | Target IP address |
| tp_ips | Array of all IPs in the log entry |
| tp_domains | Array of domain names from the request |

## Troubleshooting

### Common Issues

1. **Missing Fields**: If certain fields are consistently null, verify your ALB logging configuration includes all desired fields.

2. **Parsing Errors**: If you see parse errors:
   - Check the log format matches ALB's output
   - Verify timezone settings
   - Ensure log files aren't corrupted

3. **Performance Issues**: For large log sets:
   - Use appropriate time-based filters
   - Consider partitioning data
   - Index frequently queried columns

### Debugging

Enable debug logging for more detailed information:

```bash
TAILPIPE_LOG_LEVEL=debug tailpipe collect aws_alb_access_log.production_logs
```

## Best Practices

1. **Data Management**
   - Regularly archive old logs
   - Use appropriate S3 lifecycle policies
   - Consider compressing logs

2. **Query Optimization**
   - Filter by date ranges when possible
   - Use appropriate indexes
   - Avoid SELECT * for large datasets

3. **Monitoring**
   - Set up alerts for high error rates
   - Monitor response times
   - Track resource usage

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch
3. Add your changes
4. Submit a pull request

## License

This plugin is licensed under [appropriate license]

## Support

For issues and questions:
- GitHub Issues: [link to issues]
- Documentation: [link to docs]
- Community Forum: [link to forum]