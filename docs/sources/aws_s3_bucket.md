---
title: "Source: aws_s3_bucket - Collect logs from AWS S3 buckets"
description: "Allows users to collect logs from AWS S3 buckets."
---

# Source: aws_s3_bucket - Obtain logs from AWS S3 buckets

An AWS S3 bucket is a cloud storage resource used to store objects like data files and metadata. It serves as a central repository for logs from AWS services such as CloudTrail, ELB, and VPC flow logs, and more.

Using this source, you can collect, filter, and analyze logs stored in S3, enabling system monitoring, security investigations, and compliance reporting.

## Examples

Collect all CloudTrail logs:

```hcl
partition "aws_cloudtrail_log" "all" {
  source "aws_s3_bucket"  {
    bucket = "my-logging-bucket"
  }
}
```

Using a specific AWS connection, collect all CloudTrail logs:

```hcl
partition "aws_cloudtrail_log" "dev" {
  source "aws_s3_bucket"  {
    connection = connection.aws.dev
    bucket     = "my-logging-bucket"
  }
}
```

Collect logs from us-east-1 only:

```hcl
partition "aws_cloudtrail_log" "cloudtrail_logs" {
  source "aws_s3_bucket"  {
    bucket        = "turbot-632902152528-us-east-1"
    file_layout   = "AWSLogs/%{NUMBER:account_id}/CloudTrail/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{NUMBER:prefix}_CloudTrail_%{DATA:region_file}_%{YEAR:year_file}%{MONTHNUM:month_file}%{MONTHDAY:day_file}T%{HOUR:hour}%{MINUTE:minute}Z_%{DATA:suffix}.json.gz"
  }
}
```

## Arguments

| Argument                | Description                                                                           | Default       |
|-------------------------|---------------------------------------------------------------------------------------|---------------|
| `bucket`                | The name of the S3 bucket to collect logs from.                                      |               |
| `connection`            | The connection to use to connect to the AWS account.                                 |               |
| `file_layout`           | Regex of pattern filename layout, used to extract information such as year, month, day, etc. |               |
| `lexicographical_order` | Used to indicate log files are in lexicographical order.                             | `false`       |
| `start_after_key`       | The key to start collecting logs from.                                               |               |
