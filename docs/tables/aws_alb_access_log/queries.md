## Activity Examples

### Daily request trends

Count requests per day to identify traffic patterns over time.

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

### Top 10 clients by request count

List the top 10 client IP addresses making requests.

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

### Request distribution by target

Analyze how requests are distributed across target instances.

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

### HTTP status code distribution

Analyze the distribution of HTTP status codes returned by the load balancer.

```sql
select
  elb_status_code,
  count(*) as response_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 2) as percentage
from
  aws_alb_access_log
group by
  elb_status_code
order by
  response_count desc;
```

## Detection Examples

### Failed backend connections

Identify instances where the load balancer couldn't connect to the backend targets.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  client_ip,
  target_ip,
  request,
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

### SSL cipher vulnerabilities

Detect usage of deprecated or insecure SSL ciphers.

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

### Suspicious user agents

Identify potentially suspicious user agents making requests.

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
having
  count(*) > 100
order by
  request_count desc;
```

## Operational Examples

### Slow response times

Top 10 requests with unusually high processing times.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request,
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

### Target health issues

Identify targets that are returning a high number of errors.

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

## Volume Examples

### High traffic periods

Detect periods of unusually high request volume.

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

### Large response sizes

Track requests generating unusually large responses.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request,
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

### Count requests by Listener Type (HTTP, HTTPS, HTTP/2, WebSockets, etc.)

Identify how many requests were received for each listener type.

```sql
select
  listener_type,
  count(*) as request_count
from (
  select
    case
      when request like 'GET http://%' then 'HTTP'
      when request like 'GET https://%' and ssl_protocol is not null then 'HTTPS'
      when request like 'GET https://%' and ssl_protocol is null then 'HTTP/2'
      when request like 'GET ws://%' then 'WebSockets'
      when request like 'GET wss://%' then 'Secured WebSockets'
      else 'Other'
    end as listener_type
  from aws_alb_access_log
)
group by listener_type
order by request_count desc;
```

### Requests with Invalid Cookies

Identify requests with invalid cookies.

```sql
select
  timestamp,
  client_ip,
  request,
  error_reason
from
  aws_alb_access_log
where
  error_reason = 'AuthInvalidCookie'
order by
  timestamp desc;
```

### Requests with Invalid Lambda response

Identify requests with an invalid Lambda response.

```sql
select
  timestamp,
  client_ip,
  request,
  error_reason
from
  aws_alb_access_log
where
  elb_status_code = 502
  and error_reason like 'LambdaInvalidResponse'
order by
  timestamp desc;
```

### Requests blocked by WAF rules

Identify requests blocked by WAF rules.

```sql
select
  timestamp,
  client_ip,
  request,
  error_reason
from
  aws_alb_access_log
where
  actions_executed = 'waf'
  and elb_status_code = 403
order by
  timestamp desc;
```

### Requests Exceeding Maximum Allowed Body Size for Lambda

Identify requests that exceed the maximum allowed body size for Lambda.

```sql
select
  timestamp,
  client_ip,
  request,
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

### Requests with Unhandled Lambda Response

Identify requests with an unhandled Lambda response.

```sql
select
  timestamp,
  client_ip,
  request,
  error_reason
from
  aws_alb_access_log
where
  error_reason = 'LambdaUnhandled'
order by 
  timestamp desc;
```