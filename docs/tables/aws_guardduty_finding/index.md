---
title: "Tailpipe Table: aws_guardduty_finding - Query AWS GuardDuty Findings"
description: "AWS GuardDuty findings provide alerts and intelligence about potential security threats and suspicious activities detected in your AWS environment."
---

# Table: aws_guardduty_finding - Query AWS GuardDuty Findings

The `aws_guardduty_finding` table allows you to query data from AWS GuardDuty findings. This table provides detailed information about potential security threats detected within your AWS environment, including the affected resources, severity levels, finding types, and contextual details to help you investigate and respond to security issues.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_guardduty_finding` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_guardduty_finding#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "security_account" {
  profile = "my-security-account"
}

partition "aws_guardduty_finding" "my_findings" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-guardduty-findings-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) findings for all `aws_guardduty_finding` partitions:

```sh
tailpipe collect aws_guardduty_finding
```

Or for a single partition:

```sh
tailpipe collect aws_guardduty_finding.my_findings
```

## Query

**[Explore example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_guardduty_finding)**

### High Severity Findings

List all high severity security findings to prioritize your security response.

```sql
select
  created_at,
  account_id,
  region,
  title,
  type,
  severity,
  description
from
  aws_guardduty_finding
where
  severity >= 7.0
order by
  severity desc,
  created_at desc;
```

### Findings by Type

Group findings by type to understand the most common security threats in your environment.

```sql
select
  type,
  count(*) as finding_count,
  min(created_at) as first_seen,
  max(created_at) as last_seen,
  round(avg(severity), 2) as avg_severity
from
  aws_guardduty_finding
group by
  type
order by
  finding_count desc;
```

### Recent Findings with Resource Details

Examine recent security findings along with details of the affected resources.

```sql
select
  created_at,
  title,
  type,
  severity,
  account_id,
  region,
  resource ->> 'resource_type' as resource_type,
  resource ->> 'resource_details' as resource_details
from
  aws_guardduty_finding
where
  created_at > current_date - interval '7 days'
order by
  created_at desc;
```

## Example Configurations

### Collect findings from an S3 bucket

Collect GuardDuty findings stored in an S3 bucket using the default log file format.

```hcl
connection "aws" "security_account" {
  profile = "my-security-account"
}

partition "aws_guardduty_finding" "my_findings" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-guardduty-findings-bucket"
  }
}
```

### Collect findings from an S3 bucket with a prefix

Collect GuardDuty findings stored in an S3 bucket using a prefix.

```hcl
partition "aws_guardduty_finding" "my_findings_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-guardduty-findings-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect findings from local files

You can also collect GuardDuty findings from local files.

```hcl
partition "aws_guardduty_finding" "local_findings" {
  source "file"  {
    paths       = ["/Users/myuser/guardduty_findings"]
    file_layout = "%{DATA}.jsonl.gz"
  }
}
```

### Filter high severity findings only

Use the filter argument in your partition to focus on high severity findings, reducing the size of local storage.

```hcl
partition "aws_guardduty_finding" "high_severity_findings" {
  filter = "severity >= 7.0"

  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-guardduty-findings-bucket"
  }
}
```

### Collect findings for all accounts in an organization

For a specific organization, collect findings for all accounts and regions.

```hcl
partition "aws_guardduty_finding" "my_findings_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.security_account
    bucket      = "guardduty-findings-bucket"
    file_layout = "AWSLogs/o-aa111bb222/%{NUMBER:account_id}/GuardDuty/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.jsonl.gz"
  }
}
```

### Collect findings for a single account

For a specific account, collect findings for all regions.

```hcl
partition "aws_guardduty_finding" "my_findings_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.security_account
    bucket      = "guardduty-findings-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/GuardDuty/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.jsonl.gz"
  }
}
```

### Collect findings for a single region

For all accounts, collect findings from us-east-1.

```hcl
partition "aws_guardduty_finding" "my_findings_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.security_account
    bucket      = "guardduty-findings-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/GuardDuty/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.jsonl.gz"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/GuardDuty/%{DATA:region_path}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.jsonl.gz` |
``` 