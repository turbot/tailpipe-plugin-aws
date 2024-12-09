---
title: "Tailpipe Table: aws_elb_access_log - Query AWS Elastic Load Balancer Access Logs"
description: "Allows users to query AWS Elastic Load Balancer access logs."
---

# Table: aws_elb_access_log - Query AWS Elastic Load Balancer Access Logs

*TODO*: Add description

## Table Usage Guide

The `aws_elb_access_log` table allows you to query data from AWS Elastic Load Balancer access logs. This table provides detailed information about client requests made to your Elastic Load Balancer, including the request time, client IP address, request method, response status code, and more.

## Examples

### Top Requested URLs

Displays the top requested URLs by the number of requests, helping to identify the most frequently accessed resources.

```sql
select
  request_uri,
  count(*) as request_count
from
    aws_elb_access_log
group by
    request_uri
order by
    request_count desc
limit 10;
```
