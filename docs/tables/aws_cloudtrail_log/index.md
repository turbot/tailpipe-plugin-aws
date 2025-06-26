---
title: "Tailpipe Table: aws_cloudtrail_log - Query AWS CloudTrail Logs"
description: "AWS CloudTrail logs record detailed information about API calls and resource changes in your AWS account, helping track user activity, security analysis, and compliance auditing."
---

# Table: aws_cloudtrail_log

AWS CloudTrail logs record detailed information about API calls and resource changes in your AWS account. These logs are essential for security analysis, resource change tracking, and compliance auditing.

## Examples

### Basic log analysis
```sql
select
  event_time,
  event_name,
  event_source,
  source_ip_address,
  user_identity->>'userName' as user_name,
  error_code
from
  aws_cloudtrail_log
where
  tp_date >= current_date - interval '7 days'
order by
  event_time desc;
```

### Find unauthorized API calls
```sql
select
  event_time,
  event_name,
  event_source,
  source_ip_address,
  user_identity->>'userName' as user_name,
  error_code,
  error_message
from
  aws_cloudtrail_log
where
  error_code like '%Unauthorized%'
  or error_code like '%AccessDenied%'
order by
  event_time desc;
```

### Track root account usage
```sql
select
  event_time,
  event_name,
  event_source,
  source_ip_address,
  user_identity->>'type' as identity_type
from
  aws_cloudtrail_log
where
  user_identity->>'type' = 'Root'
order by
  event_time desc;
```

### Monitor IAM policy changes
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'policyName' as policy_name,
  request_parameters->>'policyDocument' as policy_document
from
  aws_cloudtrail_log
where
  event_source = 'iam.amazonaws.com'
  and event_name like '%Policy%'
order by
  event_time desc;
```

### List resource deletions
```sql
select
  event_time,
  event_name,
  event_source,
  user_identity->>'userName' as user_name,
  resources[0]->>'resourceType' as resource_type,
  resources[0]->>'resourceName' as resource_name
from
  aws_cloudtrail_log
where
  event_name like 'Delete%'
order by
  event_time desc;
```

### Find API calls from specific IP addresses
```sql
select
  event_time,
  event_name,
  event_source,
  source_ip_address,
  user_identity->>'userName' as user_name
from
  aws_cloudtrail_log
where
  source_ip_address = '192.0.2.1'
order by
  event_time desc;
```

### Track changes to security groups
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'groupId' as security_group_id,
  request_parameters->>'ipPermissions' as ip_permissions
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name like '%SecurityGroup%'
order by
  event_time desc;
```

## Source Configuration

### S3 Source

CloudTrail logs can be read from S3 buckets. The default file layout pattern is:

```
AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz
```

### CloudWatch Logs Source

CloudTrail logs can also be read from CloudWatch Log Groups where CloudTrail is configured to deliver logs.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_cloudtrail_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cloudtrail_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-cloudtrail-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_cloudtrail_log` partitions:

```sh
tailpipe collect aws_cloudtrail_log
```

Or for a single partition:

```sh
tailpipe collect aws_cloudtrail_log.my_logs
```

## Query Examples

### List IAM user creation events

Find all events where IAM users were created:

```sql
select
  event_time,
  event_source,
  event_name,
  user_identity->>'userName' as actor,
  request_parameters->>'userName' as created_user,
  source_ip_address,
  aws_region
from
  aws_cloudtrail_log
where
  event_name = 'CreateUser'
  and event_source = 'iam.amazonaws.com'
order by
  event_time desc;
```

### Failed API calls

Find failed API calls to investigate potential security issues or misconfigurations:

```sql
select
  event_time,
  event_source,
  event_name,
  error_code,
  error_message,
  source_ip_address,
  user_identity->>'userName' as user_name,
  aws_region
from
  aws_cloudtrail_log
where
  error_code is not null
order by
  event_time desc;
```

### Security group changes

Monitor changes to security groups:

```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'groupId' as security_group_id,
  source_ip_address,
  aws_region
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name like '%SecurityGroup%'
order by
  event_time desc;
```

### Root account activity

Monitor AWS root account usage:

```sql
select
  event_time,
  event_source,
  event_name,
  source_ip_address,
  aws_region,
  user_agent
from
  aws_cloudtrail_log
where
  user_identity->>'type' = 'Root'
order by
  event_time desc;
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default |
|------------|---------|
| file_layout | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz` |
``` 