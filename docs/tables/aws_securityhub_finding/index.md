---
title: "Tailpipe Table: aws_securityhub_finding - Query AWS Security Hub Findings"
description: "AWS Security Hub findings provide comprehensive security findings from various AWS security services and partner integrations, including details about potential security issues and compliance violations."
---

# Table: aws_securityhub_finding - Query AWS Security Hub Findings

The `aws_securityhub_finding` table allows you to query data from [AWS Security Hub findings](https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-findings.html). This table provides detailed information about potential security issues and compliance violations detected across your AWS accounts and resources, including severity levels, compliance status, affected resources, and recommended remediation steps.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_securityhub_finding` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_securityhub_finding#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "security_account" {
  profile = "my-security-account"
}

partition "aws_securityhub_finding" "my_findings" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-securityhub-findings-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) findings for all `aws_securityhub_finding` partitions:

```sh
tailpipe collect aws_securityhub_finding
```

Or for a single partition:

```sh
tailpipe collect aws_securityhub_finding.my_findings
```

## Query

**[Explore 8+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_securityhub_finding)**

### High Severity Findings

List all high severity security findings with detailed resource information.

```sql
select
  tp_timestamp,
  title,
  types,
  severity,
  description,
  region,
  resources,
  remediation.recommendation.text as remediation_text
from
  aws_securityhub_finding
where
  severity.normalized >= 70
order by
  severity.normalized desc,
  tp_timestamp desc;
```

### Findings by Type

Group findings by type with severity and temporal information.

```sql
select
  types,
  count(*) as finding_count,
  round(avg(severity.normalized), 2) as avg_severity
from
  aws_securityhub_finding
group by
  types
order by
  finding_count desc;
```

### Recent Findings with Resource Details

Examine recent security findings with comprehensive resource and remediation information.

```sql
select
  tp_timestamp,
  title,
  types,
  severity,
  resources,
  region,
  (workflow ->> 'status') as workflow_status,
  remediation.recommendation.text as remediation_text
from
  aws_securityhub_finding
where
  created_at > (current_date - interval '7 days')
order by
  tp_timestamp desc;
```

## Example Configurations

### Collect findings from an S3 bucket

Collect Security Hub findings stored in an S3 bucket using the default log file format.

```hcl
connection "aws" "security_account" {
  profile = "my-security-account"
}

partition "aws_securityhub_finding" "my_findings" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-securityhub-findings-bucket"
  }
}
```

### Collect findings from an S3 bucket with a prefix

Collect Security Hub findings stored in an S3 bucket using a prefix.

```hcl
partition "aws_securityhub_finding" "my_findings_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-securityhub-findings-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect findings from local files

You can also collect Security Hub findings from local files.

```hcl
partition "aws_securityhub_finding" "local_findings" {
  source "file"  {
    paths       = ["/Users/myuser/securityhub_findings"]
    file_layout = `%{DATA}.jsonl.gz`
  }
}
```

### Filter high severity findings only

Use the filter argument in your partition to focus on high severity findings, reducing the size of local storage.

```hcl
partition "aws_securityhub_finding" "high_severity_findings" {
  filter = "severity.normalized >= 70"

  source "aws_s3_bucket" {
    connection = connection.aws.security_account
    bucket     = "aws-securityhub-findings-bucket"
  }
}
```

### Collect findings for all accounts in an organization

For a specific organization, collect findings for all accounts and regions.

```hcl
partition "aws_securityhub_finding" "my_findings_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.security_account
    bucket      = "securityhub-findings-bucket"
    file_layout = `AWSLogs/o-aa111bb222/%{NUMBER:account_id}/%{DATA:security_hub_integrrated_product_name}/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json`
  }
}
```

### Collect findings for a single account

For a specific account, collect findings for all regions.

```hcl
partition "aws_securityhub_finding" "my_findings_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.security_account
    bucket      = "securityhub-findings-bucket"
    file_layout = `AWSLogs/123456789012/%{DATA:security_hub_integrrated_product_name}/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json`
  }
}
```

### Collect findings for a single region

For all accounts, collect findings from us-east-1.

```hcl
partition "aws_securityhub_finding" "my_findings_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.security_account
    bucket      = "securityhub-findings-bucket"
    file_layout = `AWSLogs/%{NUMBER:account_id}/%{DATA:security_hub_integrrated_product_name}/us-east-1/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json`
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/%{NUMBER:account_id}/%{DATA:security_hub_integrrated_product_name}/%{DATA:region}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{DATA}.json` | 
