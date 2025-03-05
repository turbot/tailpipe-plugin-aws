## Activity Examples

### Daily Network Traffic Trends

Count VPC flow log entries per day to identify network activity trends. This helps monitor overall network behavior and detect unusual spikes in traffic.

```sql
select
  strftime(start_time, '%Y-%m-%d') as traffic_date,
  count(*) as request_count
from
  aws_vpc_flow_log
group by
  traffic_date
order by
  traffic_date asc;
```

```yaml
folder: VPC
```

### Top 10 Source IPs Generating Traffic

Identify the top 10 source IP addresses that generated the most network traffic. This helps detect potential high-traffic sources, including misconfigured applications or malicious activities.

```sql
select
  src_addr,
  count(*) as request_count
from
  aws_vpc_flow_log
group by
  src_addr
order by
  request_count desc
limit 10;
```

```yaml
folder: VPC
```

## Detection Examples

### Identify Traffic from a Suspicious IP

Check if a specific IP is sending or receiving traffic. This is useful for investigating potential threats or monitoring known suspicious IPs.

```sql
select
  start_time,
  src_addr,
  dst_addr,
  src_port,
  dst_port,
  protocol,
  action
from 
  aws_vpc_flow_log
where 
  src_addr = '192.0.2.100'
  or dst_addr = '192.0.2.100'
order by
  start_time desc;
```

```yaml
folder: VPC
```

### Detect Suspicious Traffic from External IPs

Identify inbound traffic from external (non-VPC) IP addresses. This helps detect unauthorized or unexpected external connections.

```sql
select
  start_time,
  src_addr,
  dst_addr,
  action,
  protocol,
  region,
  vpc_id,
  action
from 
  aws_vpc_flow_log
where 
  flow_direction = 'ingress'
  and (
    src_addr not like '10.%'  -- Exclude all 10.0.0.0/8
    and src_addr not like '192.168.%' -- Exclude all 192.168.0.0/16
    and (
      src_addr < '172.16.0.0' or src_addr > '172.31.255.255' -- Exclude full 172.16.0.0/12 range
    )
    and src_addr not like '169.254.%' -- Exclude link-local 169.254.0.0/16
    and src_addr not like '127.%' -- Exclude localhost 127.0.0.0/8
  )
order by 
  start_time desc;
```

```yaml
folder: VPC
```

### Detect Unauthorized Access Attempts

Identify unauthorized attempts to access instances through uncommon ports. This helps detect brute-force attempts or suspicious access patterns.

```sql
select
  start_time,
  src_addr,
  dst_addr,
  src_port,
  dst_port,
  protocol,
  vpc_id,
  action
from
  aws_vpc_flow_log
where
  action = 'REJECT'
  and dst_port not in (22, 80, 443)
order by
  start_time desc;
```

```yaml
folder: VPC
```

## Operational Examples

### Traffic to a Specific Subnet

Retrieve failed (rejected) network connection attempts within a specific subnet. This helps analyze access control issues or misconfigurations.

```sql
select
  start_time,
  src_addr,
  dst_addr,
  src_port,
  dst_port,
  protocol,
  subnet_id,
  action
from 
  aws_vpc_flow_log
where 
  subnet_id = 'subnet-027e9a6d4add894eb'
order by 
  start_time desc;
```

```yaml
folder: VPC
```

### Identify High-Latency Network Paths

Detect network paths with high packet loss, which may indicate congestion or misconfigured routes.

```sql
select
  start_time,
  src_addr,
  dst_addr,
  packets,
  bytes,
  protocol,
  vpc_id
from
  aws_vpc_flow_log
where
  packets > 1000
  and bytes < 5000 -- Low bytes despite high packet count
order by
  start_time desc;
```

```yaml
folder: VPC
```

## Volume Examples

### Unusually Large Data Transfers

Identify unusually large outbound traffic based on bytes transferred. This helps detect data exfiltration attempts or misconfigured workloads.

```sql
select
  start_time,
  src_addr,
  dst_addr,
  bytes,
  packets,
  protocol,
  vpc_id
from
  aws_vpc_flow_log
where
  bytes > 500000000 -- 500MB
order by
  bytes desc;
```

```yaml
folder: VPC
```

### High-Volume Network Traffic

Find network sources generating a high number of requests, helping detect possible denial-of-service (DoS) attacks or heavy application usage.

```sql
select
  start_time,
  src_addr,
  count(*) as request_count,
  date_trunc('minute', start_time) as request_minute
from
  aws_vpc_flow_log
group by
  start_time, src_addr, request_minute
having
  count(*) > 100
order by
  request_count desc;
```

```yaml
folder: VPC
```

### High-Volume Rejected Traffic

Identify network sources generating a large number of rejected requests. This helps detect access control violations or attack attempts.

```sql
select
  start_time,
  src_addr,
  count(*) as rejected_requests
from
  aws_vpc_flow_log
where
  action = 'REJECT'
group by
  start_time, src_addr
having
  count(*) > 50
order by
  rejected_requests desc;
```

```yaml
folder: VPC
```

## Baseline Examples

### Traffic Outside Standard Business Hours

Identify network activity occurring outside standard working hours (e.g., between 8 PM and 6 AM).

```sql
select
  start_time,
  src_addr,
  dst_addr,
  src_port,
  dst_port,
  protocol,
  vpc_id
from
  aws_vpc_flow_log
where
  extract('hour' from start_time) >= 20 -- 8 PM
  or extract('hour' from start_time) < 6 -- 6 AM
order by
  start_time desc;
```

```yaml
folder: VPC
```

### Track Traffic to a Specific Instance

Retrieve all network traffic related to a particular EC2 instance.

```sql
select
  start_time,
  src_addr,
  dst_addr,
  src_port,
  dst_port,
  protocol,
  instance_id,
  action
from 
  aws_vpc_flow_log
where
  instance_id = 'i-085c7a43a498c2f5d'
order by 
  start_time desc;
```

```yaml
folder: VPC
```
