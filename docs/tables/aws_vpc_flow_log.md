---
title: "Tailpipe Table: aws_vpc_flow_log - Query AWS VPC Flow Logs"
description: "Allows users to query AWS VPC flow logs."
---

# Table: aws_vpc_flow_log - Query AWS VPC Flow Logs

*TODO*: Add description

## Table Usage Guide

The `aws_vpc_flow_log` table allows you to query data from AWS VPC flow logs. This table provides detailed information about network traffic within your VPC, including the source and destination IP addresses, ports, protocol, and more.

## Examples

### Top Source IP Addresses

Displays the top source IP addresses by the number of bytes sent, helping to identify the most active sources of network traffic.

```sql
select
  srcaddr,
  sum(bytes) as total_bytes_sent
from
    aws_vpc_flow_log
group by
    srcaddr
order by
    total_bytes_sent desc
limit 10;
```