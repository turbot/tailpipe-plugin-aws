---
title: "Tailpipe Table: aws_lambda_log - Query AWS Lambda Logs"
description: "Allows users to query AWS Lambda logs."
---

# Table: aws_lambda_log - Query AWS Lambda Logs

*TODO*: Add description

## Table Usage Guide

The `aws_lambda_log` table allows you to query data from AWS Lambda logs. This table provides detailed information about AWS Lambda function invocations, including the function name, request ID, duration, memory usage, and more.

## Examples

### Total Invocations by Function

Calculates the total number of invocations for each AWS Lambda function to identify the most frequently invoked functions.

```sql
select
  function_name,
  count(*) as invocation_count
from
    aws_lambda_log
group by
    function_name
order by
    invocation_count desc
limit 10;
```