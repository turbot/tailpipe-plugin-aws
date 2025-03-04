## Activity Examples

### Daily Request Trends
Count requests per day to identify traffic patterns over time. This query helps monitor daily load balancer usage, detect potential traffic spikes, and understand overall system load across different days.

```sql
select
  strftime(timestamp, '%Y-%m-%d') as request_date,
  count(*) as request_count
from
  aws_clb_access_log
group by
  request_date
order by
  request_date asc;
```

```yaml
folder: ELB
```

### Top 10 Clients by Request Count
List the top 10 client IP addresses making requests. This query helps identify the most active clients, potential sources of high traffic, and can assist in network security monitoring and capacity planning.

```sql
select
  client_ip,
  count(*) as request_count
from
  aws_clb_access_log
group by
  client_ip
order by
  request_count desc
limit 10;
```

```yaml
folder: ELB
```

### Request Distribution by Backend
Analyze how requests are distributed across backend instances. Understanding backend request distribution can help optimize resource allocation, identify potential bottlenecks, and ensure balanced load across your infrastructure.

```sql
select
  backend_ip,
  count(*) as request_count
from
  aws_clb_access_log
where
  backend_ip is not null
group by
  backend_ip
order by
  request_count desc;
```

```yaml
folder: ELB
```

### HTTP Status Code Distribution
Analyze the distribution of HTTP status codes returned by the load balancer. This query provides insights into the overall health of your application, helping you quickly identify error rates and potential issues with your backend services.

```sql
select
  elb_status_code,
  count(*) as response_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 2) as percentage
from
  aws_clb_access_log
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
Identify instances where the load balancer couldn't connect to the backend targets. This query helps detect potential backend infrastructure issues, network problems, or service disruptions that prevent successful request routing.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  client_ip,
  backend_ip,
  request_http_version,
  request_http_method,
  request_url,
  elb_status_code,
  backend_status_code
from
  aws_clb_access_log
where
  backend_status_code is null
  and elb_status_code = 502
order by
  timestamp desc;
```

```yaml
folder: ELB
```

### SSL Cipher Vulnerabilities
Detect usage of deprecated or insecure SSL ciphers. This query helps identify outdated SSL/TLS protocols that may pose security risks, allowing you to upgrade and maintain robust encryption standards.

```sql
select
  ssl_cipher,
  ssl_protocol,
  count(*) as request_count
from
  aws_clb_access_log
where
  ssl_protocol in ('TLSv1.1', 'TLSv1', 'SSLv3', 'SSLv2')
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
Identify potentially suspicious user agents making requests. This query helps detect potential bot traffic, automated scanning tools, or unusual client behaviors that might indicate security probing or potential threats.

```sql
select
  user_agent,
  count(*) as request_count
from
  aws_clb_access_log
where
  user_agent like '%bot%'
  or user_agent like '%curl%'
  or user_agent like '%wget%'
group by
  user_agent
having
  count(*) > 100
order by
  request_count desc;
```

```yaml
folder: ELB
```

## Operational Examples

### Slow Response Times
Top 10 requests with unusually high processing times. This query helps identify performance bottlenecks by highlighting requests that take longer than expected, which can guide optimization efforts and improve overall system responsiveness.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request_http_version,
  request_http_method,
  request_url,
  client_ip,
  backend_ip,
  request_processing_time,
  backend_processing_time,
  response_processing_time,
  (request_processing_time + backend_processing_time + response_processing_time) as total_time
from
  aws_clb_access_log
where
  (request_processing_time + backend_processing_time + response_processing_time) > 1
order by
  total_time desc
limit 10;
```

```yaml
folder: ELB
```

### HTTP Request Method Distribution
Analyze the distribution of HTTP request methods. This query helps understand the types of requests being made to your load balancer, which can provide insights into application usage patterns and potential areas for optimization.

```sql
select
  request_http_method,
  count(*) as request_count
from
  aws_clb_access_log
group by
  request_http_method
order by
  request_count desc;
```

```yaml
folder: ELB
```

### Backend Health Issues
Identify backend instances that are returning a high number of errors. This query helps pinpoint specific backend servers experiencing consistent issues, enabling targeted troubleshooting and potential infrastructure improvements.

```sql
select
  backend_ip,
  backend_status_code,
  count(*) as error_count
from
  aws_clb_access_log
where
  backend_status_code >= 400
group by
  backend_ip,
  backend_status_code
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
Detect periods of unusually high request volume. This query helps identify peak traffic times, potential Denial of Service (DoS) attacks, or unexpected usage patterns that might require infrastructure scaling or further investigation.

```sql
select
  date_trunc('minute', timestamp) as request_minute,
  elb,
  count(*) as request_count
from
  aws_clb_access_log
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
Track requests generating unusually large responses. This query helps identify potential data transfer bottlenecks, content delivery issues, or unusual data transfer patterns that might impact system performance.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request_http_version,
  request_http_method,
  request_url,
  client_ip,
  sent_bytes,
  received_bytes
from
  aws_clb_access_log
where
  sent_bytes > 10485760
order by
  sent_bytes desc;
```

```yaml
folder: ELB
```
