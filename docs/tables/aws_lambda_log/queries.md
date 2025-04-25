## Activity Examples

### Recent Lambda Log Activity

This query shows the 100 most recent Lambda log entries across all functions. Real-time monitoring of Lambda activity helps with troubleshooting issues and understanding the current state of your serverless applications.

```sql
select
  tp_timestamp,
  request_id,
  log_type,
  log_level,
  message,
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name
from
  aws_lambda_log
order by
  tp_timestamp desc
limit
  100;
```

```yaml
folder: Lambda
```

### Lambda Execution Trends by Hour

This query shows Lambda execution trends by hour for each function. Understanding these patterns helps with capacity planning, identifying unusual activity spikes, and optimizing resources based on time-of-day usage patterns.

```sql
select
  date_trunc('hour', tp_timestamp) as hour,
  count(*) as execution_count,
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name
from
  aws_lambda_log
where
  log_type = 'START'
group by
  hour,
  tp_source_name
order by
  hour desc,
  execution_count desc;
```

```yaml
folder: Lambda
```

### Application Log Level Distribution

This query analyzes the distribution of log levels in Lambda application logs. Reviewing log level patterns helps identify functions generating excessive logs or experiencing frequent errors that may impact application reliability.

```sql
select
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name,
  log_level,
  count(*) as log_count
from
  aws_lambda_log
where
  log_level is not null
group by
  tp_source_name,
  log_level
order by
  lambda_function_name,
  log_count desc;
```

```yaml
folder: Lambda
```

### Execution Flow for a Specific Request

This query shows the complete execution flow for a specific Lambda request. Tracing the sequence of logs for a single invocation helps debug issues and understand function behavior from start to finish.

```sql
select
  tp_timestamp,
  log_type,
  log_level,
  substring(message, 1, 200) as message_preview,
  case
    when log_type in ('START', 'END', 'REPORT', 'INIT_START') then 'System Log'
    when log_level is not null then 'Application Log'
    else 'Other'
  end as log_category,
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name
from
  aws_lambda_log
where
  request_id = '9286dcef-4fac-4706-99f6-0f0087763dbc'
order by
  tp_timestamp asc;
```

```yaml
folder: Lambda
```

## Detection Examples

### Lambda Function Error Analysis

This query finds the 100 most recent Lambda function errors and timeouts. Monitoring these errors helps identify reliability issues and functions that need error handling improvements for better application stability.

```sql
select
  tp_timestamp,
  request_id,
  message,
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name
from
  aws_lambda_log
where
  log_level = 'ERROR'
  or message like '%Task timed out%'
  or message like '%Memory Size exceeded%'
  or message like '%Process exited before completing request%'
order by
  tp_timestamp desc
limit
  100;
```

```yaml
folder: Lambda
```

### Functions with High Billing-to-Execution Time Ratio

This query identifies functions with high billing-to-execution time ratios. Optimizing these functions can reduce costs by addressing the gap between actual runtime and billed duration, especially for functions with significant billing overhead.

```sql
select
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name,
  round(avg(duration), 2) as avg_duration_ms,
  round(avg(billed_duration), 2) as avg_billed_duration_ms,
  round(avg(billed_duration - duration), 2) as avg_billing_overhead_ms,
  round(avg(billed_duration * 100.0 / duration) - 100, 2) as billing_overhead_percent,
  count(*) as execution_count
from
  aws_lambda_log
where
  log_type = 'REPORT'
  and duration > 0
  and billed_duration is not null
group by
  tp_source_name
having
  avg(billed_duration - duration) > 10
order by
  billing_overhead_percent desc;
```

```yaml
folder: Lambda
```

### Lambda Cold Start Analysis

This query analyzes Lambda cold starts by counting initialization events for each function. Identifying functions with frequent cold starts helps prioritize optimization efforts to reduce latency and improve user experience.

```sql
select
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name,
  count(distinct request_id) as cold_start_count,
  avg(duration) as avg_init_duration_ms
from
  aws_lambda_log
where
  log_type = 'INIT_START'
  or message like '%Init Duration:%'
group by
  tp_source_name
order by
  cold_start_count desc;
```

```yaml
folder: Lambda
```

## Operational Examples

### Top 10 Slowest Lambda Function Executions

This query identifies the 10 slowest Lambda function executions by examining REPORT logs. Finding these slow executions helps pinpoint specific instances that require optimization to improve overall performance.

```sql
select
  tp_timestamp,
  request_id,
  duration as duration_ms,
  billed_duration as billed_duration_ms,
  memory_size as allocated_memory_mb,
  max_memory_used as used_memory_mb,
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name
from
  aws_lambda_log
where
  log_type = 'REPORT'
  and duration is not null
order by
  duration desc
limit
  10;
```

```yaml
folder: Lambda
```

### Memory Utilization Efficiency

This query calculates memory utilization efficiency for each Lambda function. Finding the right memory allocation helps optimize costs while maintaining performance by identifying over-provisioned or under-provisioned functions.

```sql
select
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name,
  memory_size as allocated_memory_mb,
  round(avg(max_memory_used), 2) as avg_used_memory_mb,
  round(avg(max_memory_used * 100.0 / memory_size), 2) as memory_utilization_percent,
  count(*) as execution_count
from
  aws_lambda_log
where
  log_type = 'REPORT'
  and memory_size is not null
  and max_memory_used is not null
group by
  tp_source_name,
  memory_size
order by
  memory_utilization_percent desc;
```

```yaml
folder: Lambda
```

### Detailed Request Execution Analysis

This query analyzes each request's complete execution lifecycle including message patterns and timing between execution phases. This helps identify bottlenecks in function execution flow and understand the context of each request through message analysis.

```sql
with request_phases as (
  select
    request_id,
    tp_source_name,
    regexp_replace(tp_source_name, '^/aws/lambda/', '') as function_name,
    min(case when log_type = 'START' then tp_timestamp end) as start_time,
    min(case when log_type = 'END' then tp_timestamp end) as end_time,
    min(case when log_type = 'REPORT' then tp_timestamp end) as report_time,
    max(case when log_type = 'REPORT' then duration end) as duration_ms,
    max(case when log_type = 'REPORT' then billed_duration end) as billed_duration_ms,
    max(case when log_type = 'REPORT' then memory_size end) as allocated_memory_mb,
    max(case when log_type = 'REPORT' then max_memory_used end) as max_memory_used_mb,
    count(case when log_level = 'INFO' then 1 end) as info_log_count,
    count(case when log_level = 'ERROR' then 1 end) as error_log_count,
    count(case when log_level = 'WARN' then 1 end) as warn_log_count,
    count(case when log_level = 'DEBUG' then 1 end) as debug_log_count,
    bool_or(message like '%Task timed out%') as has_timeout,
    bool_or(message like '%Init Duration:%') as has_cold_start,
    bool_or(message like '%Memory Size:%' and message like '%Max Memory Used:%') as has_memory_metrics,
    bool_or(message like '%error%' or message like '%exception%' or message like '%fail%') as has_error_keywords
  from
    aws_lambda_log
  where
    request_id is not null
  group by
    request_id,
    tp_source_name,
    function_name
),
message_samples as (
  select
    request_id,
    array_agg(message) as message_samples
  from
    (select
      request_id,
      message,
      row_number() over (partition by request_id order by tp_timestamp) as rn
    from
      aws_lambda_log
    where
      request_id is not null
      and log_level is not null
    order by
      tp_timestamp) t
  where
    rn <= 3
  group by
    request_id
)
select
  rp.request_id,
  rp.start_time,
  rp.end_time,
  extract(epoch from (rp.end_time - rp.start_time)) * 1000 as total_execution_time_ms,
  rp.duration_ms as reported_duration_ms,
  rp.billed_duration_ms,
  rp.allocated_memory_mb,
  rp.max_memory_used_mb,
  rp.info_log_count,
  rp.error_log_count,
  rp.warn_log_count,
  rp.debug_log_count,
  rp.has_timeout,
  rp.has_cold_start,
  rp.has_memory_metrics,
  rp.has_error_keywords,
  ms.message_samples[1] as first_message,
  rp.function_name as lambda_function_name
from
  request_phases rp
left join
  message_samples ms on rp.request_id = ms.request_id
where
  rp.start_time is not null
  and rp.end_time is not null
order by
  rp.start_time desc
limit
  50;
```

```yaml
folder: Lambda
```

## Volume Examples

### Lambda Function Execution Summary

This query summarizes execution metrics for each Lambda function. It helps identify frequently invoked functions and their performance characteristics for cost optimization and performance tuning.

```sql
select
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name,
  count(*) as execution_count,
  avg(duration) as avg_duration_ms,
  min(duration) as min_duration_ms,
  max(duration) as max_duration_ms,
  avg(billed_duration) as avg_billed_duration_ms,
  avg(max_memory_used) as avg_memory_used_mb,
  max(memory_size) as allocated_memory_mb
from
  aws_lambda_log
where
  log_type = 'REPORT'
group by
  tp_source_name
order by
  execution_count desc;
```

```yaml
folder: Lambda
```

## Baseline Examples

### Lambda Duration Distribution

This query categorizes Lambda function executions into duration ranges. Understanding execution time distribution helps identify inconsistent performance patterns and optimize functions that show high variability in runtime.

```sql
select
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name,
  case
    when duration < 100 then '< 100ms'
    when duration < 500 then '100-500ms'
    when duration < 1000 then '500ms-1s'
    when duration < 3000 then '1s-3s'
    when duration < 10000 then '3s-10s'
    else '> 10s'
  end as duration_range,
  count(*) as execution_count
from
  aws_lambda_log
where
  log_type = 'REPORT'
  and duration is not null
group by
  tp_source_name,
  duration_range
order by
  lambda_function_name,
  case
    when duration_range = '< 100ms' then 1
    when duration_range = '100-500ms' then 2
    when duration_range = '500ms-1s' then 3
    when duration_range = '1s-3s' then 4
    when duration_range = '3s-10s' then 5
    else 6
  end;
```

```yaml
folder: Lambda
```

### Most Common Error Messages by Function

This query identifies the most common error messages for each Lambda function. Finding recurring error patterns helps prioritize which issues to fix first and understand the reliability challenges affecting specific functions.

```sql
select
  regexp_replace(tp_source_name, '^/aws/lambda/', '') as lambda_function_name,
  regexp_replace(message, '^.*Error: ', '') as error_message_pattern,
  count(*) as occurrence_count
from
  aws_lambda_log
where
  log_level = 'ERROR'
  or message like '%Error:%'
  or message like '%Exception:%'
  or message like '%Task timed out%'
group by
  tp_source_name,
  error_message_pattern
order by
  lambda_function_name,
  occurrence_count desc;
```

```yaml
folder: Lambda
```

### Error Message Pattern Analysis

This query analyzes error patterns across all Lambda functions to identify systematic issues. Detecting common error signatures across multiple functions helps identify widespread problems that may require architectural changes rather than function-specific fixes.

```sql
with error_patterns as (
  select
    regexp_replace(message, '([a-f0-9]{8}(-[a-f0-9]{4}){3}-[a-f0-9]{12}|[\d\.]+|"[^"]*"|''[^'']*'')', 'X') as normalized_error,
    count(*) as pattern_count,
    count(distinct regexp_replace(tp_source_name, '^/aws/lambda/', '')) as affected_functions_count,
    array_agg(distinct regexp_replace(tp_source_name, '^/aws/lambda/', '')) as affected_functions
  from
    aws_lambda_log
  where
    log_level = 'ERROR'
    or message like '%Error:%'
    or message like '%Exception:%'
    or message like '%Task timed out%'
  group by
    normalized_error
  having
    count(*) > 1
)
select
  normalized_error as error_pattern,
  pattern_count as occurrence_count,
  affected_functions_count,
  affected_functions
from
  error_patterns
order by
  pattern_count desc,
  affected_functions_count desc
limit
  20;
```

```yaml
folder: Lambda
```
