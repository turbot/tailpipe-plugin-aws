---
title: "Tailpipe Table: aws_cloudtrail_log - Query AWS CloudTrail Logs"
description: "AWS CloudTrail logs capture API activity and user actions within your AWS account."
---

# Table: aws_cloudtrail_log - Query AWS CloudTrail Logs

The `aws_cloudtrail_log` table allows you to query data from [AWS CloudTrail logs](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-working-with-log-files.html). This table provides detailed information about API calls made within your AWS account, including the event name, source IP address, user identity, and more.

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

## Query

**[Explore 100+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_cloudtrail_log)**

### Root Activity

Find any actions taken by the root user.

```sql
select
  event_time,
  event_name,
  source_ip_address,
  user_agent,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  user_identity.type = 'Root'
order by
  event_time desc;
```

### Top 10 Events

List the top 10 events and how many times they were called.

```sql
select
  event_source,
  event_name,
  count(*) as event_count
from
  aws_cloudtrail_log
group by
  event_source,
  event_name,
order by
  event_count desc
limit 10;
```

### High Volume S3 Access Requests

Find users generating a high volume of S3 access requests to identify potential anomalous activity.

```sql
select
  user_identity.arn as user_arn,
  count(*) as event_count,
  date_trunc('minute', event_time) as event_minute
from
  aws_cloudtrail_log
where
  event_source = 's3.amazonaws.com'
  and event_name in ('GetObject', 'ListBucket')
group by
  user_identity.arn,
  event_minute
having
  count(*) > 100
order by
  event_count desc;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect CloudTrail logs stored in an S3 bucket that use the [default log file name format](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/get-and-view-cloudtrail-log-files.html).

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

### Collect logs from an S3 bucket with a prefix

Collect CloudTrail logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_cloudtrail_log" "my_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-cloudtrail-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Exclude read-only events

Use the filter argument in your partition to exclude read-only events and reduce the size of local log storage.

```hcl
partition "aws_cloudtrail_log" "my_logs_write" {
  # Avoid saving read-only events, which can drastically reduce local log size
  filter = "not read_only"

  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-cloudtrail-logs-bucket"
  }
}
```

### Collect logs for all accounts in an organization

For a specific organization, collect logs for all accounts and regions.

```hcl
partition "aws_cloudtrail_log" "my_logs_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = `AWSLogs/o-aa111bb222/%{NUMBER:account_id}/CloudTrail/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz`
  }
}
```

### Collect logs for a single account

For a specific account, collect logs for all regions.

```hcl
partition "aws_cloudtrail_log" "my_logs_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = `AWSLogs/(%{DATA:org_id}/)?123456789012/CloudTrail/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz`
  }
}
```

### Collect logs for a single region

For all accounts, collect logs from us-east-1.

```hcl
partition "aws_cloudtrail_log" "my_logs_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz`
  }
}
```

### Collect logs for multiple regions

For all accounts, collect logs from us-east-1 and us-east-2.

```hcl
partition "aws_cloudtrail_log" "my_logs_regions" {
  source "aws_s3_bucket"  {
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/(us-east-1|us-east-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz`
  }
}
```

### Collect logs from a CloudWatch log group

Collect CloudTrail logs from all log streams in a CloudWatch log group.

```hcl
partition "aws_cloudtrail_log" "cw_log_group_logs" {
  source "aws_cloudwatch_log_group" {
    connection     = connection.aws.logging_account
    log_group_name = "aws-cloudtrail-logs-123456789012-fd33b044"
    region         = "us-east-1"
  }
}
```

### Collect logs from a CloudWatch log group for a specific account and region

Collect CloudTrail logs for a single region in an account.

```hcl
partition "aws_cloudtrail_log" "cw_log_group_logs_specific" {
  source "aws_cloudwatch_log_group" {
    connection       = connection.aws.default
    log_group_name   = "aws-cloudtrail-logs-123456789012-fd33b044"
    log_stream_names = ["456789012345_CloudTrail_us-east-1*"]
    region           = "us-east-1"
  }
}
```

### Collect logs from a CloudWatch log group for all regions in an account

Collect CloudTrail logs for all regions in an account.

```hcl
partition "aws_cloudtrail_log" "cw_log_group_logs_all_regions" {
  source "aws_cloudwatch_log_group" {
    connection       = connection.aws.default
    log_group_name   = "aws-cloudtrail-logs-123456789012-fd33b044"
    log_stream_names = ["456789012345_CloudTrail_*"]
    region           = "us-east-1"
  }
}
```

### Collect logs from local files

You can also collect CloudTrail logs from local files, like the [flaws.cloud public dataset](https://summitroute.com/blog/2020/10/09/public_dataset_of_cloudtrail_logs_from_flaws_cloud/).

```hcl
partition "aws_cloudtrail_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/cloudtrail_logs"]
    file_layout = `%{DATA}.json.gz`
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

AWS does not offer a default or recommended S3 bucket file path for Security Hub Findings export. The default file layout used in this table is based on the implementation provided in the [AWS Security Hub Findings Export Samples](https://github.com/aws-samples/aws-security-hub-findings-export) repository.

| Argument    | Default                                                                                                                                   |
| ----------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| file_layout | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz` |
