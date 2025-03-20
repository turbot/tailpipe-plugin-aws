## Activity Examples

### Daily Connection Trends
Count connections per day to identify traffic patterns over time. This query provides a comprehensive view of daily connection volume, helping you understand usage patterns, peak hours, and potential seasonal variations in network traffic.

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

```yaml
folder: ELB
```

### Top 10 Clients by Connection Count
List the top 10 client IP addresses making connections. This query helps identify the most active clients, potential sources of high traffic, and can assist in network security monitoring and capacity planning.

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

```yaml
folder: ELB
```

### Connection Distribution by Destination
Analyze how connections are distributed across destination instances. Understanding destination-level connection patterns can help optimize network configuration, identify potential bottlenecks, and ensure balanced traffic across different backend resources.

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

```yaml
folder: ELB
```

### TLS Protocol Version Distribution
Analyze the distribution of TLS protocol versions used. This query provides insights into the security and encryption standards of incoming connections, helping identify potential security upgrades or legacy system interactions.

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

```yaml
folder: ELB
```

## Detection Examples

### Failed TLS Handshakes
Identify instances where TLS handshakes failed or encountered alerts. This query helps detect potential security issues, misconfigured clients, or network problems that prevent successful encrypted connections.

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

```yaml
folder: ELB
```

### TLS Cipher Vulnerabilities
Detect usage of deprecated or insecure TLS ciphers. This query helps identify outdated SSL/TLS protocols that may pose security risks, allowing you to upgrade and maintain robust encryption standards.

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

```yaml
folder: ELB
```

## Operational Examples

### Slow Connection Times
Top 10 connections with unusually high connection establishment times. This query helps identify performance bottlenecks in connection initialization, which can impact overall network responsiveness and user experience.

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

```yaml
folder: ELB
```

### TLS Handshake Performance
Identify TLS handshakes that took an unusually long time. This query helps detect potential cryptographic performance issues or misconfigured security settings that could slow down connection establishment.

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

```yaml
folder: ELB
```

## Volume Examples

### High Traffic Periods
Detect periods of unusually high connection volume. This query helps identify peak traffic times, potential Denial of Service (DoS) attacks, or unexpected usage patterns that might require infrastructure scaling.

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

```yaml
folder: ELB
```

### Large Data Transfers
Track connections transferring unusually large amounts of data. This query helps identify potential data exfiltration, backup processes, or unusual data transfer patterns that might impact network performance.

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

```yaml
folder: ELB
```

### TLS Cipher Usage by Protocol
Analyze the distribution of TLS ciphers across different protocol versions. This query provides detailed insights into encryption method diversity and potential security improvements within specific TLS protocol versions.

```sql
select
  tls_protocol_version,
  tls_cipher,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (partition by tls_protocol_version), 3) as percentage_within_protocol
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

```yaml
folder: ELB
```

### ALPN Protocol Distribution
Analyze the distribution of ALPN (Application-Layer Protocol Negotiation) protocols used in frontend connections. This query helps understand application-level protocol preferences and potential modernization opportunities.

```sql
select
  alpn_fe_protocol,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 3) as percentage
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

```yaml
folder: ELB
```

### TLS Named Group Distribution
Analyze the distribution of TLS named groups used in connections. This query provides insights into the cryptographic key exchange methods and elliptic curve preferences in your network connections.

```sql
select
  tls_named_group,
  count(*) as connection_count,
  round(count(*) * 100.0 / sum(count(*)) over (), 3) as percentage
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

```yaml
folder: ELB
```

### Certificate Usage Analysis
Analyze which certificates are being used for TLS connections. This query helps track certificate utilization, identify potential consolidation opportunities, and monitor certificate lifecycle.

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

```yaml
folder: ELB
```

### Domain Name Analysis
Analyze which domain names are being requested. This query provides insights into the types of services and endpoints being accessed through your network load balancer.

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

```yaml
folder: ELB
```

### Connection Time Distribution
Analyze the distribution of connection establishment times. This query helps understand the overall performance characteristics of network connections, identifying potential latency issues.

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

```yaml
folder: ELB
```