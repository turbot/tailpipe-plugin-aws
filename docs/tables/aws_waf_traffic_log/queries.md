## Activity Examples

### Daily request trends

Count events per day to identify request trends over time.

```sql
select
  strftime(timestamp, '%Y-%m-%d') as access_date,
  count(*) as requests
from
  aws_waf_traffic_log
group by
  access_date
order by
  access_date asc;
```

### Top 10 client IPs blocked by AWS WAF

Identify the top 10 client IP addresses that were blocked by AWS WAF, helping to detect high-volume attack sources or suspicious traffic patterns.

```sql
select
  http_source_name,
  http_source_id,
  http_request.clientIp as client_ip,
  count(*) as request_count
from
  aws_waf_traffic_log
where
  action = 'BLOCK'
group by
  http_source_name,
  http_source_id,
  client_ip
order by
  request_count desc
limit 10;
```

### Top HTTP methods by source

Analyzes the AWS WAF traffic logs to identify the most frequently used HTTP methods (GET, POST, PUT, DELETE, etc.) across different sources.

```sql
select
  http_source_name,
  http_source_id,
  http_request.httpMethod,
  count(*) as request_count
from
  aws_waf_traffic_log
group by
  http_source_name,
  http_source_id,
  http_request.httpMethod
order by
  request_count desc;
```

## Operational Examples

### Requests blocked by specific WAF rules

This query retrieves the number of requests blocked by each WAF rule. It groups the blocked requests by WAF rule name and action type, providing insights into which rules are most actively blocking traffic. This helps in fine-tuning security policies, identifying potential threats, and optimizing WAF rules.

```sql
with blocked_rule as (
  select
    timestamp,
    unnest(
      from_json(rule_group_list, '["JSON"]')
    ) as rule_group,
    action,
    http_request
  from
    aws_waf_traffic_log
   where
    action = 'BLOCK'
)
select
  timestamp,
  action,
  http_request.clientIp as client_ip,
  http_request.uri,
  (rule_group -> 'terminatingRule' ->> 'ruleId') as terminating_rule_id,
  (rule_group -> 'terminatingRule' ->> 'action') as terminating_rule_action
from
  blocked_rule;
```

### Most Targeted URLs

Finds which URLs or endpoints are most frequently targeted URL. This helps identify high-risk areas in your application.

```sql
select
  http_source_name,
  http_source_id,
  http_request.uri,
  action,
  count(*) as request_count
from
  aws_waf_traffic_log
group by
  http_source_name,
  http_source_id,
  http_request.uri,
  action
order by
  request_count desc
limit 10;
```

## Detection Examples

### Requests without labels

Labels in AWS WAF are metadata tags applied to web requests that match specific rules within a Web ACL. These labels provide context on why a request was flagged, blocked, or allowed.

```sql
select
  timestamp,
  http_request.clientIp as client_ip,
  action,
  request_headers_inserted
from
  aws_waf_traffic_log,
where
  labels is null
order by
  timestamp;
```

### XSS (Cross-Site Scripting) Attack Attempts

Finds requests blocked due to Cross-Site Scripting (XSS) attempts. These attacks try to inject malicious scripts into your web pages.

```sql
with blocked_rule as (
  select
    timestamp,
    http_source_name,
    http_source_id,
    unnest(
    from_json(rule_group_list, '["JSON"]')
    ) as rule_group,
    action,
    http_request
  from
    aws_waf_traffic_log
  where
    action = 'BLOCK'
)
select
  timestamp,
  http_source_name,
  http_source_id,
  http_request.clientIp as client_ip,
  http_request.httpMethod as http_method,
  http_request.uri,
from
  blocked_rule
where
  (rule_group -> 'terminatingRule' ->> 'ruleId') like '%XSS%';
```

### Detect high volume of blocked requests

Identify IPs generating a high volume of blocked requests.

```sql
select
  http_request.clientIp as client_ip,
  count(*) as block_count
from
  aws_waf_traffic_log
where
  action = 'BLOCK'
group by
  client_ip
order by
  block_count desc
limit 10;
```

### Detect requests triggering multiple WAF rules

Find requests that matched multiple non-terminating WAF rules.

```sql
select
  timestamp,
  http_request.clientIp as client_ip,
  http_source_name,
  http_source_id,
  json_array_length(non_terminating_matching_rules) as matched_rules,
from
  aws_waf_traffic_log
where
  json_array_length(non_terminating_matching_rules) > 1
order by
  matched_rules desc;
```

### Detect IPs bypassing CAPTCHA challenges

Find IPs that repeatedly triggered CAPTCHA but continued making requests.

```sql
select
  http_request.clientIp as client_ip,
  count(*) as captcha_challenges
from
  aws_waf_traffic_log
where
 terminating_rule_type = 'CAPTCHA'
group by
  client_ip
order by
  captcha_challenges desc;
```

### Detect IP addresses in request headers

Identify web requests where the client IP is present in HTTP headers such as X-Forwarded-For, Client-IP, True-Client-IP, or X-Real-IP. This helps track proxied traffic, potential IP spoofing, or forwarded requests passing through CDN, load balancers, or proxies.

```sql
select
  timestamp,
  http_request.clientIp as client_ip,
  action,
  request_headers_inserted
from 
  aws_waf_traffic_log
where 
  (request_headers_inserted ->> 'name') in ('X-Forwarded-For', 'Client-IP', 'True-Client-IP', 'X-Real-IP')
  and (request_headers_inserted ->> 'value') is not null
order by 
  timestamp desc;
```

## Volume Examples

### High volume of blocked request

Analyze high-volume blocked requests and provide statistics on blocked traffic trends.

```sql
select 
  date_trunc('hour', timestamp) as request_hour,
  http_request.clientIp as client_ip,
  http_source_name,
  http_request.uri as request_uri,
  count(*) as block_count
from 
  aws_waf_traffic_log
where 
  action = 'BLOCK'
group by 
  request_hour, 
  client_ip,
  http_source_name,
  request_uri
having 
  count(*) > 100
order by 
  block_count desc;
```

### Most frequently triggered WAF rules

Identify which WAF rules are blocking/allowing the most traffic.

```sql
select
  terminating_rule_id,
  terminating_rule_type,
  count(*) as rule_trigger_count
from
  aws_waf_traffic_log
group by 
  terminating_rule_id,
  terminating_rule_type
having 
  count(*) > 100
order by
  rule_trigger_count desc;
```

## Baseline Examples

### Requests triggering multiple rules

Find requests that matched more than one rule.

```sql
select 
  http_request.clientIp as client_ip,
  json_array_length(non_terminating_matching_rules) as matched_rules,
  count(*) as request_count
from 
  aws_waf_traffic_log
where 
  json_array_length(non_terminating_matching_rules) > 1
group by 
  client_ip,
  matched_rules
order by 
  request_count desc;
```

### Top requests by country

Find the top countries where WAF rules are being triggered.

```sql
select 
  http_request.country as country,
  count(*) as request_count
from 
  aws_waf_traffic_log
where 
  terminating_rule_id is not null
group by 
  country
order by 
  request_count desc;
```
