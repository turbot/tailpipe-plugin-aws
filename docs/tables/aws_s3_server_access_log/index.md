---
title: "Tailpipe Table: aws_s3_server_access_log - Query AWS S3 Server Access Logs"
description: "AWS S3 server access logs provide detailed information about requests made to your S3 buckets, including request source, operations performed, and response details."
---

# Table: aws_s3_server_access_log - Query AWS S3 server access logs

The `aws_s3_server_access_log` table allows you to query AWS S3 server access logs. This table captures detailed request and access information for S3 buckets, helping to analyze access patterns, troubleshoot issues, and enhance security monitoring.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_s3_server_access_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_s3_server_access_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_s3_server_access_log" "my_s3_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "s3-server-access-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_s3_server_access_log` partitions:

```sh
tailpipe collect aws_s3_server_access_log
```

Or for a single partition:

```sh
tailpipe collect aws_s3_server_access_log.my_s3_logs
```

## Query

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_s3_server_access_log)**

### Find all failed requests

```sql
select
  timestamp,
  bucket,
  request_uri,
  http_status,
  error_code,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  http_status is not null
  and http_status >= 400
order by
  timestamp desc;
```

### Identify top 10 users accessing a bucket

```sql
select
  requester,
  count(*) as request_count
from
  aws_s3_server_access_log
where
  bucket = 'test-tailpipe-source-pc'
group by
  requester
order by
  request_count desc
limit 10;
```

### Detect unusually large S3 downloads

```sql
select
  timestamp,
  bucket,
  key,
  bytes_sent,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  bytes_sent is not null
  and bytes_sent > 50000000 -- 50MB
order by
  bytes_sent desc;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect S3 server access logs stored in an S3 bucket.

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_s3_server_access_log" "my_s3_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "s3-server-access-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect S3 server access logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_s3_server_access_log" "my_s3_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "s3-server-access-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs using a prefix with year and month

Include the year and month in the prefix to speed up collection when using [non-date-based partition](https://docs.aws.amazon.com/AmazonS3/latest/userguide/ServerLogs.html#server-log-keyname-format).

```hcl
partition "aws_s3_server_access_log" "my_s3_logs_year_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "s3-server-access-logs-bucket"
    prefix     = "my/prefix/2025-01"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `(%{NUMBER:account_id}/%{DATA:region}/%{DATA:bucket_name}/%{YEAR:partition_year}/%{MONTHNUM:partition_month}/%{MONTHDAY:partition_day}/)?%{YEAR:year}-%{MONTHNUM:month}-%{MONTHDAY:day}-%{HOUR:hour}-%{MINUTE:minute}-%{SECOND:second}-%{DATA:suffix}` |
