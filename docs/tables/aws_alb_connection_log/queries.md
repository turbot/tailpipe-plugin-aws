## Activity Examples

### Daily Connection Trends
Count connections per day to identify traffic patterns over time. This query provides a comprehensive view of daily connection volume, helping you understand usage patterns, peak hours, and potential seasonal variations in network traffic.

```sql
select
  strftime(timestamp, '%Y-%m-%d') as connection_date,
  count(*) as connection_count
from
  aws_alb_connection_log
group by
  connection_date
order by
  connection_date asc;
```

```yaml
folder: ELB
```

### Top 10 Clients by Connection Count
List the top 10 client IP addresses making connection attempts. This query helps identify the most active clients, potential sources of high traffic, and can assist in network security monitoring and capacity planning.

```sql
select
  client_ip,
  count(*) as connection_count
from
  aws_alb_connection_log
group by
  client_ip
order by
  connection_count desc
limit 10;
```

```yaml
folder: ELB
```

### Connection Distribution by Listener Port
Analyze how connections are distributed across listener ports. Understanding port-level connection patterns can help optimize network configuration, identify potential bottlenecks, and ensure balanced traffic across different services.

```sql
select
  listener_port,
  count(*) as connection_count
from
  aws_alb_connection_log
group by
  listener_port
order by
  connection_count desc;
```

```yaml
folder: ELB
```

### TLS Protocol Distribution
Analyze the distribution of TLS protocols used by clients. This query provides insights into the security and encryption standards of incoming connections, helping identify potential security upgrades or legacy system interactions.

```sql
select
  tls_protocol,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 2) as percentage
from
  aws_alb_connection_log
where
  tls_protocol is not null
group by
  tls_protocol
order by
  connection_count desc;
```

```yaml
folder: ELB
```

## Detection Examples

### Failed TLS Handshakes
Identify connections with TLS handshake verification failures. This query helps detect potential security issues, misconfigured clients, or network problems that prevent successful encrypted connections.

```sql
select
  timestamp,
  tp_index as client_ip,
  conn_trace_id,
  client_port,
  tls_protocol,
  tls_cipher,
  tls_verify_status
from
  aws_alb_connection_log
where
  tls_verify_status like 'Failed:%'
order by
  timestamp desc;
```

```yaml
folder: ELB
```

### Deprecated TLS Protocols
Detect usage of deprecated or insecure TLS protocols. This query helps identify outdated SSL/TLS protocols that may pose security risks, allowing you to upgrade and maintain robust encryption standards.

```sql
select
  tls_protocol,
  tls_cipher,
  count(*) as connection_count
from
  aws_alb_connection_log
where
  tls_protocol in ('TLSv1.1', 'TLSv1', 'SSLv3', 'SSLv2') -- Insecure protocols
group by
  tls_protocol,
  tls_cipher
order by
  connection_count desc;
```

```yaml
folder: ELB
```

## Operational Examples

### Slow TLS Handshakes
Top 10 connections with unusually high TLS handshake latency. This query helps identify performance bottlenecks in the TLS negotiation process, which can impact overall connection establishment times and user experience.

```sql
select
  timestamp,
  tp_index as client_ip,
  conn_trace_id,
  client_port,
  tls_protocol,
  tls_cipher,
  tls_handshake_latency
from
  aws_alb_connection_log
where
  tls_handshake_latency > 1 -- Handshakes taking longer than 1 second
order by
  tls_handshake_latency desc
limit 10;
```

```yaml
folder: ELB
```

## Volume Examples

### TLS Cipher Usage
Analyze the distribution of TLS ciphers used by connections. This query provides detailed insights into the encryption methods clients are using, helping assess cryptographic diversity and potential security improvements.

```sql
select
  tls_cipher,
  tls_protocol,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 3) as percentage
from
  aws_alb_connection_log
group by
  tls_cipher,
  tls_protocol
order by
  connection_count desc;
```

```yaml
folder: ELB
```

### Connection Failure Rate by Time Period
Analyze the rate of connection failures over time. This query helps identify temporal patterns in connection failures, potentially revealing systemic issues, network problems, or security-related connection challenges.

```sql
select
  strftime(timestamp, '%Y-%m-%d %H:00:00') as hour,
  count(*) as total_connections,
  sum(case when tls_verify_status like 'Failed:%' then 1 else 0 end) as failed_connections,
  round(sum(case when tls_verify_status like 'Failed:%' then 1 else 0 end) * 100.0 / count(*), 2) as failure_rate
from
  aws_alb_connection_log
group by
  hour
order by
  hour desc;
```

```yaml
folder: ELB
```

### Connection Trace Correlation
Link connection logs to access logs using the connection trace ID. This query enables deep investigation of connection lifecycle by correlating low-level connection details with HTTP access information.

```sql
select
  c.timestamp as connection_timestamp,
  c.client_ip,
  c.client_port,
  c.tls_protocol,
  c.tls_handshake_latency,
  a.timestamp as access_timestamp,
  a.request_url,
  a.request_http_method,
  a.request_http_version,
  a.elb_status_code
from
  aws_alb_connection_log c
join
  aws_alb_access_log a
on
  c.conn_trace_id = a.conn_trace_id
order by
  c.timestamp desc
limit 10;
```

```yaml
folder: ELB
```