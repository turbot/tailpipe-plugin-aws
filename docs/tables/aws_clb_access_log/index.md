---
title: "Tailpipe Table: aws_clb_access_log - Query AWS CLB Access Logs"
description: "AWS CLB access logs capture detailed information about requests processed by a Classic Load Balancer, including client information, backend responses, and SSL details. This table provides a structured representation of the log data."
---

# Table: aws_clb_access_log - Query AWS CLB access logs

The `aws_clb_access_log` table allows you to query AWS Classic Load Balancer (CLB) access logs. This table provides detailed information about requests processed by your load balancers, including client and backend details, processing times, and SSL parameters.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_clb_access_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_clb_access_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_clb_access_log" "my_clb_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-clb-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_clb_access_log` partitions:

```sh
tailpipe collect aws_clb_access_log
```

Or for a single partition:

```sh
tailpipe collect aws_clb_access_log.my_clb_logs
```

## Query

### Failed Requests

Find failed HTTP requests (with status codes 400 and above) to troubleshoot load balancer issues.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  client_ip,
  backend_ip,
  elb_status_code,
  backend_status_code,
  request_http_version,
  request_http_method,
  request_url
from
  aws_clb_access_log
where
  elb_status_code >= 400
order by
  timestamp desc;
```

### Slow Response Times

Identify requests where the combined processing time (request + backend + response) exceeds 1 second.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request_http_version,
  request_http_method,
  request_url,
  client_ip,
  backend_ip,
  request_processing_time,
  backend_processing_time,
  response_processing_time,
  (request_processing_time + backend_processing_time + response_processing_time) as total_time
from
  aws_clb_access_log
where
  (request_processing_time + backend_processing_time + response_processing_time) > 1
order by
  total_time desc;
```

### SSL Cipher Vulnerabilities

Detect usage of deprecated or insecure SSL ciphers.

```sql
select
  ssl_cipher,
  ssl_protocol,
  count(*) as request_count
from
  aws_clb_access_log
where
  ssl_protocol in ('TLSv1.1', 'TLSv1', 'SSLv3', 'SSLv2')
group by
  ssl_cipher,
  ssl_protocol
order by
  request_count desc;
```

## Example Configurations

### Collect logs from an S3 bucket

```hcl
partition "aws_clb_access_log" "my_clb_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-clb-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

```hcl
partition "aws_clb_access_log" "my_clb_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-clb-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from local files

```hcl
partition "aws_clb_access_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/clb_logs"]
    file_layout = "%{DATA}.log"
  }
}
```

### Exclude successful requests

Use the filter argument in your partition to exclude successful requests to reduce the size of local log storage and focus on troubleshooting failed requests.

```hcl
partition "aws_clb_access_log" "my_logs_write" {
  filter = "elb_status_code != 200"

  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-clb-logs-bucket"
  }
}
```

### Collect logs for all accounts in an organization

```hcl
partition "aws_clb_access_log" "my_logs_org" {
  source "aws_s3_bucket" {
    connection  = connection.aws.logging_account
    bucket      = "aws-clb-logs-bucket"
    file_layout = "AWSLogs/o-aa111bb222/%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_%{DATA}.log"
  }
}
```

### Collect logs for a single account

```hcl
partition "aws_clb_access_log" "my_logs_account" {
  source "aws_s3_bucket" {
    connection  = connection.aws.logging_account
    bucket      = "aws-clb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_%{DATA}.log"
  }
}
```

### Collect logs for a single region

```hcl
partition "aws_clb_access_log" "my_logs_region" {
  source "aws_s3_bucket" {
    connection  = connection.aws.logging_account
    bucket      = "aws-clb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_%{DATA}.log"
  }
}
```

### Collect logs for multiple regions

```hcl
partition "aws_clb_access_log" "my_logs_regions" {
  source "aws_s3_bucket" {
    connection  = connection.aws.logging_account
    bucket      = "aws-clb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/(us-east-1|us-east-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_%{DATA}.log"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_%{DATA}.log` |
