---
title: "Source: aws_cloudwatch_log_group - Obtain logs from AWS CloudWatch Log Groups"
description: "Allows users to collect logs from AWS CloudWatch Log Groups."
---

# Source: aws_cloudwatch_log_group - Obtain logs from AWS CloudWatch Log Groups

An AWS CloudWatch Log Group is a collection of log streams, which are sequences of log events. Log events are records of activity in a system, such as actions taken by users or applications. CloudWatch Log Groups are used to store and manage log data, which can be searched, monitored, and archived.

## Configuration

| Property            | Description                                                | Default                         |
|---------------------|------------------------------------------------------------|---------------------------------|
| `connection`        | The connection to use to connect to the AWS account.       | -                               |
| `log_group_name`    | The name of the CloudWatch Log Group to collect logs from. | -                               |
| `log_stream_prefix` | A prefix to find matching log streams.                     | -                               |
| `start_time`        | The start time to collect logs from.                       | -                               |
| `end_time`          | The end time to collect logs until.                        | Defaults to current time        |
| `region`            | The AWS region where the log group is located.             | Defaults to connections default |