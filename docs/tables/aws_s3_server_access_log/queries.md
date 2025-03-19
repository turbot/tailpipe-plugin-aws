## Activity Examples

### Daily Access Trends

Count access log entries per day to identify trends over time.

```sql
select
  strftime(timestamp, '%Y-%m-%d') as access_date,
  count(*) AS requests
from
  aws_s3_server_access_log
group by
  access_date
order by
  access_date asc;
```

```yaml
folder: S3
```

### Top 10 Accessed Objects

List the 10 most frequently accessed S3 objects.

```sql
select
  bucket,
  key,
  count(*) as requests
from
  aws_s3_server_access_log
where
  key is not null
group by
  bucket,
  key
order by
  requests desc
limit 10;
```

```yaml
folder: S3
```

### Top 10 Requester IP Addresses

List the top 10 requester IP addresses.

```sql
select
  remote_ip,
  count(*) as request_count
from
  aws_s3_server_access_log
group by
  remote_ip
order by
  request_count desc
limit 10;
```

```yaml
folder: S3
```

### Top Error Codes

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

```yaml
folder: S3
```

## Detection Examples

### Unusually Large File Downloads

Detect unusually large downloads based on file size.

```sql
select
  timestamp,
  bucket,
  key,
  bytes_sent,
  operation,
  request_uri,
  requester,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  bytes_sent > 50000000 -- 50MB
  and http_status = 200
order by
  bytes_sent desc;
```

```yaml
folder: S3
```

### Requests from Unapproved IAM Roles and Users

Flag requests from IAM roles and users outside an approved list (by AWS account ID in this example).

```sql
select
  timestamp,
  bucket,
  operation,
  requester,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  requester is not null -- Exclude unauthenticated requests
  and requester not like 'arn:%:%:%:123456789012:%'
order by
  timestamp desc;
```

```yaml
folder: S3
```

## Operational Examples

### Failed Object Upload Requests

List failed object upload requests along with the error codes.

```sql
select
  timestamp,
  bucket,
  key,
  requester,
  remote_ip,
  http_status,
  error_code
from
  aws_s3_server_access_log
where
  operation = 'REST.PUT.OBJECT'
  and http_status >= 400
order by
  timestamp desc;
```

```yaml
folder: S3
```

### Unauthenticated Requests

List all unauthenticated requests.

```sql
select
  timestamp,
  bucket,
  operation,
  request_uri,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  requester is null
order by
  timestamp desc;
```

```yaml
folder: S3
```

## Volume Examples

### High Volume of Requests

Detect unusually high number of requests by remote IP address.

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

```yaml
folder: S3
```

### High Volume of Failed Requests

Identify remote IPs with a high number of failed requests.

```sql
select
  remote_ip,
  bucket,
  count(*) as failed_requests
from
  aws_s3_server_access_log
where
  http_status >= 400
group by
  remote_ip,
  bucket
having
  count(*) > 50
order by
  failed_requests desc;
```

```yaml
folder: S3
```

## Baseline Examples

### Requests Outside of Normal Hours

Flag requests occurring outside of standard working hours, e.g., between 8 PM and 6 AM.

```sql
select
  timestamp,
  bucket,
  operation,
  requester,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  extract('hour' from timestamp) >= 20 -- 8 PM
  or extract('hour' from timestamp) < 6 -- 6 AM
order by
  timestamp desc;
```

```yaml
folder: S3
```
