---
title: "Source: aws_cloudwatch_log_group - Collect logs from AWS CloudWatch Log Groups"
description: "Allows users to collect logs from AWS CloudWatch log groups."
---

# Source: aws_cloudwatch_log_group - Collect logs from AWS CloudWatch log groups

AWS CloudWatch log groups are collections of log streams that share the same retention, monitoring, and access control settings. They serve as containers for log data from AWS services, containerized applications, and custom applications.

Using this source, you can collect and analyze logs from CloudWatch log groups, enabling real-time monitoring, troubleshooting, and analysis of your AWS resources and applications.

## Example Configurations

### Collect CloudTrail logs

Collect CloudTrail logs for all accounts and regions.

```hcl
connection "aws" "default" {
  profile = "my-aws-profile"
}

partition "aws_cloudtrail_log" "cw_log_group_logs" {
  source "aws_cloudwatch_log_group" {
    connection     = connection.aws.default
    log_group_name = "aws-cloudtrail-logs-123456789012-fd33b044"
    region         = "us-east-1"
  }
}
```

### Collect CloudTrail logs with a prefix

Collect CloudTrail logs for account ID `456789012345` in us-east-1.

```hcl
partition "aws_cloudtrail_log" "cw_log_group_logs_prefix" {
  source "aws_cloudwatch_log_group" {
    connection        = connection.aws.default
    log_group_name    = "aws-cloudtrail-logs-123456789012-fd33b044"
    log_stream_prefix = "456789012345_CloudTrail_us-east-1"
    region            = "us-east-1"
  }
}
```

## Arguments

| Argument          | Type             | Required | Default                  | Description                                                                                                                   |
| ----------------- | ---------------- | -------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------------- |
| connection        | `connection.aws` | No       | `connection.aws.default` | The [AWS connection](https://hub.tailpipe.io/plugins/turbot/aws#connection-credentials) to use to connect to the AWS account. |
| log_group_name    | String           | Yes      |                          | The name of the CloudWatch log group to collect logs from.                                                                    |
| log_stream_prefix | String           | No       |                          | Collect logs from log streams whose names begin the specified prefix.                                                         |
| region            | String           | Yes      |                          | The AWS region where the log group is located.                                                                                |
