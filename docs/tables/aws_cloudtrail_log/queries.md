---
title: "Sample Queries: aws_cloudtrail_log - Query AWS CloudTrail Logs"
description: "Sample queries for AWS CloudTrail logs to help with security analysis, resource tracking, and compliance auditing."
---

# Sample Queries: aws_cloudtrail_log

## Security Analysis

### Unauthorized API calls

Find unauthorized API calls to investigate potential security breaches:

```sql
select
  event_time,
  event_source,
  event_name,
  error_code,
  error_message,
  source_ip_address,
  user_identity->>'userName' as user_name,
  aws_region
from
  aws_cloudtrail_log
where
  error_code = 'UnauthorizedOperation'
  or error_code = 'AccessDenied'
order by
  event_time desc;
```

### Console login activity

Monitor AWS Management Console login events:

```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  source_ip_address,
  aws_region,
  user_agent,
  error_message
from
  aws_cloudtrail_log
where
  event_name = 'ConsoleLogin'
order by
  event_time desc;
```

### List all failed console logins
```sql
select
  event_time,
  source_ip_address,
  user_identity->>'userName' as user_name,
  error_code,
  error_message
from
  aws_cloudtrail_log
where
  event_name = 'ConsoleLogin'
  and error_code is not null
order by
  event_time desc;
```

### Find root account usage
```sql
select
  event_time,
  event_name,
  event_source,
  source_ip_address,
  user_identity->>'type' as identity_type
from
  aws_cloudtrail_log
where
  user_identity->>'type' = 'Root'
order by
  event_time desc;
```

### Track IAM policy changes
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'policyName' as policy_name,
  request_parameters->>'policyDocument' as policy_document
from
  aws_cloudtrail_log
where
  event_source = 'iam.amazonaws.com'
  and event_name like '%Policy%'
order by
  event_time desc;
```

### Monitor security group changes
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'groupId' as security_group_id,
  request_parameters->>'ipPermissions' as ip_permissions
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name like '%SecurityGroup%'
order by
  event_time desc;
```

## Resource Changes

### EC2 instance changes

Track EC2 instance lifecycle events:

```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'instanceId' as instance_id,
  source_ip_address,
  aws_region
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name in (
    'RunInstances',
    'StartInstances',
    'StopInstances',
    'TerminateInstances'
  )
order by
  event_time desc;
```

### S3 bucket policy changes

Monitor changes to S3 bucket policies:

```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'bucketName' as bucket_name,
  source_ip_address,
  aws_region
from
  aws_cloudtrail_log
where
  event_source = 's3.amazonaws.com'
  and event_name in (
    'PutBucketPolicy',
    'DeleteBucketPolicy',
    'PutBucketAcl'
  )
order by
  event_time desc;
```

### Monitor S3 bucket configuration changes
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'bucketName' as bucket_name
from
  aws_cloudtrail_log
where
  event_source = 's3.amazonaws.com'
  and event_name like '%Bucket%'
order by
  event_time desc;
```

### Track EC2 instance state changes
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'instanceId' as instance_id
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name in ('StartInstances', 'StopInstances', 'TerminateInstances')
order by
  event_time desc;
```

### List resource deletions
```sql
select
  event_time,
  event_name,
  event_source,
  user_identity->>'userName' as user_name,
  resources[0]->>'resourceType' as resource_type,
  resources[0]->>'resourceName' as resource_name
from
  aws_cloudtrail_log
where
  event_name like 'Delete%'
order by
  event_time desc;
```

## Compliance Monitoring

### IAM policy changes

Track changes to IAM policies:

```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'policyName' as policy_name,
  source_ip_address,
  aws_region
from
  aws_cloudtrail_log
where
  event_source = 'iam.amazonaws.com'
  and event_name like '%Policy%'
order by
  event_time desc;
```

### Track changes to security groups
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'groupId' as security_group_id
from
  aws_cloudtrail_log
where
  event_source = 'ec2.amazonaws.com'
  and event_name like '%SecurityGroup%'
order by
  event_time desc;
```

### Monitor KMS key usage
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'keyId' as key_id
from
  aws_cloudtrail_log
where
  event_source = 'kms.amazonaws.com'
order by
  event_time desc;
```

### List changes to CloudWatch alarms
```sql
select
  event_time,
  event_name,
  user_identity->>'userName' as user_name,
  request_parameters->>'alarmName' as alarm_name
from
  aws_cloudtrail_log
where
  event_source = 'monitoring.amazonaws.com'
  and event_name like '%Alarm%'
order by
  event_time desc;
```

## User Activity

### Find API calls from specific IP addresses
```sql
select
  event_time,
  event_name,
  event_source,
  source_ip_address,
  user_identity->>'userName' as user_name
from
  aws_cloudtrail_log
where
  source_ip_address = '192.0.2.1'
order by
  event_time desc;
```

### Track user activity by service
```sql
select
  event_source,
  count(*) as event_count,
  count(distinct user_identity->>'userName') as unique_users
from
  aws_cloudtrail_log
where
  tp_date >= current_date - interval '7 days'
group by
  event_source
order by
  event_count desc;
```

### List read-only vs write operations
```sql
select
  read_only,
  count(*) as event_count
from
  aws_cloudtrail_log
where
  tp_date >= current_date - interval '7 days'
group by
  read_only
order by
  event_count desc;
``` 