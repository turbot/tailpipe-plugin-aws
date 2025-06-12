---
title: "Tailpipe Table: aws_alb_connection_log - Query AWS ALB Connection Logs"
description: "AWS ALB Connection logs capture detailed information about connection attempts to an Application Load Balancer, including TLS handshake details, client certificate data, and connection traceability identifiers."
---

# Table: aws_alb_connection_log - Query AWS ALB connection logs

The `aws_alb_connection_log` table allows you to query AWS Application Load Balancer (ALB) connection logs. This table provides detailed information about connection attempts to your load balancers, including TLS handshake details, client certificate data, and connection traceability identifiers.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_alb_connection_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_alb_connection_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_alb_connection_log" "my_alb_conn_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-alb-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_alb_connection_log` partitions:

```sh
tailpipe collect aws_alb_connection_log
```

Or for a single partition:

```sh
tailpipe collect aws_alb_connection_log.my_alb_conn_logs
```

## Query

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_alb_connection_log)**

### Failed TLS handshakes

Find failed TLS handshakes to troubleshoot connection issues.

```sql
select
  timestamp,
  tp_index as conn_trace_id,
  client_ip,
  client_port,
  tls_protocol,
  tls_cipher,
  tls_verify_status
from
  aws_alb_connection_log
where
  tls_verify_status like 'Failed:%'
order by
  timestamp desc;
```

### Slow TLS handshakes

Identify connections with unusually high TLS handshake latency that might indicate performance issues.

```sql
select
  timestamp,
  tp_index as conn_trace_id,
  client_ip,
  client_port,
  tls_protocol,
  tls_cipher,
  tls_handshake_latency
from
  aws_alb_connection_log
where
  tls_handshake_latency > 1 -- Handshakes taking longer than 1 second
order by
  tls_handshake_latency desc
limit 10;
```

### Deprecated TLS protocols

Identify connections using deprecated or insecure TLS protocols.

```sql
select
  tls_protocol,
  tls_cipher,
  count(*) as connection_count
from
  aws_alb_connection_log
where
  tls_protocol in ('TLSv1.1', 'TLSv1', 'SSLv3', 'SSLv2') -- Insecure protocols
group by
  tls_protocol,
  tls_cipher
order by
  connection_count desc;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect ALB connection logs stored in an S3 bucket using the default log file format.

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_alb_connection_log" "my_alb_conn_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-alb-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect ALB connection logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_alb_connection_log" "my_alb_conn_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-alb-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from local files

You can also collect ALB connection logs from local files.

```hcl
partition "aws_alb_connection_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/elb_logs"]
    file_layout = "%{DATA}.log.gz"
  }
}
```

### Collect logs for all accounts in an organization

For a specific organization, collect logs for all accounts and regions.

```hcl
partition "aws_alb_connection_log" "my_logs_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-alb-logs-bucket"
    file_layout = "AWSLogs/o-aa111bb222/%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/conn_log.%{DATA}.log.gz"
  }
}
```

### Collect logs for a single account

For a specific account, collect logs for all regions.

```hcl
partition "aws_alb_connection_log" "my_logs_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-alb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/conn_log.%{DATA}.log.gz"
  }
}
```

### Collect logs for a single region

For all accounts, collect logs from us-east-1.

```hcl
partition "aws_alb_connection_log" "my_logs_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-alb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/conn_log.%{DATA}.log.gz"
  }
}
```

### Collect logs for multiple regions

For all accounts, collect logs from us-east-1 and us-east-2.

```hcl
partition "aws_alb_connection_log" "my_logs_regions" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-alb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/(us-east-1|us-east-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/conn_log.%{DATA}.log.gz"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/conn_log.%{DATA}.log.gz` |