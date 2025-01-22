---
title: "Tailpipe Table: aws_cloudtrail_log - Query AWS CloudTrail Logs"
description: "Allows users to query AWS CloudTrail logs."
---

# Table: aws_cloudtrail_log - Query AWS CloudTrail logs

The `aws_cloudtrail_log` table allows you to query data from AWS CloudTrail logs. This table provides detailed information about API calls made within your AWS account, including the event name, source IP address, user identity, and more.

## Queries

For a full list of example queries, please see [aws_cloudtrail_log queries](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_cloudtrail_log).

### Root activity

Find any actions taken by the root user.

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

### Top 10 events

List the top 10 events and how many times they were called.

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

### High volume S3 access requests

Find users generating a high volume of S3 access requests to identify potential anomalous activity.

```sql
select
  user_identity.arn as user_arn,
  count(*) as event_count,
  date_trunc('minute', event_time) as event_minute
from
  aws_cloudtrail_log
where
  event_source = 's3.amazonaws.com'
  and event_name in ('GetObject', 'ListBucket')
group by
  user_identity.arn,
  event_minute
having
  count(*) > 100
order by
  event_count desc;
```
