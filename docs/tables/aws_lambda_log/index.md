---
title: "Tailpipe Table: aws_lambda_log - Query AWS Lambda Logs"
description: "AWS Lambda logs capture invocation details and function output within your AWS account."
---

# Table: aws_lambda_log - Query AWS Lambda Logs

The `aws_lambda_log` table allows you to query data from [AWS Lambda logs](https://docs.aws.amazon.com/lambda/latest/dg/monitoring-cloudwatchlogs.html). This table provides detailed information about Lambda function invocations, including request ID, log level, message content, timestamps, and more.

## Message Format and Parsing

The `aws_lambda_log` table provides multiple message fields to handle different log formats across Lambda runtimes:

### Message Fields

- **`message`** (string) – Extracted and parsed message content when possible. Always populated if a message can be extracted, including from JSON-formatted logs.
- **`message_json`** (json) – Extracted message parsed as JSON. Populated only if the extracted message is valid JSON and can be converted.
- **`raw_message`** (string) – Complete original message string as received from the log source. Always populated regardless of format.
- **`raw_message_json`** (json) – Full original message parsed as JSON. Populated only if the original message is native JSON or convertible to JSON.

### Runtime Behavior

The table handles logs consistently across most AWS Lambda runtimes (Node.js, Python, .NET, Go, Ruby), with automatic parsing of both plain text and JSON-formatted logs.

**Note:** PowerShell runtime emits logs in a different format compared to other runtimes, which may affect message extraction and parsing. Refer to the [AWS Lambda runtime documentation](https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html) for runtime-specific log format details.

### JSON Log Handling

- JSON-formatted logs are automatically parsed and stored in `raw_message_json`
- If the extracted message portion is also valid JSON, it's additionally parsed into `message_json`
- This dual approach ensures you can query both structured JSON data and extract specific message content

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_lambda_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_lambda_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_lambda_log" "my_logs" {
  source "aws_cloudwatch_log_group" {
    connection     = connection.aws.logging_account
    log_group_name = "/aws/lambda/my-function"
    region         = "us-east-1"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_lambda_log` partitions:

```sh
tailpipe collect aws_lambda_log
```

Or for a single partition:

```sh
tailpipe collect aws_lambda_log.my_logs
```

## Query

**[Explore 16+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_lambda_log)**

### Recent Error Messages

Find recent error messages from Lambda functions.

```sql
select
  timestamp,
  tp_source_name as function_name,
  log_level,
  message
from
  aws_lambda_log
where
  log_level = 'ERROR'
order by
  timestamp desc
limit 100;
```

### Slow Function Executions

Identify Lambda functions with long execution times.

```sql
select
  tp_source_name as function_name,
  request_id,
  duration * 1000 as duration_ms,
  timestamp
from
  aws_lambda_log
where
  duration > 5
order by
  duration desc
limit 20;
```

### Memory Utilization

Find Lambda functions approaching their memory limits.

```sql
select
  tp_source_name as function_name,
  request_id,
  max_memory_used as memory_used_mb,
  memory_size as memory_limit_mb,
  round((max_memory_used::float / memory_size::float) * 100, 2) as memory_utilization_percent,
  timestamp
from
  aws_lambda_log
where
  max_memory_used is not null
  and memory_size is not null
  and (max_memory_used::float / memory_size::float) > 0.8
order by
  memory_utilization_percent desc;
```

## Example Configurations

### Collect all log streams explicitly

Collect logs from all streams in a log group by explicitly setting the wildcard pattern.

```hcl
partition "aws_lambda_log" "all_streams" {
  source "aws_cloudwatch_log_group" {
    connection      = connection.aws.logging_account
    log_group_name  = "/aws/lambda/my-function"
    log_stream_names = ["*"]
    region          = "us-east-1"
  }
}
```

### Filter logs by log stream prefix

Collect Lambda logs from streams with a specific prefix.

```hcl
partition "aws_lambda_log" "prefix_filtered_logs" {
  source "aws_cloudwatch_log_group" {
    connection      = connection.aws.logging_account
    log_group_name  = "/aws/lambda/my-function"
    log_stream_names = ["PROD_*", "2023/07/*"]
    region          = "us-east-1"
  }
}
```

### Collect logs from an S3 bucket

Collect Lambda logs archived to an S3 bucket.

```hcl
partition "aws_lambda_log" "s3_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "lambda-logs-bucket"
    prefix     = "lambda-logs"
  }
}
```

### Collect logs from local files

You can also collect Lambda logs from local files.

```hcl
partition "aws_lambda_log" "local_logs" {
  source "file" {
    paths       = ["/Users/myuser/lambda_logs"]
    file_layout = `%{DATA}.log`
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default                                                         |
| ----------- | --------------------------------------------------------------- |
| file_layout | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/%{DATA:region}/%{DATA:function_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{DATA}.log.zst` |
