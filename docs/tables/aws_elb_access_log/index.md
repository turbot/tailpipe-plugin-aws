---
title: "Tailpipe Table: aws_elb_access_log - Query AWS ELB Access Logs"
description: "AWS ELB Access logs capture detailed information about the requests that are processed by an Elastic Load Balancer. This table provides a structured representation of the log data, including request and response details, client and target information, processing times, and security parameters."
---

# Table: aws_elb_access_log - Query AWS ELB access logs

The `aws_elb_access_log` table allows you to query AWS Elastic Load Balancer (ELB) access logs. This table provides detailed information about requests processed by your load balancers, including client and target details, processing times, and security parameters.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_elb_access_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_elb_access_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_elb_access_log" "my_elb_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-elb-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_elb_access_log` partitions:

```sh
tailpipe collect aws_elb_access_log
```

Or for a single partition:

```sh
tailpipe collect aws_elb_access_log.my_elb_logs
```

## Query

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_elb_access_log)**

### Failed requests

Find failed HTTP requests (with status codes 400 and above) to troubleshoot load balancer issues.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  client_ip,
  target_ip,
  elb_status_code,
  target_status_code,
  request
from
  aws_elb_access_log
where
  elb_status_code >= 400
order by
  timestamp desc;
```

### Slow response times

Identify requests where the combined processing time (request + target + response) exceeds 1 second. This includes the time taken to process the request at the load balancer, the target's processing time, and the response processing time.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request,
  client_ip,
  target_ip,
  request_processing_time,  -- Time taken by load balancer to process request
  target_processing_time,   -- Time taken by target to process request
  response_processing_time, -- Time taken to process response
  (request_processing_time + target_processing_time + response_processing_time) as total_time
from
  aws_elb_access_log
where
  (request_processing_time + target_processing_time + response_processing_time) > 1 -- Requests taking longer than 1 second
order by
  total_time desc
limit 10;
```

### SSL cipher vulnerabilities

Detect usage of deprecated or insecure SSL ciphers.

```sql
select
  ssl_cipher,
  ssl_protocol,
  count(*) as request_count
from
  aws_elb_access_log
where
  ssl_protocol in ('TLSv1.1', 'TLSv1', 'SSLv3', 'SSLv2') -- Insecure protocols (TLSv1.1, TLSv1, SSLv3, SSLv2)
group by
  ssl_cipher,
  ssl_protocol
order by
  request_count desc;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect ELB access logs stored in an S3 bucket using the default log file format.

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_elb_access_log" "my_elb_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-elb-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect ELB access logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_elb_access_log" "my_elb_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-elb-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from local files

You can also collect ELB access logs from local files.

```hcl
partition "aws_elb_access_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/elb_logs"]
    file_layout = "%{DATA}.log.gz"
  }
}
```

### Exclude read-only events

Use the filter argument in your partition to exclude read-only events and reduce the size of local log storage.

```hcl
partition "aws_elb_access_log" "my_logs_write" {
  # Avoid saving read-only events, which can drastically reduce local log size
  filter = "not read_only"

  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-elb-logs-bucket"
  }
}
```

### Collect logs for all accounts in an organization

For a specific organization, collect logs for all accounts and regions.

```hcl
partition "aws_elb_access_log" "my_logs_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "elb-logs-bucket"
    file_layout = "AWSLogs/o-aa111bb222/%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log.gz"
  }
}
```

### Collect logs for a single account


For a specific account, collect logs for all regions.

```hcl
partition "aws_elb_access_log" "my_logs_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "elb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log.gz"
  }
}
```

### Collect logs for a single region

For all accounts, collect logs from us-east-1.

```hcl
partition "aws_elb_access_log" "my_logs_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "elb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log.gz"
  }
}
```

### Collect logs for multiple regions

For all accounts, collect logs from us-east-1 and us-east-2.

```hcl
partition "aws_elb_access_log" "my_logs_regions" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "elb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/(us-east-1|us-east-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log.gz"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.log.gz` |