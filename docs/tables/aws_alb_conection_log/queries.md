## Activity Examples

### Daily connection trends

Count connections per day to identify traffic patterns over time.

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

### Top 10 clients by connection count

List the top 10 client IP addresses making connection attempts.

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

### Connection distribution by listener port

Analyze how connections are distributed across listener ports.

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

### TLS protocol distribution

Analyze the distribution of TLS protocols used by clients.

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

## Detection Examples

### Failed TLS handshakes

Identify connections with TLS handshake verification failures.

```sql
select
  timestamp,
  tp_index as conn_trace_id,
  client_ip,
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

### Deprecated TLS protocols

Detect usage of deprecated or insecure TLS protocols.

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

## Operational Examples

### Slow TLS handshakes

Top 10 connections with unusually high TLS handshake latency.

```sql
select
  timestamp,
  tp_index as conn_trace_id,
  client_ip,
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

## Volume Examples

### TLS cipher usage

Analyze the distribution of TLS ciphers used by connections.

```sql
select
  tls_cipher,
  tls_protocol,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 2) as percentage
from
  aws_alb_connection_log
group by
  tls_cipher,
  tls_protocol
order by
  connection_count desc;
```

### Connection failure rate by time period

Analyze the rate of connection failures over time.

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

### Connection trace correlation

Link connection logs to access logs using the connection trace ID.

```sql
select
  c.timestamp as connection_timestamp,
  c.client_ip,
  c.client_port,
  c.tls_protocol,
  c.tls_handshake_latency,
  a.timestamp as access_timestamp,
  a.request,
  a.elb_status_code
from
  aws_alb_connection_log c
join
  aws_alb_access_log a
on
  c.tp_index = a.conn_trace_id
order by
  c.timestamp desc
limit 10;
```