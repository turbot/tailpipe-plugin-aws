---
title: "Tailpipe Table: aws_alb_access_log - Query AWS ALB Access Logs"
description: "AWS ALB access logs capture detailed information about the requests that are processed by an Application Load Balancer. This table provides a structured representation of the log data, including request and response details, client and target information, processing times, and security parameters."
---

# Table: aws_alb_access_log - Query AWS ALB Access Logs

The `aws_alb_access_log` table allows you to query AWS Application Load Balancer (ALB) access logs. This table provides detailed information about requests processed by your load balancers, including client and target details, processing times, and security parameters.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_alb_access_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_alb_access_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_alb_access_log" "my_alb_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-alb-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_alb_access_log` partitions:

```sh
tailpipe collect aws_alb_access_log
```

Or for a single partition:

```sh
tailpipe collect aws_alb_access_log.my_alb_logs
```

## Query

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_alb_access_log)**

### Failed Requests

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
  request_url,
  request_http_method,
  request_http_version
from
  aws_alb_access_log
where
  elb_status_code >= 400
order by
  timestamp desc;
```

### Slow Response Times

Identify requests where the combined processing time (request + target + response) exceeds 1 second. This includes the time taken to process the request at the load balancer, the target's processing time, and the response processing time.

```sql
select
  timestamp,
  elb,
  tp_index as account_id,
  request_url,
  request_http_method,
  request_http_version,
  client_ip,
  target_ip,
  request_processing_time,  -- Time taken by load balancer to process request
  target_processing_time,   -- Time taken by target to process request
  response_processing_time, -- Time taken to process response
  (request_processing_time + target_processing_time + response_processing_time) as total_time
from
  aws_alb_access_log
where
  (request_processing_time + target_processing_time + response_processing_time) > 1 -- Requests taking longer than 1 second
order by
  total_time desc
limit 10;
```

### SSL Cipher Vulnerabilities

Detect usage of deprecated or insecure SSL ciphers.

```sql
select
  ssl_cipher,
  ssl_protocol,
  count(*) as request_count
from
  aws_alb_access_log
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

Collect ALB access logs stored in an S3 bucket using the default log file format.

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_alb_access_log" "my_alb_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-alb-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect ALB access logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_alb_access_log" "my_alb_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-alb-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from local files

You can also collect ALB access logs from local files.

```hcl
partition "aws_alb_access_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/elb_logs"]
    file_layout = "%{DATA}.log.gz"
  }
}
```

### Exclude successful requests

Use the filter argument in your partition to exclude successful requests to reduce the size of local log storage and focus on troubleshooting failed requests.

```hcl
partition "aws_alb_access_log" "my_alb_logs_filtered" {
  filter = "elb_status_code != 200"

  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-alb-logs-bucket"
  }
}
```

### Collect logs for all accounts in an organization

For a specific organization, collect logs for all accounts and regions.

```hcl
partition "aws_alb_access_log" "my_logs_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-alb-logs-bucket"
    file_layout = "AWSLogs/o-aa111bb222/%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_app.%{DATA}.log.gz"
  }
}
```

### Collect logs for a single account

For a specific account, collect logs for all regions.

```hcl
partition "aws_alb_access_log" "my_logs_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-alb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_app.%{DATA}.log.gz"
  }
}
```

### Collect logs for a single region

For all accounts, collect logs from us-east-1.

```hcl
partition "aws_alb_access_log" "my_logs_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-alb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_app.%{DATA}.log.gz"
  }
}
```

### Collect logs for multiple regions

For all accounts, collect logs from us-east-1 and us-east-2.

```hcl
partition "aws_alb_access_log" "my_logs_regions" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.logging_account
    bucket      = "aws-alb-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/(us-east-1|us-east-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_app.%{DATA}.log.gz"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/elasticloadbalancing/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:account_id}_elasticloadbalancing_%{DATA:region}_app.%{DATA}.log.gz` |
