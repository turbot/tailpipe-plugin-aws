## Activity Examples

### Daily access trends

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

### Top 10 accessed objects

List the 10 most frequently accessed IAM objects.

```sql
select
bucket,
key,
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

### Top 10 requester IP addresses

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
