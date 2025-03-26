---
title: "Tailpipe Table: aws_detailed_billing_report - Query AWS Detailed Billing Reports"
description: "AWS Detailed Billing Report provides comprehensive information on AWS costs, usage, and billing details for your AWS accounts."
---

# Table: aws_detailed_billing_report - Query AWS Detailed Billing Reports

The `aws_detailed_billing_report` table allows you to query AWS Detailed Billing Report data from AWS. This table provides insights into your AWS billing, usage, resource details, and cost allocations, offering a historical view of your AWS spending before the introduction of the Cost and Usage Reports.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_detailed_billing_report` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_detailed_billing_report#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_detailed_billing_report" "my_dbr" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-dbr-billing-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) data for all `aws_detailed_billing_report` partitions:

```sh
tailpipe collect aws_detailed_billing_report
```

Or for a single partition:

```sh
tailpipe collect aws_detailed_billing_report.my_dbr
```

## Query

**[Explore 12+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_detailed_billing_report)**

### Monthly Cost Breakdown

Retrieve the total cost for each month, grouped by AWS account.

```sql
select
  date_trunc('month', usage_start_date) as billing_month,
  linked_account_id as account_id,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
group by
  billing_month, account_id
order by
  billing_month desc;
```

### Top 10 Most Expensive AWS Services

List the top 10 AWS services with the highest costs.

```sql
select
  product_name,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
group by
  product_name
order by
  total_cost desc
limit 10;
```

### Cost by Operation

Analyze costs by operation type.

```sql
select
  operation,
  sum(total_cost) as total_cost,
  sum(usage_quantity) as total_usage
from
  aws_detailed_billing_report
where
  operation is not null
group by
  operation
order by
  total_cost desc
limit 10;
```

## Example Configurations

### Collect Detailed Billing Reports from an S3 Bucket

Collect AWS DBR files stored in an S3 bucket.

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_detailed_billing_report" "my_dbr" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-dbr-billing-bucket"
  }
}
```

### Collect Reports from an S3 Bucket with a Prefix

Collect AWS DBR files stored in an S3 bucket using a prefix.

```hcl
partition "aws_detailed_billing_report" "my_dbr_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-dbr-billing-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect Reports from Local Files

You can also collect AWS DBR files from local files.

```hcl
partition "aws_detailed_billing_report" "local_dbr" {
  source "file"  {
    paths       = ["/Users/myuser/aws_dbr"]
    file_layout = "%{DATA}.csv"
  }
}
```

### Filter Only Compute Costs

Use the filter argument in your partition to collect only compute-related costs.

```hcl
partition "aws_detailed_billing_report" "compute_costs" {
  filter = "product_name = 'Amazon Elastic Compute Cloud'"

  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-dbr-billing-bucket"
  }
}
```

### Collect Reports for All Accounts in an AWS Organization

For a specific AWS Organization, collect DBR data for all accounts.

```hcl
partition "aws_detailed_billing_report" "org_dbr" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.billing_account
    bucket      = "aws-dbr-org-bucket"
    file_layout = "%{DATA}.csv"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default       |
| ----------- | ------------- |
| file_layout | `%{DATA}.csv` |
