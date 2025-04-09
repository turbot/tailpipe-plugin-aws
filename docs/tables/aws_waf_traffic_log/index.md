---
title: "Tailpipe Table: aws_waf_traffic_log - Query AWS WAF Traffic Logs"
description: "AWS WAF traffic logs capture detailed information about web requests inspected by AWS WAF, helping analyze threats, monitor rule effectiveness, and improve security posture."
---

# Table: aws_waf_traffic_log - Query AWS WAF Traffic Logs

The `aws_waf_traffic_log` table allows you to query data from [AWS WAF traffic logs](https://docs.aws.amazon.com/waf/latest/developerguide/logging.html). This table provides detailed insights into incoming web requests, including the request source, matched WAF rules, rule actions, and threat indicators. Use this data to monitor traffic patterns, detect anomalies, and fine-tune WAF rules.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_waf_traffic_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_waf_traffic_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "security_account" {
  profile = "my-security-account"
}

partition "aws_waf_traffic_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-waf-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_waf_traffic_log` partitions:

```sh
tailpipe collect aws_waf_traffic_log
```

Or for a single partition:

```sh
tailpipe collect aws_waf_traffic_log.my_logs
```

## Query

**[Explore 18+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_waf_traffic_log)**

### Blocked Requests

Find all blocked requests recorded by AWS WAF.

```sql
select
  tp_timestamp,
  http_request.clientIp as client_ip,
  http_request.country as country,
  action
from
  aws_waf_traffic_log
where
  action = 'BLOCK'
order by
  tp_timestamp desc;
```

### Top 10 Blocked Client IPs

List the top 10 IPs frequently blocked by AWS WAF.

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

### Blocked Requests With SQL Injection

Find web requests that matched AWS WAF's SQL injection detection.

```sql
select
  timestamp,
  http_request.uri as request_uri,
  http_request.clientIp as client_ip,
  action,
  terminating_rule
from
  aws_waf_traffic_log,
  unnest(
    from_json(terminating_rule_match_details, '["JSON"]')
  ) as terminating_rule
where
  json_contains(terminating_rule, '{"conditionType":"SQL_INJECTION"}')
order by
  timestamp desc;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect WAF traffic logs stored in an S3 bucket that uses the default log file format.

```hcl
connection "aws" "security_account" {
  profile = "my-security-account"
}

partition "aws_waf_traffic_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-waf-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_waf_traffic_log" "my_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-waf-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from a CloudWatch log group

Collect WAF traffic logs from a CloudWatch log group without any log stream filtering.

```hcl
partition "aws_waf_traffic_log" "cw_log_group_logs" {
  source "aws_cloudwatch_log_group" {
    connection     = connection.aws.logging_account
    log_group_name = "aws-waf-log-testLogGroup2"
    region         = "us-east-1"
  }
}
```

### Collect logs from a CloudWatch log group with a log stream

Collect WAF traffic logs for a specific Web ACL in us-east-1 by using an exact log stream name match.

```hcl
partition "aws_waf_traffic_log" "cw_log_group_logs_specific" {
  source "aws_cloudwatch_log_group" {
    connection = connection.aws.logging_account
    log_group_name = "aws-waf-log-testLogGroup2"
    log_stream_names = ["us-east-1_TestWebACL_123456"]
    region = "us-east-1"
  }
}
```

### Collect WAF logs from CloudWatch for all Web ACLs in a region

Collect WAF traffic logs for all Web ACLs in us-east-1 by using a wildcard pattern in the log stream name.

```hcl
partition "aws_waf_traffic_log" "cw_log_group_logs_all" {
  source "aws_cloudwatch_log_group" {
    connection = connection.aws.logging_account
    log_group_name = "aws-waf-log-testLogGroup2"
    log_stream_names = ["us-east-1_*"]
    region = "us-east-1"
  }
}
```

### Collect logs from local files

You can also collect logs from local files.

```hcl
partition "aws_waf_traffic_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/aws_waf_traffic_logs"]
    file_layout = `%{DATA}.json.gz`
  }
}
```

### Collect logs for all WAF ACLs in an organization

For a specific organization, collect logs for all WAF ACLs.

```hcl
partition "aws_waf_traffic_log" "my_logs_org" {
  source "aws_s3_bucket" {
    connection  = connection.aws.security_account
    bucket      = "waf-traffic-logs-bucket"
    file_layout = `AWSLogs/o-aa111bb222/%{NUMBER:account_id}/WAFLogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.json.gz`
  }
}
```

### Collect logs for a single account

For a specific account, collect logs for all regions.

```hcl
partition "aws_waf_traffic_log" "my_logs_account" {
  source "aws_s3_bucket" {
    connection  = connection.aws.security_account
    bucket      = "waf-traffic-logs-bucket"
    file_layout = `AWSLogs/(%{DATA:org_id}/)?123456789012/WAFLogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.json.gz`
  }
}
```

### Collect logs from local files

You can also collect logs from local files.

```hcl
partition "aws_waf_traffic_log" "my_logs" {
  source "file"  {
    paths       = ["/Users/myuser/aws_waf_traffic_log"]
    file_layout = `%{DATA}.txt`
  }
}
```

### Exclude GET requests

Use the filter argument in your partition to exclude read-only requests and reduce the size of local log storage.

```hcl
partition "aws_waf_traffic_log" "my_logs_write" {
  filter = "(http_request ->> 'httpMethod') not like 'GET'"

  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "waf-traffic-logs-bucket"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default                                                                                                                                                                                                                |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| file_layout | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/WAFLogs/%{DATA:cloudfront_or_region}/%{DATA:cloudfront_name_or_resource_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.gz` |
