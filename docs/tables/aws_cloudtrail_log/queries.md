## Activity Examples

### Daily activity trends

Count events per day to identify activity trends over time.

```sql
select
  strftime(event_time, '%y-%m-%d') AS event_date,
  count(*) AS event_count
from
  aws_cloudtrail_log
group by
  event_date
order by
  event_date asc;
```

### Top 10 events

List the 10 most frequently called events.

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
  event_count desc
limit 10;
```

### Top 10 events (exclude read-only)

List the top 10 most frequently called events, excluding read-only events.

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
  event_count desc
limit 10;
```

### Top events by account

Count and group events by account ID, event source, and event name to analyze activity across accounts.

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

### Top error codes

Identify the most frequent error codes.

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

### Default EBS encryption disabled in a region

Detect when default EBS encryption was disabled in a region.

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

### CloudTrail trail logging stopped

Detect when logging was stopped for a CloudTrail trail.

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
  event_source = 'cloudtrail.amazonaws.com'
  and event_name in ('StopLogging', 'DeleteTrail')
  and error_code is not null
order by
  event_time desc;
```

### Unsuccessful AWS console login attempts

Find failed console login attempts, highlighting potential unauthorized access attempts.

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

### Root activity

Track any actions performed by the root user.

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

### Activity in unapproved regions

Identify actions occurring in AWS regions outside an approved list.

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

### Activity from unapproved IP addresses

Flag activity originating from IP addresses outside an approved list.

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

### VPC security group rule updates

Track changes to VPC security group ingress and egress rules.

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

### IAM user permission updates

List events where an IAM user has added or removed permissions through managed policies, inline policies, or groups.

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
  and event_name in ('AddUserToGroup', 'AttachUserPolicy', 'DeleteUserPolicy', 'DetachUserPolicy', 'PutUserPolicy', 'RemoveUserFromGroup')
order by
  event_time desc;
```

## Volume Examples

### High volume of S3 bucket access requests

Detect unusually high access activity to S3 buckets and objects.

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

### Excessive IAM role assumptions

Identify IAM roles being assumed at an unusually high frequency.

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

### Unrecognized user source IP addresses

Detect user activity from unexpected or new source IP addresses.

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

### Activity outside of normal hours

Flag activity occurring outside of standard working hours, e.g., activity bewteen 8 PM and 6 AM.

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
  cast(strftime(event_time, '%H') as integer) >= 20 -- 8 PM
  or cast(strftime(event_time, '%H') as integer) < 6 -- 6 AM
order by
  event_time desc;
