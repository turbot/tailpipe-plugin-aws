## Activity Examples

### Daily Request Trends

Count requests per day to identify traffic patterns over time. This query helps visualize usage trends, detect potential traffic anomalies, and understand the overall load on your application load balancer across different days.

```sql
select
  strftime(timestamp, '%Y-%m-%d') as request_date,
  count(*) as request_count
from
  aws_alb_access_log
group by
  request_date
order by
  request_date asc;
```

```yaml
folder: ELB
```

### Top 10 Clients by Request Count

List the top 10 client IP addresses making requests. This query helps identify the most active clients, potentially revealing heavy users, bot traffic, or unusual access patterns that might require further investigation.

```sql
select
  client_ip,
  count(*) as request_count
from
  aws_alb_access_log
group by
  client_ip
order by
  request_count desc
limit 10;
```

```yaml
folder: ELB
```

### Request Distribution by Target

Analyze how requests are distributed across target instances. This query provides insights into load balancing effectiveness, helping identify potential bottlenecks or uneven traffic distribution among backend servers.

```sql
select
  target_ip,
  count(*) as request_count
from
  aws_alb_access_log
where
  target_ip is not null
group by
  target_ip
order by
  request_count desc;
```

```yaml
folder: ELB
```

### HTTP Status Code Distribution

Analyze the distribution of HTTP status codes returned by the load balancer. This query helps understand the overall health of your application, identifying success rates, client errors, and server errors.

```sql
select
  elb_status_code,
  count(*) as response_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 3) as percentage
from
  aws_alb_access_log
group by
  elb_status_code
order by
  response_count desc;
```

```yaml
folder: ELB
```

## Detection Examples

### Failed Backend Connections

Identify instances where the load balancer couldn't connect to the backend targets. This query helps detect backend infrastructure issues, network problems, or potential service disruptions.

```sql
select
  timestamp,
  elb,
  client_ip,
  target_ip,
  request_url,
  request_http_method,
  request_http_version,
  elb_status_code,
  target_status_code,
  error_reason
from
  aws_alb_access_log
where
  target_status_code is null
  and elb_status_code = 502
order by
  timestamp desc;
```

```yaml
folder: ELB
```

### SSL Cipher Vulnerabilities

Detect usage of deprecated or insecure SSL ciphers. This query helps identify potential security risks by highlighting the use of outdated or vulnerable SSL protocols and ciphers.

```sql
select
  ssl_cipher,
  ssl_protocol,
  count(*) as request_count
from
  aws_alb_access_log
where
  ssl_protocol in ('TLSv1.1', 'TLSv1', 'SSLv3', 'SSLv2') -- Insecure protocols (TLSv1.1, TLSv1, SSLv3, SSLv2)
group by
  ssl_cipher,
  ssl_protocol
order by
  request_count desc;
```

```yaml
folder: ELB
```

### Suspicious User Agents

Identify potentially suspicious user agents making requests. This query helps detect potential bot traffic, web scrapers, or automated scanning tools that might be interacting with your application.

```sql
select
  user_agent,
  count(*) as request_count
from
  aws_alb_access_log
where
  user_agent like '%bot%'
  or user_agent like '%curl%'
  or user_agent like '%wget%'
group by
  user_agent
order by
  request_count desc;
```

```yaml
folder: ELB
```

## Operational Examples

### Slow Response Times

Top 10 requests with unusually high processing times. This query helps identify performance bottlenecks, slow backend services, or potential optimization opportunities in your application infrastructure.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request_url,
  request_http_method,
  request_http_version,
  client_ip,
  target_ip,
  request_processing_time,  -- Time taken by load balancer to process request
  target_processing_time,   -- Time taken by target to process request
  response_processing_time, -- Time taken to process response
  (request_processing_time + target_processing_time + response_processing_time) as total_time
from
  aws_alb_access_log
where
  (request_processing_time + target_processing_time + response_processing_time) > 1 -- Requests taking longer than 1 second
order by
  total_time desc
limit 10;
```

```yaml
folder: ELB
```

### Target Health Issues

Identify targets that are returning a high number of errors. This query helps detect backend service problems, potential server misconfigurations, or application-level issues affecting specific targets.

```sql
select
  target_ip,
  target_status_code,
  count(*) as error_count
from
  aws_alb_access_log
where
  target_status_code >= 400
group by
  target_ip,
  target_status_code
having
  count(*) > 100
order by
  error_count desc;
```

```yaml
folder: ELB
```

## Volume Examples

### High Traffic Periods

Detect periods of unusually high request volume. This query helps identify traffic peaks, potential denial-of-service scenarios, or unexpected surges in application usage.

```sql
select
  date_trunc('minute', timestamp) as request_minute,
  elb,
  count(*) as request_count
from
  aws_alb_access_log
group by
  request_minute,
  elb
having
  count(*) > 1000
order by
  request_count desc;
```

```yaml
folder: ELB
```

### Large Response Sizes

Track requests generating unusually large responses. This query helps identify potential data exfiltration attempts, performance issues with large payloads, or unusual data transfer patterns.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request_url,
  request_http_method,
  request_http_version,
  client_ip,
  sent_bytes,
  received_bytes
from
  aws_alb_access_log
where
  sent_bytes > 10485760 -- 10MB
order by
  sent_bytes desc;
```

```yaml
folder: ELB
```

### Count Requests by Listener Type

Identify how many requests were received for each listener type. This query provides insights into the types of traffic your load balancer is handling, helping understand protocol usage and potential security configurations.

```sql
select
  listener_type,
  count(*) as request_count
from (
  select
    case
      when request_http_method = 'GET' and request_url like 'http://%' then 'HTTP'
      when request_http_method = 'GET' and request_url like 'https://%' and request_http_version = 'HTTP/1.1' then 'HTTPS'
      when request_http_method = 'GET' and request_url like 'https://%' and request_http_version = 'HTTP/2' then 'HTTP/2'
      when request_http_method = 'GET' and request_url like 'ws://%' then 'WebSockets'
      when request_http_method = 'GET' and request_url like 'wss://%' then 'Secured WebSockets'
      else 'Other'
    end as listener_type
  from aws_alb_access_log
) t
group by 
  listener_type
order by
  request_count desc;
```

```yaml
folder: ELB
```

### Count by HTTP Method

Identify the distribution of requests by HTTP method. This query helps understand how clients interact with your application, providing insights into usage patterns and potential security risks.

```sql
select
  request_http_method,
  count(*) as request_count
from
  aws_alb_access_log
where
  request_http_method is not null
group by
  request_http_method
order by
  request_count desc;
```

```yaml
folder: ELB
```

### Requests with Invalid Cookies

Identify requests with invalid cookies. This query helps detect potential security issues, client-side problems, or application configuration errors related to cookie handling.

```sql
select
  timestamp,
  client_ip,
  request_url,
  request_http_method,
  request_http_version,
  error_reason
from
  aws_alb_access_log
where
  error_reason = 'AuthInvalidCookie'
order by
  timestamp desc;
```

```yaml
folder: ELB
```

### Requests with Invalid Lambda Response

Identify requests with an invalid Lambda response. This query helps troubleshoot issues with serverless function integrations and detect potential problems in Lambda function implementations.

```sql
select
  timestamp,
  client_ip,
  request_url,
  request_http_method,
  request_http_version,
  error_reason
from
  aws_alb_access_log
where
  error_reason = 'LambdaInvalidResponse'
order by
  timestamp desc;
```

```yaml
folder: ELB
```

### Requests Blocked by WAF Rules

Identify requests blocked by WAF rules. This query helps understand potential security threats, analyze the effectiveness of web application firewall configurations, and track potential malicious traffic.

```sql
select
  timestamp,
  client_ip,
  request_url,
  request_http_method,
  request_http_version,
  error_reason
from
  aws_alb_access_log
where
  actions_executed @> ['waf']
  and elb_status_code = 403
order by
  timestamp desc;
```

```yaml
folder: WAF
```

### Requests Exceeding Maximum Allowed Body Size for Lambda

Identify requests that exceed the maximum allowed body size for Lambda. This query helps detect potential issues with request payload sizes and understand limitations in serverless function integrations.

```sql
select
  timestamp,
  client_ip,
  request_url,
  request_http_method,
  request_http_version,
  sent_bytes,
  received_bytes,
  error_reason
from
  aws_alb_access_log
where
  error_reason = 'LambdaResponseTooLarge'
order by
  sent_bytes desc;
```

```yaml
folder: ELB
```