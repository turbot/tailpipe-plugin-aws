---
title: "Tailpipe Table: aws_waf_traffic_log - Query AWS WAF Traffic Logs"
description: "AWS WAF Traffic Logs capture detailed information about web requests inspected by AWS WAF, helping analyze threats, monitor rule effectiveness, and improve security posture."
---

# Table: aws_waf_traffic_log - Query AWS WAF Traffic Logs

The `aws_waf_traffic_log` table allows you to query data from AWS WAF traffic logs. This table provides detailed insights into incoming web requests, including the request source, matched WAF rules, rule actions, and threat indicators. Use this data to monitor traffic patterns, detect anomalies, and fine-tune WAF rules.

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
    bucket     = "aws-waf-traffic-logs-bucket"
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

**[Explore 11+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_waf_traffic_log)**

### Blocked Requests by WAF

Find all blocked requests recorded by AWS WAF.

```sql
select
  tp_timestamp,
  (http_request ->> 'client_ip') as client_ip,
  (http_request ->> 'country') as country,
  rule_group.name as rule_matched,
  action
from
  aws_waf_traffic_log
where
  action = 'BLOCK'
order by
  tp_timestamp desc;
```

### Top Sources of WAF-Blocked Traffic

Identify IPs frequently blocked by AWS WAF.

```sql
select
  (http_request ->> 'clientIp') as client_ip,
  count(*) as block_count
from
  aws_waf_traffic_log
where
  action = 'BLOCK'
group by
  (http_request ->> 'clientIp')
order by
  block_count desc
limit 10;
```

### Requests Matching SQL Injection Rule

Find web requests that matched AWS WAF’s SQL Injection detection.

```sql
select
  timestamp,
  (http_request ->> 'uri') as request_uri,
  (terminating_rule_match_details ->> 'conditionType') as condition_type,
  (http_request ->> 'clientIp') as client_ip,
  action
from
  aws_waf_traffic_log
where
  (rule_match_details ->> 'condition_type') = 'SQL_INJECTION'
order by
  timestamp desc;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect WAF logs stored in an S3 bucket that uses the default log file format.

```hcl
connection "aws" "security_account" {
  profile = "my-security-account"
}

partition "aws_waf_traffic_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-waf-traffic-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect WAF logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_waf_traffic_log" "my_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-waf-traffic-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from local files

You can also collect AWS WAF logs from local files.

```hcl
partition "aws_waf_traffic_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/aws_waf_traffic_logs"]
    file_layout = "%{DATA}.json.gz"
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
    file_layout = "AWSLogs/o-aa111bb222/%{NUMBER:account_id}/WAFLogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.json.gz"
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
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/WAFLogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.json.gz"
  }
}
```

### Collect logs for a single region

For all accounts, collect logs from `us-east-1`.

```hcl
partition "aws_waf_traffic_log" "my_logs_region" {
  source "aws_s3_bucket" {
    connection  = connection.aws.security_account
    bucket      = "waf-traffic-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/WAFLogs/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.json.gz"
  }
}
```

### Collect logs for multiple regions

For all accounts, collect logs from `us-east-1` and `us-west-2`.

```hcl
partition "aws_waf_traffic_log" "my_logs_regions" {
  source "aws_s3_bucket" {
    bucket      = "waf-traffic-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/WAFLogs/(us-east-1|us-west-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.json.gz"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default                                                                                                                                                                                |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| file_layout | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/WAFLogs/%{DATA:log_group_name}/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{MINUTE:minute}/%{DATA}.gz` |
