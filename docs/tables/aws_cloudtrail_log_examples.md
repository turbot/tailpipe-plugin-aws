## Count Examples

### Count events per day

Return a count of events per day to see trends across time.

```sql
select
  strftime(event_time, '%Y-%m-%d') AS event_date,
  count(*) AS event_count
from
  aws_cloudtrail_log
group by
  event_date
order by
  event_date ASC;
```

### Count events by event name

Return a count of events by event name see the most operations.

```sql
select
  event_source,
  event_name,
  count(*) as event_count
from
  aws_cloudtrail_log
group by
  event_source,
  event_name,
order by
  event_count desc;
```

### Count events by event name (exclude read-only events)

```sql
select
  event_source,
  event_name,
  count(*) as event_count
from
  aws_cloudtrail_log
where
  not read_only
group by
  event_source,
  event_name,
order by
  event_count desc;
```

### Count events grouped by account ID

```sql
select
  event_source,
  event_name,
  recipient_account_id,
  count(*) as event_count
from
  aws_cloudtrail_log
group by
  event_source,
  event_name,
  recipient_account_id
order by
  event_count desc;
```

### Count errors by frequency

```sql
select
  error_code,
  count(*) as event_count
from
  aws_cloudtrail_log
where
  error_code is not null
group by
  error_code
order by
  event_count desc;
```

## Detection Examples

### Detect default EBS encryption disabled in a region

```sql
select
  event_time,
  event_source,
  event_name,
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name = 'DisableEbsEncryptionByDefault'
  and error_code is null
order by
  event_time desc;
```

### Detect attempts to stop CloudTrail trail logging

```sql
select
  event_time,
  event_source,
  event_name,
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  event_name in ('StopLogging', 'DeleteTrail')
  and event_source = 'cloudtrail.amazonaws.com'
  and error_code is not null
order by
  event_time desc;
```

### Detect unsuccessful AWS console login attempts

```sql
select
  event_time,
  event_name,
  user_identity.arn as user_arn,
  source_ip_address,
  user_agent,
  aws_region,
  recipient_account_id as account_id,
  additional_event_data
from
  aws_cloudtrail_log
where
  event_source = 'signin.amazonaws.com'
  and event_name = 'ConsoleLogin'
  and error_code is not null
order by
  event_time desc;
```

### Detect root account activity

```sql
select
  event_time,
  event_name,
  source_ip_address,
  user_agent,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  user_identity.type = 'Root'
order by
  event_time desc;
```

### Detect actions in unapproved regions

```sql
select
  event_time,
  event_source,
  event_name,
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  aws_region not in ('us-east-1', 'us-west-1')
order by
  event_time desc;
```

### Detect actions from unapproved IP addresses

```sql
select
  event_time,
  event_source,
  event_name,
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  source_ip_address not in ('192.0.2.146', '206.253.208.100')
order by
  event_time desc;
```

## Operational Examples

### List VPC security group rule updates

```sql
select
  event_time,
  event_source,
  event_name,
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name in ('AuthorizeSecurityGroupEgress', 'AuthorizeSecurityGroupIngress', 'RevokeSecurityGroupEgress', 'RevokeSecurityGroupIngress')
  and error_code is null
order by
  event_time desc;
```

### List IAM user policy updates

```sql
select
  event_time,
  event_source,
  event_name,
  request_parameters,
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  event_source = 'iam.amazonaws.com'
  and event_name in ('AttachUserPolicy', 'DetachUserPolicy', 'PutUserPolicy')
order by
  event_time desc;
```

## Volume Examples

### Detect high volume of S3 bucket access requests

```sql
select
  user_identity.arn as user_arn,
  count(*) as event_count,
  date_trunc('minute', event_time) as event_minute
from
  aws_cloudtrail_log
where
  event_source = 's3.amazonaws.com'
  and event_name in ('GetObject', 'PutObject', 'ListBucket')
group by
  user_identity.arn,
  event_minute
having
  count(*) > 100
order by
  event_count desc;
```

### Detect excessive IAM role assumption

```sql
select
  user_identity.arn as user_arn,
  count(*) as event_count,
  date_trunc('hour', event_time) as event_hour
from
  aws_cloudtrail_log
where
  event_source = 'sts.amazonaws.com'
  and event_name = 'AssumeRole'
group by
  user_identity.arn,
  event_hour
having
  count(*) > 10
order by
  event_hour desc,
  event_count desc;
```

## Baseline Examples

### Detect unusual source IP addresses for users

```sql
select
  user_identity.arn as user_arn,
  source_ip_address,
  count(*) as access_count,
  date_trunc('day', event_time) as access_day
from
  aws_cloudtrail_log
where
  source_ip_address not like '%.amazonaws.com'
  and source_ip_address not in (select distinct source_ip_address from aws_cloudtrail_log)
group by
  user_identity.arn, source_ip_address, access_day
having
  access_count > 5
order by
  access_count desc;
```

### Detect activity outside of normal hours (between 8 PM and 6 AM)

```sql
select
  event_time,
  event_source,
  event_name,
  request_parameters,
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region,
  recipient_account_id as account_id
from
  aws_cloudtrail_log
where
  cast(strftime(event_time, '%H') as integer) >= 20
  or cast(strftime(event_time, '%H') as integer) < 6
order by
  event_time desc;
