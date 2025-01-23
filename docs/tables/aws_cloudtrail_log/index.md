---
title: "Tailpipe Table: aws_cloudtrail_log - Query AWS CloudTrail Logs"
description: "Allows users to query AWS CloudTrail logs."
---

# Table: aws_cloudtrail_log - Query AWS CloudTrail logs

The `aws_cloudtrail_log` table allows you to query data from AWS CloudTrail logs. This table provides detailed information about API calls made within your AWS account, including the event name, source IP address, user identity, and more.

## Configuration

CloudTrail logs are normally stored within S3 buckets, so a typical configuration would look like:

```hcl
connection "aws" "aws_profile" {
  profile = "my-profile"
}

partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.aws_profile
    bucket     = "aws-cloudtrail-logs-bucket"
  }
}
```

To reduce the amount of logs stored locally, you can use the `filter` argument in your partition:

```hcl
connection "aws" "aws_profile" {
  profile = "my-profile"
}

partition "aws_cloudtrail_log" "my_logs" {
  # Avoid saving read-only events, which often make up to 90% of all log entries
  filter = "not read_only"

  source "aws_s3_bucket" {
    connection = connection.aws.aws_profile
    bucket     = "aws-cloudtrail-logs-bucket"
  }
}
```

The `filter` argument values use SQL `where` clause syntax with any of the columns in the `aws_cloudtrail_log` table.

You can also collect logs for a single account:

```hcl
connection "aws" "aws_profile" {
  profile = "my-profile"
}

partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.aws_profile
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/CloudTrail/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"
  }
}
```

Or for a single region:

```hcl
connection "aws" "aws_profile" {
  profile = "my-profile"
}

partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.aws_profile
    bucket      = "cloudtrail-s3-log-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/CloudTrail/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json.gz"
  }
}
```

For more examples using the `aws_s3_bucket` source, please see [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket).

You can also work with local files, like the [public dataset from flaws.cloud](https://summitroute.com/blog/2020/10/09/public_dataset_of_cloudtrail_logs_from_flaws_cloud/):

```hcl
partition "aws_cloudtrail_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/mscott/cloudtrail_logs"]
    file_layout = "%{DATA}.json.gz"
  }
}
```

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
