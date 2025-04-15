---
title: "Tailpipe Table: aws_vpc_flow_log - Query AWS VPC Flow Logs"
description: "AWS VPC flow logs capture information about IP traffic going to and from network interfaces in your VPC."
---

# Table: aws_vpc_flow_log - Query AWS VPC Flow Logs

The `aws_vpc_flow_log` table allows you to query data from [AWS VPC Flow Logs](https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html). This table provides detailed insights into network traffic within your VPC, including source and destination IP addresses, ports, protocols, and more.

**Note**:
- Using `aws_s3_bucket` source:
  - The `format` block is **optional**. If not provided, the log format will be inferred from the **first line of the object** in the bucket.
  - However, when logs are **exported from CloudWatch to S3**, they typically **do not include a header line**, so the `format` block becomes **required**. If not specified, a **default format** will be applied.
- Using`aws_cloudwatch_log_group` source:
  - The `format` block is **required** to define the structure of the logs. If omitted, the system will fall back to a **default format**.
- For timestamp information, the `start` field will be used first, with the `end` field as a fallback. If neither field is available, then that log line will not be collected and Tailpipe will return an error.
- When determining each log's index, the table uses the following order of precedence:
  - `interface_id`
  - `subnet_id`
  - `vpc_id`
  - If neither field is present, the log will use `default` as the index instead of an AWS account ID.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_vpc_flow_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_vpc_flow_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "vpc_logging" {
  profile = "my-vpc-logging"
}

partition "aws_vpc_flow_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.vpc_logging
    bucket     = "aws-vpc-flow-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_vpc_flow_log` partitions:

```sh
tailpipe collect aws_vpc_flow_log
```

Or for a single partition:

```sh
tailpipe collect aws_vpc_flow_log.my_logs
```

## Query

**[Explore 12+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_vpc_flow_log)**

### Rejected Traffic

Identify rejected traffic within your VPC.

```sql
select
  start_time,
  src_addr,
  dst_addr,
  protocol,
  action
from
  aws_vpc_flow_log
where
  action = 'REJECT'
order by
  start_time desc;
```

### High-Volume Network Traffic

Identify network interfaces generating high-volume network traffic.

```sql
select
  interface_id,
  count(*) as packet_count,
  sum(coalesce(bytes, 0)) as total_bytes,
  date_trunc('minute', start_time) as event_minute
from
  aws_vpc_flow_log
where
  bytes is not null
group by
  interface_id,
  event_minute
order by
  total_bytes desc;
```

### Top 10 IP Addresses by Request Count

Identify the top 10 source IP addresses that generated the most traffic.

```sql
select
  src_addr,
  count(*) as request_count
from
  aws_vpc_flow_log
where
  src_addr is not null
group by
  src_addr
order by
  request_count desc
limit 10;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect VPC Flow Logs stored in an S3 bucket.

```hcl
connection "aws" "vpc_logging" {
  profile = "my-vpc-logging"
}

partition "aws_vpc_flow_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.vpc_logging
    bucket     = "aws-vpc-flow-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect VPC Flow Logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_vpc_flow_log" "my_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.vpc_logging
    bucket     = "aws-vpc-flow-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from local files

You can also collect VPC Flow Logs from local files.

```hcl
partition "aws_vpc_flow_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/vpc_flow_logs"]
    file_layout = `%{DATA}.json.gz`
  }
}
```

### Collect logs for all accounts in an organization

For a specific organization, collect logs for all accounts and regions.

```hcl
partition "aws_vpc_flow_log" "my_logs_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.vpc_logging
    bucket      = "vpc-flow-logs-bucket"
    file_layout = `AWSLogs/o-aa111bb222/%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz`
  }
}
```

### Collect logs for a single account

For a specific account, collect logs for all regions.

```hcl
partition "aws_vpc_flow_log" "my_logs_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.vpc_logging
    bucket      = "vpc-flow-logs-bucket"
    file_layout = `AWSLogs/(%{DATA:org_id}/)?123456789012/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz`
  }
}
```

### Collect logs for a single region

For all accounts, collect logs from `us-east-1`.

```hcl
partition "aws_vpc_flow_log" "my_logs_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.vpc_logging
    bucket      = "vpc-flow-logs-bucket"
    file_layout = `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz`
  }
}
```

### Collect logs for multiple regions

For all accounts, collect logs from us-east-1 and us-east-2.

```hcl
partition "aws_vpc_flow_log" "my_logs_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.vpc_logging
    bucket      = "vpc-flow-logs-bucket"
    file_layout = `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/(us-east-1|us-east-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz`
  }
}
```

### Exclude rejected traffic logs

Use the filter argument in your partition to exclude rejected traffic logs and reduce the size of local log storage.

```hcl
partition "aws_vpc_flow_log" "my_logs_instance" {
  filter = "action <> 'REJECT'"

  source "aws_s3_bucket" {
    connection  = connection.aws.vpc_logging
    bucket      = "vpc-flow-logs-bucket"
  }
}
```

### Collect VPC flow logs exported from CloudWatch to an S3 bucket

Collect VPC flow logs that were exported from CloudWatch to an S3 bucket, where the logs do not include a header line. Since the exported logs contain raw log lines only, a custom format block is defined to map the layout of each field.

**Important Note**: The field, `export-timestamp`, is not part of the original log event. It is automatically added by AWS during the export process from CloudWatch to S3. This timestamp represents the time when the log was exported, not when the event occurred.

```hcl
format "aws_vpc_flow_log" "exported_log_format" {
  layout = "export-timestamp instance-id interface-id pkt-srcaddr pkt-dstaddr pkt-src-aws-service az-id flow-direction start log-status packets protocol srcaddr dstaddr srcport end subnet-id"
}
```

```hcl
partition "aws_vpc_flow_log" "cw_log_group_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    format     = format.aws_vpc_flow_log.exported_log_format
    bucket     = "aws-vpc-flow-logs-123456789012-exported"
    region     = "us-east-1"
  }
}
```

### Collect logs from a CloudWatch log group

Collect VPC flow logs from all log streams in a CloudWatch log group.

```hcl
partition "aws_vpc_flow_log" "cw_log_group_logs" {
  source "aws_cloudwatch_log_group" {
    connection     = connection.aws.logging_account
    log_group_name = "aws-vpc-flow-logs-123456789012-fd33b044"
    region         = "us-east-1"
  }
}
```

### Collect VPC flow logs from a CloudWatch log group with custom log format

Collect VPC flow logs from an AWS CloudWatch Log Group where the log format does not include a header line, or the format differs from the default.

```hcl
format "aws_vpc_flow_log" "custom" {
  layout = `account-id action az-id bytes dstaddr dstport end flow-direction instance-id interface-id log-status packets pkt-dst-aws-service pkt-dstaddr pkt-src-aws-service pkt-srcaddr protocol region reject-reason srcaddr srcport start sublocation-id sublocation-type subnet-id tcp-flags traffic-path type version vpc-id`
}
```

```hcl
partition "aws_vpc_flow_log" "cw_log_group_logs" {
  source "aws_cloudwatch_log_group" {
    connection     = connection.aws.logging_account
    format         = format.aws_vpc_flow_log.custom
    log_group_name = "aws-vpc-flow-logs-123456789012-fd33b044"
    region         = "us-east-1"
  }
}
```

### Collect VPC flow logs from a CloudWatch log group for all ENIs with wildcard stream names

Collect VPC flow logs from a CloudWatch log group where each log stream corresponds to an ENI (Elastic Network Interface).

```hcl
format "aws_vpc_flow_log" "log_format" {
  layout = "instance-id interface-id pkt-srcaddr pkt-dstaddr pkt-src-aws-service az-id flow-direction start log-status packets protocol srcaddr dstaddr srcport end subnet-id"
}
```

```hcl
partition "aws_vpc_flow_log" "cw_log_group_logs" {
  source "aws_cloudwatch_log_group" {
    connection       = connection.aws.default
    format           = format.aws_vpc_flow_log.log_format
    log_group_name   = "aws-vpc-flow-logs-123456789012-fd33b044"
    log_stream_names = ["eni-*"]
    region           = "us-east-1"
  }
}
```

### Exclude logs with no status from a CloudWatch log group

Use the filter argument in your partition to exclude [records that are skipped or have no data](https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs-records-examples.html#flow-log-example-no-data) and reduce the size of local log storage.

```hcl
partition "aws_vpc_flow_log" "cw_flow_log" {
  filter = "log_status <> 'NODATA' and log_status <> 'SKIPDATA'"

  source "aws_cloudwatch_log_group"  {
    connection = connection.aws.default
    format     = format.aws_vpc_flow_log.log_layout
    log_group_name = "aws-vpc-flow-logs-123456789012-fd33b044"
    region = "us-east-1"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default                                                                                                                                                     |
| ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| file_layout | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz` |
