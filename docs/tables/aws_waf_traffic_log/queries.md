## Activity Examples

### Daily request trends

Count the number of requests per day to analyze traffic trends over time.

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

### Top 10 frequently accessed URIs

Identify the most accessed URIs along with the top client IPs, grouped by the action taken (ALLOW, BLOCK, CHALLENGE, CAPTCHA).

```sql
select
  terminating_rule_id,
  http_request.clientIp as client_ip,
  http_request.uri as request_uri,
  count(*) as request_count
from
  aws_waf_traffic_log
group by
  terminating_rule_id,
  request_uri,
  client_ip
order by
  request_count desc
limit 10;
```

### Analyze CAPTCHA and CHALLENGE failures

Analyze the total requests and categorize failures based on CAPTCHA and CHALLENGE response reasons.

```sql
select
  http_request.clientIp as client_ip,
  count(*) as total_requests,

  -- Count of CHALLENGE & CAPTCHA actions
  sum(case when action = 'CHALLENGE' then 1 else 0 end) as challenge_count,
  sum(case when action = 'CAPTCHA' then 1 else 0 end) as captcha_count,

  -- CAPTCHA & CHALLENGE Failure Reasons
  sum(case when captcha_response.failureReason = 'TOKEN_INVALID' then 1 else 0 end) as challenge_token_invalid,
  sum(case when captcha_response.failureReason = 'TOKEN_INVALID' then 1 else 0 end) as captcha_token_invalid,
  sum(case when captcha_response.failureReason = 'TOKEN_DOMAIN_MISMATCH' then 1 else 0 end) as challenge_token_domain_mismatch,
  sum(case when captcha_response.failureReason = 'TOKEN_DOMAIN_MISMATCH' then 1 else 0 end) as captcha_token_domain_mismatch,
  sum(case when captcha_response.failureReason = 'TOKEN_EXPIRED' then 1 else 0 end) as challenge_token_expired,
  sum(case when captcha_response.failureReason = 'TOKEN_EXPIRED' then 1 else 0 end) as captcha_token_expired,
  sum(case when captcha_response.failureReason = 'TOKEN_MISSING' then 1 else 0 end) as challenge_token_missing,
  sum(case when captcha_response.failureReason = 'TOKEN_MISSING' then 1 else 0 end) as captcha_token_missing
from 
  aws_waf_traffic_log
group by 
  client_ip
order by 
  total_requests desc;
```

## Operational Examples

### Retrieve terminating rule matched data for requests
Extract details of requests that were terminated due to rule matches, helping analyze why a request was blocked or allowed.

```sql
with terminating_rule_match_details as (
  select
    timestamp,
    http_request.clientIp as client_ip,
    http_request.uri as request_uri,
    action,
    unnest(from_json(terminating_rule_match_details, '["JSON"]')) as match_details
  from
    aws_waf_traffic_log
  where
    json_array_length(terminating_rule_match_details) > 0
)
select
  timestamp,
  client_ip,
  request_uri,
  action,
  match_details ->> 'conditionType' as condition_type,
  match_details ->> 'sensitivityLevel' as sensitivity_level,
  match_details ->> 'location' as location,
  match_details ->> 'matchedData' as matched_data
from
  terminating_rule_match_details;
```

### Identify requests missing labels

Retrieve requests that do not contain labels, which help in categorizing and identifying the reason for request handling.

```sql
select
  timestamp,
  action,
  http_request.clientIp as client_ip,
  request_headers_inserted
from
  aws_waf_traffic_log
where
  labels is null
order by
  timestamp;
```

### Analyze blocked requests by rule

Retrieve the number of requests blocked by each rule, providing insights into which rules are most frequently triggered.

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

### Identify the most targeted URLs

Find which URLs or endpoints are most frequently accessed, helping detect high-risk areas in the application.

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

### Detect requests with CAPTCHA failures

Retrieve requests where CAPTCHA validation failed, indicating unsuccessful user verification.

```sql
select
  timestamp,
  http_request.clientIp as client_ip,
  http_request.uri as request_uri,
  action,
  captcha_response.responseCode as response_code,
  captcha_response.solveTimestamp as solve_timestamp,
  captcha_response.failureReason as failure_reason
from
  aws_waf_traffic_log
where
  response_code > 0
  and failure_reason is not null;
```

### Detect high volume of blocked requests

Identify IP addresses generating a high volume of blocked requests.

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

### Detect requests triggering multiple rules

Find requests that matched multiple non-terminating rules within a single evaluation.

```sql
select
  timestamp,
  http_request.clientIp as client_ip,
  http_source_name,
  http_source_id,
  json_array_length(non_terminating_matching_rules) as matched_rules
from
  aws_waf_traffic_log
where
  json_array_length(non_terminating_matching_rules) > 1
order by
  matched_rules desc;
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

### Analyze high volume of blocked requests

Identify patterns in blocked traffic over time.

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

### Identify the most frequently triggered rules

Analyze rules that are most frequently triggered to assess their effectiveness.

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

### Get non-terminating rules that detected SQL injection

Evaluate rule effectiveness in detecting SQL injection before enforcing stricter rules.

```sql
with not_terminating_rule as (
  select
    timestamp,
    http_request.clientIp as client_ip,
    action,
    unnest(from_json(non_terminating_matching_rules, '["JSON"]')) as rules
  from
    aws_waf_traffic_log
  where
    json_array_length(non_terminating_matching_rules) > 0
)
select
  timestamp,
  client_ip,
  action,
  (rules ->> 'ruleMatchDetails') as rule_match_details
from
  not_terminating_rule
where
  json_contains(rule_match_details, '{"conditionType": "SQL_INJECTION"}');
```