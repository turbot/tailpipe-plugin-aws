---
title: "Source: aws_s3_bucket - Collect logs from AWS S3 buckets"
description: "Allows users to collect logs from AWS S3 buckets."
---

# Source: aws_s3_bucket - Collect logs from AWS S3 buckets

An AWS S3 bucket is a cloud storage resource used to store objects like data files and metadata. It serves as a central repository for logs from AWS services such as CloudTrail, ALB, VPC flow logs, and more.

Using this source, you can collect, filter, and analyze logs stored in S3 buckets, enabling system monitoring, security investigations, and compliance reporting.

Most AWS tables define a default `file_layout` for the `aws_s3_bucket` source, so if your AWS logs are stored in default log locations, you don't need to override the `file_layout` argument.

The trailing `/` is not automatically included in the `prefix`. If your log path requires it, be sure to add it explicitly.

## Example Configurations

### Collect CloudTrail logs

Collect CloudTrail logs for all accounts and regions.

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

### Collect CloudTrail logs with a prefix

Collect CloudTrail logs stored with an S3 key prefix.

```hcl
partition "aws_cloudtrail_log" "my_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-cloudtrail-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect CloudTrail logs with a custom path

Collect CloudTrail logs stored in an S3 bucket with a custom log file format.

```hcl
partition "aws_cloudtrail_log" "my_logs_custom_path" {
  source "aws_s3_bucket" {
    connection  = connection.aws.logging_account
    bucket      = "aws-cloudtrail-logs-bucket"
    file_layout = `CustomLogs/Dev/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz`
  }
}
```

### Collect S3 access logs for a specific date from non-date based partition logs

Collect logs in an S3 bucket stored with [non-date based partitioning](https://docs.aws.amazon.com/AmazonS3/latest/userguide/ServerLogs.html) using a prefix to only retrieve files for a specific day.

```hcl
partition "aws_s3_server_access_log" "my_s3_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-s3-server-access-logs"
    prefix     = "2025-06-07"
  }
}
```

## Arguments

| Argument     | Type            | Required | Default                  | Description                                                                                                                   |
|-------------|------------------|----------|--------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| bucket      | String           | Yes      |                          | The name of the S3 bucket to collect logs from.                                                                               |
| connection  | `connection.aws` | No       | `connection.aws.default` | The [AWS connection](https://hub.tailpipe.io/plugins/turbot/aws#connection-credentials) to use to connect to the AWS account. |
| file_layout | String           | No       |                          | The Grok pattern that defines the log file structure.                                                                         |
| prefix      | String           | No       |                          | The S3 key prefix that comes after the name of the bucket you have designated for log file delivery.                          |

### Table Defaults

The following tables define their own default values for certain source arguments:

- **[aws_alb_access_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_alb_access_log#aws_s3_bucket)**
- **[aws_clb_access_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_clb_access_log#aws_s3_bucket)**
- **[aws_cloudtrail_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cloudtrail_log#aws_s3_bucket)**
- **[aws_guardduty_finding](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_guardduty_finding#aws_s3_bucket)**
- **[aws_nlb_access_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_nlb_access_log#aws_s3_bucket)**
- **[aws_s3_server_access_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_s3_server_access_log#aws_s3_bucket)**
- **[aws_cost_and_usage_focus](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_and_usage_focus#aws_s3_bucket)**
- **[aws_cost_and_usage_report](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_and_usage_report#aws_s3_bucket)**
- **[aws_cost_optimization_recommendation](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_optimization_recommendation#aws_s3_bucket)**
- **[aws_vpc_flow_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_vpc_flow_log#aws_s3_bucket)**
- **[aws_waf_traffic_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_waf_traffic_log#aws_s3_bucket)**
