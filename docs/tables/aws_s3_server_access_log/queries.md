## Access Log Examples

### Daily access trends

Count access log entries per day to identify trends over time.

```sql
select
  strftime(timestamp, '%y-%m-%d') AS access_date,
  count(*) AS access_count
from
  aws_s3_server_access_log
group by
  access_date
order by
  access_date asc;
```

### Top 10 accessed objects

List the 10 most frequently accessed S3 objects.

```sql
select
  timestamp,
  requester,
  bucket,
  key,
  count(*) as access_count
from
  aws_s3_server_access_log
group by
  timestamp,
  requester,
  bucket,
  key
order by
  access_count desc
limit 10;
```

### Top 10 requesters

Identify the top 10 requesters generating the most traffic.

```sql
select
  timestamp,
  requester,
  count(*) as request_count
from
  aws_s3_server_access_log
group by
  timestamp,
  requester
order by
  request_count desc
limit 10;
```

### Top error codes

Identify the most frequent error codes.

```sql
select
  timestamp,
  requester,
  error_code,
  count(*) as access_count
from
  aws_s3_server_access_log
where
  error_code is not null
group by
  timestamp,
  requester,
  error_code
order by
  access_count desc;
```

## Detection Examples

### Unusual large file downloads

Detect unusually large downloads from S3.

```sql
select
  timestamp,
  requester,
  bucket,
  key,
  bytes_sent,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  bytes_sent is not null
  and bytes_sent > 50000000 -- 50MB
order by
  bytes_sent desc;
```

### Requests from unknown IPs

Identify S3 access from unknown or unapproved IP addresses.

```sql
select
  timestamp,
  requester,
  bucket,
  remote_ip,
  operation
from
  aws_s3_server_access_log
where
  remote_ip not in ('192.0.2.146', '206.253.208.100')
order by
  timestamp desc;
```

## Operational Examples

### Find failed requests

Identify all S3 requests that resulted in an error.

```sql
select
  timestamp,
  requester,
  bucket,
  request_uri,
  http_status,
  error_code,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  http_status is not null
  and http_status >= 400
order by
  timestamp desc;
```

### Detect ACL-related access issues

Find failed requests where ACL permissions were required, helping identify potential permission misconfigurations or access control issues.

```sql
select
  timestamp,
  requester,
  bucket,
  key,
  operation,
  error_code,
  http_status,
  acl_required
from
  aws_s3_server_access_log
where
  acl_required = true
  and (
    http_status >= 400
    or error_code is not null
  )
order by
  timestamp desc;
```

## Volume Examples

### High volume of access requests

Detect unusually high access activity to S3 buckets and objects.

```sql
select
  timestamp,
  requester,
  bucket,
  count(*) as access_count,
  date_trunc('minute', timestamp) as access_minute
from
  aws_s3_server_access_log
group by
  timestamp,
  requester,
  bucket,
  access_minute
having
  count(*) > 100
order by
  access_count desc;
```

### High volume of failed requests

Identify accounts with a high number of failed requests.

```sql
select
  timestamp,
  requester,
  count(*) as failed_requests
from
  aws_s3_server_access_log
where
  http_status is not null
  and http_status >= 400
group by
  timestamp,
  requester
having
  count(*) > 50
order by
  failed_requests desc;
```

## Baseline Examples

### Unrecognized user source IP addresses

Detect user access from unexpected or new source IP addresses.

```sql
select
  timestamp,
  requester,
  remote_ip,
  count(*) as access_count,
  date_trunc('day', timestamp) as access_day
from
  aws_s3_server_access_log
where
  remote_ip not in (select distinct remote_ip from aws_s3_server_access_log)
group by
  timestamp,
  requester,
  remote_ip,
  access_day
having
  access_count > 5
order by
  access_count desc;
```

### Access outside of normal hours

Flag access occurring outside of standard working hours, e.g., between 8 PM and 6 AM.

```sql
select
  timestamp,
  requester,
  bucket,
  remote_ip,
  operation
from
  aws_s3_server_access_log
where
  cast(strftime(timestamp, '%H') as integer) >= 20 -- 8 PM
  or cast(strftime(timestamp, '%H') as integer) < 6 -- 6 AM
order by
  timestamp desc;
```
