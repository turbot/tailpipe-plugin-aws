---
title: "Source: aws_cloudwatch_log_group - Collect logs from AWS CloudWatch log groups"
description: "Allows users to collect logs from AWS CloudWatch Log Groups."
---

# Source: aws_cloudwatch_log_group - Collect logs from AWS CloudWatch log groups

AWS CloudWatch Log Groups are collections of log streams that share the same retention, monitoring, and access control settings. They serve as containers for log data from AWS services, containerized applications, and custom applications.

Using this source, you can collect and analyze logs from CloudWatch Log Groups, enabling real-time monitoring, troubleshooting, and analysis of your AWS resources and applications.

## Example Configurations

### Collect all logs from a Log Group

Collect all logs from a specific CloudWatch Log Group.

```hcl
connection "aws" "default" {
  profile = "my-aws-profile"
}

partition "aws_cloudwatch_log" "application_logs" {
  source "aws_cloudwatch_log_group" {
    connection    = connection.aws.default
    log_group_name = "/aws/lambda/my-function"
  }
}
```

### Collect logs from specific Log Streams

Collect logs from Log Streams that match a specific prefix.

```hcl
partition "aws_cloudwatch_log" "filtered_logs" {
  source "aws_cloudwatch_log_group" {
    connection       = connection.aws.default
    log_group_name  = "/aws/ecs/my-cluster"
    log_stream_prefix = "my-service"
  }
}
```

### Collect logs from a specific region

Collect logs from a Log Group in a specific AWS region.

```hcl
partition "aws_cloudwatch_log" "regional_logs" {
  source "aws_cloudwatch_log_group" {
    connection    = connection.aws.default
    log_group_name = "/aws/containerinsights/my-cluster"
    region        = "us-west-2"
  }
}
```

## Arguments

| Argument          | Type             | Required | Default                  | Description                                                                                                                   |
| ----------------- | ---------------- | -------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------------- |
| log_group_name    | String           | Yes      |                          | The name of the CloudWatch Log Group to collect logs from.                                                                    |
| connection        | `connection.aws` | No       | `connection.aws.default` | The [AWS connection](https://hub.tailpipe.io/plugins/turbot/aws#connection-credentials) to use to connect to the AWS account. |
| log_stream_prefix | String           | No       |                          | The prefix to filter Log Streams within the Log Group.                                                                        |
| region            | String           | No       | "us-east-1"              | The AWS region where the Log Group is located.                                                                                |

### Table Defaults

The following tables define their own default values for certain source arguments:

- **[aws_waf_traffic_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_waf_traffic_log#aws_cloudwatch_log_group)**
- **[aws_cloudtrail_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cloudtrail_log#aws_cloudwatch_log_group)**
