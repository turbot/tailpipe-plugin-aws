---
title: "Tailpipe Table: aws_lambda_log - Query AWS Lambda Logs"
description: "AWS Lambda logs capture invocation details and function output within your AWS account."
---

# Table: aws_lambda_log - Query AWS Lambda Logs

The `aws_lambda_log` table allows you to query data from [AWS Lambda logs](https://docs.aws.amazon.com/lambda/latest/dg/monitoring-cloudwatchlogs.html). This table provides detailed information about Lambda function invocations, including request ID, log level, message content, timestamps, and more.

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
  function_name,
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
  function_name,
  request_id,
  duration_ms,
  timestamp
from
  aws_lambda_log
where
  duration_ms > 5000
order by
  duration_ms desc
limit 20;
```

### Memory Utilization

Find Lambda functions approaching their memory limits.

```sql
select
  function_name,
  request_id,
  memory_used_mb,
  memory_limit_mb,
  round((memory_used_mb::float / memory_limit_mb::float) * 100, 2) as memory_utilization_percent,
  timestamp
from
  aws_lambda_log
where
  memory_used_mb is not null
  and memory_limit_mb is not null
  and (memory_used_mb::float / memory_limit_mb::float) > 0.8
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

// TODO: Change/remove the file layout
This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default                                                         |
| ----------- | --------------------------------------------------------------- |
| file_layout | `%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log.gz` |

### aws_cloudwatch_log_group

This table sets the following defaults for the [aws_cloudwatch_log_group source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_cloudwatch_log_group#arguments):

| Argument         | Default |
| ---------------- | ------- |
| log_stream_names | `["*"]` |
