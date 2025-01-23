---
title: "Source: aws_s3_bucket - Collect logs from AWS S3 buckets"
description: "Allows users to collect logs from AWS S3 buckets."
---

# Source: aws_s3_bucket - Collect logs from AWS S3 buckets

An AWS S3 bucket is a cloud storage resource used to store objects like data files and metadata. It serves as a central repository for logs from AWS services such as CloudTrail, ELB, VPC flow logs, and more.

Using this source, you can collect, filter, and analyze logs stored in S3 buckets, enabling system monitoring, security investigations, and compliance reporting.

Each AWS table defines a default `file_path` for the `aws_s3_bucket` source, so if your logs are stored in default AWS log locations, you don't need to override the `file_path` argument.

## CloudTrail Log Examples

### Collect logs

Collect CloudTrail logs for all accounts and regions.

```hcl
partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket"  {
    bucket = "cloudtrail-s3-log-bucket"
  }
}
```

### Specify a prefix

Collect CloudTrail logs stored with an S3 key prefix.

```hcl
partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket"  {
    bucket = "cloudtrail-s3-log-bucket"
    prefix = "sample/prefix/"
  }
}
```

### Use an AWS connection

Use a specific AWS connection when connecting to the AWS account.

```hcl
connection "aws" "dev" {
  profile = "dev"
}

partition "aws_cloudtrail_log" "dev" {
  source "aws_s3_bucket"  {
    connection = connection.aws.dev
    bucket     = "cloudtrail-s3-log-bucket"
  }
}
```

### Collect for an account

Collect CloudTrail logs for a specific account in all regions.

```hcl
partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket"  {
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/CloudTrail/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"
  }
}
```

### Collect for a specific region

Collect CloudTrail logs for all accounts in us-east-1.

```hcl
partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket"  {
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"
  }
}
```

### Collect for multiple regions

Collect CloudTrail logs for all accounts in us-east-1 and us-east-2.

```hcl
partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket"  {
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/(us-east-1|us-east-2)/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"
  }
}
```

### Collect logs in a custom path

```hcl
partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket"  {
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = "CustomLogs/Dev/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"
  }
}
```

## Arguments

| Argument      | Description                                                                                           |
|---------------|-------------------------------------------------------------------------------------------------------|
| `bucket`      | The name of the S3 bucket to collect logs from.                                                       |
| `connection`  | The connection to use to connect to the AWS account.                                                  |
| `file_layout` | Pattern filename layout using Grok pattern, used to extract information such as year, month, day, etc.|
| `prefix`      | The S3 key prefix that comes after the name of the bucket you have designated for log file delivery.  |
