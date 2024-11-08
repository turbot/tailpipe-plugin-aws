
# AWS Security Threat Detection Queries

## Detect AWS Console unsuccessful login attempts
Flags actors who failed to login to the console.

```sql
select
  epoch_ms(event_time) as event_time,
  event_name,
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region,
  recipient_account_id as account_id,
  user_agent,
  additional_event_data,
  request_parameters,
  response_elements,
  service_event_details,
  resources,
  user_identity
from
  aws_cloudtrail_log
where
  event_source = 'signin.amazonaws.com'
  and event_name = 'ConsoleLogin'
  and (additional_event_data::json ->> 'MFAUsed') = 'Yes'
  and (response_elements::json ->> 'ConsoleLogin') = 'Failure'
order by
  event_time desc;
```

## Detect Access from Non-Whitelisted Locations
Flags actors accessing from unapproved locations, which may indicate unauthorized access or policy violations.

```sql
select
  user_identity.arn as user_arn,
  source_ip_address,
  aws_region
from
  aws_cloudtrail_log
where
  source_ip_address not in ('trusted_location_1', 'trusted_location_2')
group by
  user_identity.arn, source_ip_address, aws_region;
```

## Detect high volume of S3 bucket access in short time
Flags accounts with a high number of S3 access events within a short period, indicating potential brute force or data exfiltration attempts.

```sql
select
  user_identity.arn as user_arn,
  count(*) as access_count,
  date_trunc('minute', event_time) as event_minute
from
  aws_cloudtrail_log
where
  event_source = 's3.amazonaws.com'
  and event_name in ('GetObject', 'PutObject', 'ListBucket')
group by
  user_identity.arn, event_minute
having
  count(*) > 100
order by
  access_count desc;
```

## Detect attempts to disable CloudTrail logging
Detects if any actor attempted to stop CloudTrail logging, which may indicate an attempt to cover tracks.

```sql
select
  user_identity.arn as user_arn,
  event_name,
  event_time
from
  aws_cloudtrail_log
where
  event_name in ('StopLogging', 'DeleteTrail')
  and event_source = 'cloudtrail.amazonaws.com'
order by
  event_time desc;
```

## Detect IAM policy modifications
Flags actors modifying IAM policies, which may indicate privilege escalation attempts.

```sql
select
  user_identity.arn as user_arn,
  event_name,
  request_parameters,
  event_time
from
  aws_cloudtrail_log
where
  event_source = 'iam.amazonaws.com'
  and event_name in ('PutUserPolicy', 'AttachUserPolicy', 'DetachUserPolicy')
order by
  event_time desc;
```

## Detect unusual source IP addresses for users
Flags users accessing from unfamiliar IP addresses, which may indicate account compromise.

```sql
select
  user_identity.arn as user_arn,
  source_ip_address,
  count(*) as access_count,
  date_trunc('day', event_time) as access_day
from
  aws_cloudtrail_log
where
  source_ip_address not in (select distinct source_ip_address from aws_cloudtrail_log)
group by
  user_identity.arn, source_ip_address, access_day
having
  access_count > 5
order by
  access_count desc;
```

## Detect root account activity
Detects usage of the root account, which is generally a high-risk action.

```sql
select
  event_time,
  event_name,
  source_ip_address,
  aws_region,
  user_identity.arn as user_arn
from
  aws_cloudtrail_log
where
  user_identity.type = 'Root'
order by
  event_time desc;
```

## Detect EC2 instance start or stop events
Flags EC2 instance start or stop actions, which could indicate unauthorized resource usage.

```sql
select
  user_identity.arn as user_arn,
  event_name,
  request_parameters,
  event_time
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name in ('StartInstances', 'StopInstances')
order by
  event_time desc;
```

## Detect IAM user deletions
Flags deletion of IAM users, which could indicate malicious attempts to remove evidence.

```sql
select
  user_identity.arn as user_arn,
  event_name,
  request_parameters,
  event_time
from
  aws_cloudtrail_log
where
  event_source = 'iam.amazonaws.com'
  and event_name = 'DeleteUser'
order by
  event_time desc;
```

## Detect excessive role assumption
Flags actors assuming roles more frequently than usual, indicating potential privilege misuse.

```sql
select
  user_identity.arn as user_arn,
  count(*) as assumption_count,
  date_trunc('hour', event_time) as assumption_hour
from
  aws_cloudtrail_log
where
  event_source = 'sts.amazonaws.com'
  and event_name = 'AssumeRole'
group by
  user_identity.arn, assumption_hour
having
  count(*) > 10
order by
  assumption_count desc;
```

