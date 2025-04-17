---
title: "Tailpipe Table: aws_vpc_flow_log - Query AWS VPC Flow Logs"
description: "AWS VPC flow logs capture information about IP traffic going to and from network interfaces in your VPC."
---

# Table: aws_vpc_flow_log - Query AWS VPC Flow Logs

The `aws_vpc_flow_log` table allows you to query data from [AWS VPC Flow Logs](https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html). This table provides detailed insights into network traffic within your VPC, including source and destination IP addresses, ports, protocols, and more.

**Note**: For timestamp information, the `start` field will be used first, with the `end` field as a fallback. If neither field is available, then that log line will not be collected and Tailpipe will return an error.

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

### Exclude skipped logs and logs with no data

Use the filter argument in your partition to exclude [records that are skipped or have no data](https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs-records-examples.html#flow-log-example-no-data) and reduce the size of local log storage.

```hcl
partition "aws_vpc_flow_log" "my_logs_status_ok" {
  filter = "log_status = 'OK'"

  source "aws_s3_bucket" {
    connection  = connection.aws.vpc_logging
    bucket      = "vpc-flow-logs-bucket"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/vpcflowlogs/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/(%{NUMBER:hour}/)?%{DATA}.log.gz` |

