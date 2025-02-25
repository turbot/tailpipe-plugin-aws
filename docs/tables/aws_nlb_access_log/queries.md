## Activity Examples

### Daily connection trends

Count connections per day to identify traffic patterns over time.

```sql
select
  strftime(timestamp, '%Y-%m-%d') as connection_date,
  count(*) as connection_count
from
  aws_nlb_access_log
group by
  connection_date
order by
  connection_date asc;
```

### Top 10 clients by connection count

List the top 10 client IP addresses making connections.

```sql
select
  client_ip,
  count(*) as connection_count
from
  aws_nlb_access_log
group by
  client_ip
order by
  connection_count desc
limit 10;
```

### Connection distribution by destination

Analyze how connections are distributed across destination instances.

```sql
select
  destination_ip,
  count(*) as connection_count
from
  aws_nlb_access_log
group by
  destination_ip
order by
  connection_count desc;
```

### TLS protocol version distribution

Analyze the distribution of TLS protocol versions used.

```sql
select
  tls_protocol_version,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 2) as percentage
from
  aws_nlb_access_log
where
  tls_protocol_version is not null
group by
  tls_protocol_version
order by
  connection_count desc;
```

## Detection Examples

### Failed TLS handshakes

Identify instances where TLS handshakes failed or encountered alerts.

```sql
select
  timestamp,
  elb,
  client_ip,
  client_port,
  incoming_tls_alert,
  tls_protocol_version,
  tls_cipher
from
  aws_nlb_access_log
where
  incoming_tls_alert is not null
  and incoming_tls_alert != ''
order by
  timestamp desc;
```

### TLS cipher vulnerabilities

Detect usage of deprecated or insecure TLS ciphers.

```sql
select
  tls_cipher,
  tls_protocol_version,
  count(*) as connection_count
from
  aws_nlb_access_log
where
  tls_protocol_version in ('tlsv1.1', 'tlsv1', 'sslv3', 'sslv2') -- Insecure protocols
group by
  tls_cipher,
  tls_protocol_version
order by
  connection_count desc;
```

## Operational Examples

### Slow connection times

Top 10 connections with unusually high connection establishment times.

```sql
select
  timestamp,
  elb,
  client_ip,
  client_port,
  destination_ip,
  destination_port,
  connection_time, -- Time taken to establish connection in ms
  tls_handshake_time -- Time taken for TLS handshake in ms
from
  aws_nlb_access_log
where
  connection_time > 1000 -- Connections taking longer than 1 second
order by
  connection_time desc
limit 10;
```

### TLS handshake performance

Identify TLS handshakes that took unusually long time.

```sql
select
  timestamp,
  elb,
  client_ip,
  client_port,
  tls_protocol_version,
  tls_cipher,
  tls_handshake_time -- Time taken for TLS handshake in ms
from
  aws_nlb_access_log
where
  tls_handshake_time > 500 -- TLS handshakes taking longer than 500ms
order by
  tls_handshake_time desc
limit 10;
```

## Volume Examples

### High traffic periods

Detect periods of unusually high connection volume.

```sql
select
  date_trunc('minute', timestamp) as connection_minute,
  elb,
  count(*) as connection_count
from
  aws_nlb_access_log
group by
  connection_minute,
  elb
having
  count(*) > 1000
order by
  connection_count desc;
```

### Large data transfers

Track connections transferring unusually large amounts of data.

```sql
select
  timestamp,
  elb,
  client_ip,
  client_port,
  destination_ip,
  destination_port,
  sent_bytes,
  received_bytes,
  (sent_bytes + received_bytes) as total_bytes
from
  aws_nlb_access_log
where
  (sent_bytes + received_bytes) > 10485760 -- 10MB
order by
  total_bytes desc;
```

### TLS cipher usage by protocol

Analyze the distribution of TLS ciphers across different protocol versions.

```sql
select
  tls_protocol_version,
  tls_cipher,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (partition by tls_protocol_version), 2) as percentage_within_protocol
from
  aws_nlb_access_log
where
  tls_protocol_version is not null
  and tls_cipher is not null
group by
  tls_protocol_version,
  tls_cipher
order by
  tls_protocol_version,
  connection_count desc;
```

### ALPN protocol distribution

Analyze the distribution of ALPN protocols used in frontend connections.

```sql
select
  alpn_fe_protocol,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 2) as percentage
from
  aws_nlb_access_log
where
  alpn_fe_protocol is not null
  and alpn_fe_protocol != ''
group by
  alpn_fe_protocol
order by
  connection_count desc;
```

### TLS named group distribution

Analyze the distribution of TLS named groups used in connections.

```sql
select
  tls_named_group,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 2) as percentage
from
  aws_nlb_access_log
where
  tls_named_group is not null
  and tls_named_group != ''
group by
  tls_named_group
order by
  connection_count desc;
```

### Certificate usage analysis

Analyze which certificates are being used for TLS connections.

```sql
select
  chosen_cert_arn,
  count(*) as connection_count
from
  aws_nlb_access_log
where
  chosen_cert_arn is not null
  and chosen_cert_arn != ''
group by
  chosen_cert_arn
order by
  connection_count desc;
```

### Domain name analysis

Analyze which domain names are being requested.

```sql
select
  domain_name,
  count(*) as connection_count
from
  aws_nlb_access_log
where
  domain_name is not null
  and domain_name != ''
group by
  domain_name
order by
  connection_count desc;
```

### Connection time distribution

Analyze the distribution of connection establishment times.

```sql
select
  case
    when connection_time < 10 then '0-10ms'
    when connection_time < 50 then '10-50ms'
    when connection_time < 100 then '50-100ms'
    when connection_time < 200 then '100-200ms'
    when connection_time < 500 then '200-500ms'
    when connection_time < 1000 then '500-1000ms'
    else '1000ms+'
  end as connection_time_range,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 2) as percentage
from
  aws_nlb_access_log
group by
  connection_time_range
order by
  case connection_time_range
    when '0-10ms' then 1
    when '10-50ms' then 2
    when '50-100ms' then 3
    when '100-200ms' then 4
    when '200-500ms' then 5
    when '500-1000ms' then 6
    else 7
  end;
```