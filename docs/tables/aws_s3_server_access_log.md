---
title: "Tailpipe Table: aws_s3_server_access_log - Query AWS S3 Server Access Logs"
description: "Allows users to query AWS S3 server access logs."
---

# Table: aws_s3_server_access_log - Query AWS S3 Server Access Logs

*TODO*: Add description

## Table Usage Guide

The `aws_s3_server_access_log` table allows you to query data from AWS S3 server access logs. This table provides detailed information about requests made to your S3 bucket, including the request time, requestor IP address, request method, response status code, and more.

## Examples

### Top Requested URLs

Displays the top requested URLs by the number of requests, helping to identify the most frequently accessed resources.

```sql
select
  request_uri,
  count(*) as request_count
from
    aws_s3_server_access_log
group by
    request_uri
order by
    request_count desc
limit 10;
```