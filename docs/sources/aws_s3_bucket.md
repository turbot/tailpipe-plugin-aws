---
title: "Source: aws_s3_bucket - Obtain logs from AWS S3 buckets"
description: "Allows users to collect logs from AWS S3 buckets."
---

# Source: aws_s3_bucket - Obtain logs from AWS S3 buckets

An AWS S3 Bucket is a public cloud storage resource available in Amazon Web Services' (AWS) Simple Storage Service (S3). It is used to store objects, which consist of data and its descriptive metadata. S3 makes it possible to store and retrieve varying amounts of data, at any time, from anywhere on the web.

## Configuration 

| Property | Description                                                                                  | Default                   |
| - |----------------------------------------------------------------------------------------------|---------------------------|
| `connection` | The connection to use to connect to the AWS account.                                         | -                         |
| `bucket` | The name of the S3 bucket to collect logs from.                                              | -                         |
| `prefix` | The prefix to filter objects in the bucket.                                                  | Defaults to bucket root.  |
| `region` | The AWS region where the bucket is located.                                                  | Defaults to `us-east-1`   |
| `extensions` | The file extensions to collect.                                                              | Defaults to all files.    |
| `lexicographical_order` | Used to indicate log files are in lexicographical order.                                     | Defaults to `false`       |
| `start_after_key` | The key to start collecting logs from.                                                       | -                         |
| `file_layout` | Regex of pattern filename layout, used to extract information such as year, month, day, etc. | Default depends on Table. |
