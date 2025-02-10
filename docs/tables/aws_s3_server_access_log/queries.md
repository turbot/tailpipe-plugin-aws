## Access Log Examples

### Daily access trends

Count access log entries per day to identify trends over time.

```sql
select
  strftime(timestamp, '%y-%m-%d') as date,
  count(*) AS requests
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
  bucket,
  key as object,
  count(*) as access_count
from
  aws_s3_server_access_log
where
  key is not null
group by
  key,
  bucket
order by
  access_count desc
limit 10;
```

### Top 10 requesters

Identify the top 10 requesters generating the most traffic.

```sql
select
  case
    when requester = '-' then 'Unauthenticated'
    else requester
  end as requester,
  count(*) as request_count
from
  aws_s3_server_access_log
group by
  requester
order by
  request_count desc
limit 10;
```

### Top error codes

Identify the most frequent error codes.

```sql
select
  http_status,
  error_code,
  count(*) as error_count
from
  aws_s3_server_access_log
where
  error_code is not null
group by
  http_status,
  error_code
order by
  error_count desc;
```

## Detection Examples

### Unusual large file downloads

Detect unusually large downloads from S3.

```sql
select
  timestamp,
  bucket,
  key,
  bytes_sent,
  requester,
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
  bucket,
  operation,
  requester,
  remote_ip
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
  bucket,
  operation,
  requester,
  remote_ip,
  http_status,
  error_code,
  user_agent
from
  aws_s3_server_access_log
where
  http_status is not null
  and http_status >= 400
order by
  timestamp desc;
```

### Find authenticated requests

Identify all authenticated requests to S3.

```sql
select
  timestamp,
  bucket,
  operation,
  requester,
  remote_ip
from
  aws_s3_server_access_log
where
  requester != '-'
order by
  timestamp desc;
```

## Volume Examples

### High volume of access requests

Detect unusually high access activity to S3 buckets and objects.

```sql
select
  remote_ip,
  bucket,
  count(*) as request_count,
  date_trunc('minute', timestamp) as request_minute
from
  aws_s3_server_access_log
group by
  remote_ip,
  bucket,
  request_minute
having
  count(*) > 100
order by
  request_count desc;
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
