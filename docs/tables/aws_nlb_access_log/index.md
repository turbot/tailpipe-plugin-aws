---
title: "Tailpipe Table: aws_nlb_access_log - Query AWS NLB Access Logs"
description: "AWS NLB Access logs capture detailed information about the requests that are processed by a Network Load Balancer. This table provides a structured representation of the log data, including client and destination information, connection times, TLS parameters, and network traffic statistics."
---

# Table: aws_nlb_access_log - Query AWS NLB access logs

The `aws_nlb_access_log` table allows you to query AWS Network Load Balancer (NLB) access logs. This table provides detailed information about connections processed by your load balancers, including client and destination details, connection times, TLS parameters, and network traffic statistics.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_nlb_access_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_nlb_access_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_nlb_access_log" "my_nlb_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-nlb-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_nlb_access_log` partitions:

```sh
tailpipe collect aws_nlb_access_log
```

Or for a single partition:

```sh
tailpipe collect aws_nlb_access_log.my_nlb_logs
```

## Query

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_nlb_access_log)**

### High latency connections

Identify connections with unusually high connection establishment times to troubleshoot network latency issues.

```sql
select
  timestamp,
  elb,
  client_ip,
  client_port,
  destination_ip,
  destination_port,
  connection_time
from
  aws_nlb_access_log
where
  connection_time > 1000 -- Connections taking longer than 1 second to establish
order by
  connection_time desc
limit 10;
```

### TLS handshake performance issues

Find connections with slow TLS handshake times that might indicate security configuration or network issues.

```sql
select
  timestamp,
  elb,
  client_ip,
  client_port,
  destination_ip,
  destination_port,
  tls_handshake_time,
  tls_cipher,
  tls_protocol_version
from
  aws_nlb_access_log
where
  tls_handshake_time > 500 -- TLS handshakes taking longer than 500ms
order by
  tls_handshake_time desc
limit 10;
```

### TLS protocol vulnerabilities

Detect usage of deprecated or insecure TLS protocols.

```sql
select
  tls_protocol_version,
  tls_cipher,
  count(*) as connection_count
from
  aws_nlb_access_log
where
  tls_protocol_version in ('tlsv1.1', 'tlsv1', 'sslv3', 'sslv2') -- Insecure protocols
group by
  tls_protocol_version,
  tls_cipher
order by
  connection_count desc;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect NLB access logs stored in an S3 bucket using the default log file format.

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_nlb_access_log" "my_nlb_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-nlb-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect NLB access logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_nlb_access_log" "my_nlb_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-nlb-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from local files

You can also collect NLB access logs from local files.

```hcl
partition "aws_nlb_access_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/elb_logs"]
    file_layout = "%{DATA}.log.gz"
  }
}
```

### Exclude read-only events

Use the filter argument in your partition to exclude read-only events and reduce the size of local log storage.

```hcl
partition "aws_nlb_access_log" "my_logs_write" {
  # Avoid saving read-only events, which can drastically reduce local log size
  filter = "not read_only"

  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-nlb-logs-bucket"
  }
}
```

### Collect logs for all accounts in an organization

For a specific organization, collect logs for all accounts and regions.

```hcl
partition "aws_nlb_access_log" "my_logs_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-nlb-logs-bucket"
    file_layout = "AWSLogs/o-aa111bb222/%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_net.[^_]+_[^_]+_[^.]+.log.gz"
  }
}
```

### Collect logs for a single account

For a specific account, collect logs for all regions.

```hcl
partition "aws_nlb_access_log" "my_logs_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-nlb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_net.[^_]+_[^_]+_[^.]+.log.gz"
  }
}
```

### Collect logs for a single region

For all accounts, collect logs from us-east-1.

```hcl
partition "aws_nlb_access_log" "my_logs_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-nlb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_net.[^_]+_[^_]+_[^.]+.log.gz"
  }
}
```

### Collect logs for multiple regions

For all accounts, collect logs from us-east-1 and us-east-2.

```hcl
partition "aws_nlb_access_log" "my_logs_regions" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-nlb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/(us-east-1|us-east-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_net.[^_]+_[^_]+_[^.]+.log.gz"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_net.[^_]+_[^_]+_[^.]+.log.gz` |