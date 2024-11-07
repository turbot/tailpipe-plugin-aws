# ALB Log Analysis Queries and Results

## Client Type Distribution
Analyzes the distribution of different client types accessing the application.

```sql
SELECT
    CASE
        WHEN user_agent LIKE '%Mobile%' THEN 'Mobile'
        WHEN user_agent LIKE '%Chrome%' THEN 'Chrome'
        WHEN user_agent LIKE '%Firefox%' THEN 'Firefox'
        WHEN user_agent LIKE '%Safari%' THEN 'Safari'
        WHEN user_agent LIKE '%bot%' OR user_agent LIKE '%Bot%' THEN 'Bot'
        ELSE 'Other'
    END as client_type,
    COUNT(*) as request_count
FROM aws_alb_access_log
GROUP BY client_type
ORDER BY request_count DESC;
```

Results show most traffic comes from "Other" clients (6,337 requests), followed by Mobile (1,842) and Chrome (1,821).

## SSL/TLS Configuration Analysis
Examines the distribution of SSL protocols and cipher suites.

```sql
SELECT
    ssl_protocol,
    ssl_cipher,
    COUNT(*) as count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage
FROM aws_alb_access_log
WHERE ssl_protocol != '-'
GROUP BY ssl_protocol, ssl_cipher
ORDER BY count DESC;
```

Key findings:
- TLSv1.2 and TLSv1.3 are evenly distributed
- Most common cipher is TLS_AES_128_GCM_SHA256 (17.21%)
- All modern cipher suites used

## Traffic Distribution by Hour
Shows the distribution of requests across different hours of the day.

```sql
SELECT
    EXTRACT(hour FROM timestamp) as hour_of_day,
    COUNT(*) as request_count
FROM aws_alb_access_log
GROUP BY hour_of_day
ORDER BY request_count DESC;
```

Traffic is relatively evenly distributed between hours 0-14, with a drop-off at hour 15.

## Error Rate Analysis
Examines error rates by hour.

```sql
SELECT
    DATE_TRUNC('hour', timestamp) as hour,
    COUNT(*) as total_requests,
    SUM(CASE WHEN alb_status_code >= 400 THEN 1 ELSE 0 END) as error_count,
    ROUND(SUM(CASE WHEN alb_status_code >= 400 THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2) as error_rate
FROM aws_alb_access_log
GROUP BY hour
ORDER BY hour;
```

Error rates consistently range between 19-23% across all hours.

## ALB Performance Analysis
Analyzes performance metrics across different ALBs.

```sql
SELECT
    tp_index as alb_name,
    COUNT(*) as request_count,
    COUNT(DISTINCT client_ip) as unique_clients,
    ROUND(AVG(request_processing_time + target_processing_time + response_processing_time), 3) as avg_total_time
FROM aws_alb_access_log
GROUP BY tp_index
ORDER BY request_count DESC;
```

Findings:
- prod-web-alb handles most traffic (3,383 requests)
- staging-alb has higher average processing time (0.523s)
- prod environments maintain better performance (~0.35s)

## Security Analysis

### SQL Injection Attempts
```sql
SELECT
    client_ip,
    user_agent,
    COUNT(*) AS attempt_count,
    STRING_AGG(sample_request, ' | ') AS sample_requests,
    MIN(timestamp) AS first_seen,
    MAX(timestamp) AS last_seen
FROM filtered_requests
WHERE row_num <= 3
GROUP BY client_ip, user_agent
HAVING COUNT(*) > 1
ORDER BY attempt_count DESC;
```

Detected multiple SQL injection attempts from IPs in ranges:
- 185.181.x.x
- 193.27.228.x
- 45.155.205.x

### Cross-ALB Attack Patterns
```sql
WITH suspicious_paths AS (...)
SELECT
    client_ip,
    user_agent,
    COUNT(DISTINCT alb_name) as albs_targeted,
    COUNT(*) as total_probes,
    STRING_AGG(DISTINCT alb_name, ', ') as targeted_albs,
    STRING_AGG(sample_request, ' | ') as sample_requests,
    MIN(timestamp) as first_seen,
    MAX(timestamp) as last_seen,
    EXTRACT(MINUTES FROM MAX(timestamp) - MIN(timestamp)) as campaign_duration_mins
FROM suspicious_paths
WHERE row_num <= 3
GROUP BY client_ip, user_agent
HAVING
    COUNT(DISTINCT alb_name) > 1 AND
    COUNT(*) >= 3
ORDER BY albs_targeted DESC, total_probes DESC;
```

Key findings:
- Multiple IPs targeting all three ALBs
- Campaign durations ranging from 13 to 54 minutes
- Common attack tools identified: subfinder, WhatWeb, dirbuster, Nmap
- Systematic probing of infrastructure suggesting coordinated attacks

### Notable Attack Patterns
- Sustained campaigns lasting 30+ minutes
- Multiple tools used per attacker
- Systematic targeting of all ALB environments
- Focus on common vulnerabilities (actuator, debug endpoints, config files)