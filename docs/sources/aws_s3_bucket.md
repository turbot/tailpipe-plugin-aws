---
title: "Source: aws_s3_bucket - Collect logs from AWS S3 buckets"
description: "Allows users to collect logs from AWS S3 buckets."
---

# Source: aws_s3_bucket - Collect logs from AWS S3 buckets

An AWS S3 bucket is a cloud storage resource used to store objects like data files and metadata. It serves as a central repository for logs from AWS services such as CloudTrail, ELB, VPC flow logs, and more.

Using this source, you can collect, filter, and analyze logs stored in S3 buckets, enabling system monitoring, security investigations, and compliance reporting.

Most AWS tables define a default `file_path` for the `aws_s3_bucket` source, so if your AWS logs are stored in default log locations, you don't need to override the `file_path` argument.

Table example configurations:

- **[aws_cloudtrail_log](https://tailpipe.io/plugins/turbot/aws/tables/aws_cloudtrail_log#example-configurations)**

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

```hcl
partition "aws_cloudtrail_log" "my_logs_custom_path" {
  source "aws_s3_bucket" {
    connection  = connection.aws.logging_account
    bucket      = "aws-cloudtrail-logs-bucket"
    file_layout = "CustomLogs/Dev/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"
  }
}
```

## Arguments

| Argument      | Required | Default                  | Description                                                                                                                |
|---------------|----------|--------------------------|----------------------------------------------------------------------------------------------------------------------------|
| bucket        | Yes      |                          | The name of the S3 bucket to collect logs from.                                                                            |
| connection    | No       | `connection.aws.default` | The [AWS connection](https://tailpipe.io/docs/reference/config-files/connection/aws) to use to connect to the AWS account. |
| file_layout   | No       |                          | The Grok pattern that defines the log file structure.                                                                      |
| prefix        | No       |                          | The S3 key prefix that comes after the name of the bucket you have designated for log file delivery.                       |

### Table Defaults

The following tables define their own default values for certain source arguments:

- **[aws_cloudtrail_log](https://tailpipe.io/plugins/turbot/aws/tables/aws_cloudtrail_log#aws_s3_bucket)**
