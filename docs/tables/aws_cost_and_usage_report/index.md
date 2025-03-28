---
title: "Tailpipe Table: aws_cost_and_usage_report - Query AWS Cost and Usage Reports"
description: "AWS Cost and Usage Reports contain the most comprehensive set of cost and usage data available for your AWS account."
---
# Table: aws_cost_and_usage_report - Query AWS Cost and Usage Reports

The `aws_cost_and_usage_report` table allows you to query AWS Cost and Usage Report (CUR) data from AWS. This table provides insights into your AWS billing, usage, cost categories, and discounts.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_cost_and_usage_report` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_and_usage_report#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_and_usage_report" "my_cur" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) data for all `aws_cost_and_usage_report` partitions:

```sh
tailpipe collect aws_cost_and_usage_report
```

Or for a single partition:

```sh
tailpipe collect aws_cost_and_usage_report.my_cur
```

## Query

**[Explore 12+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_cost_and_usage_report)**

### monthly cost breakdown

Retrieve the total cost for each month, grouped by AWS account.

```sql
select
  date_trunc('month', bill_billing_period_start_date) as billing_month,
  line_item_usage_account_id as account_id,
  sum(line_item_unblended_cost) as total_cost
from
  aws_cost_and_usage_report
group by
  billing_month, account_id
order by
  billing_month desc;
```

### Top 10 most expensive services

List the top 10 AWS services with the highest costs.

```sql
select
  product_service_code,
  sum(line_item_unblended_cost) as total_cost
from
  aws_cost_and_usage_report
group by
  product_service_code
order by
  total_cost desc
limit 10;
```

### High-Volume data transfer usage

Identify accounts with high outbound data transfer usage.

```sql
select
  line_item_usage_account_id as account_id,
  sum(cast(line_item_usage_amount as double)) as total_data_transfer_gb
from
  aws_cost_and_usage_report
where
  line_item_usage_type like 'DataTransfer-Out%'
group by
  account_id
order by
  total_data_transfer_gb desc
limit 10;
```

### Usage by region

Get a breakdown of usage and cost by AWS region.

```sql
select
  product_region_code,
  sum(cast(line_item_usage_amount as double)) as total_usage,
  sum(line_item_unblended_cost) as total_cost
from
  aws_cost_and_usage_report
group by
  product_region_code
order by
  total_cost desc;
```

## Example Configurations

### Collect cost and usage reports from an S3 bucket

Collect AWS CUR files stored in an S3 bucket.

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_and_usage_report" "my_cur" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
  }
}
```

### Collect reports from an S3 bucket with a prefix

Collect AWS CUR files stored in an S3 bucket using a prefix.

```hcl
partition "aws_cost_and_usage_report" "my_cur_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect reports from local files

You can also collect AWS CUR files from local files.

```hcl
partition "aws_cost_and_usage_report" "local_cur" {
  source "file"  {
    paths       = ["/Users/myuser/aws_cur"]
    file_layout = "%{DATA}.csv.gz"
  }
}
```

### Filter only compute costs

Use the filter argument in your partition to collect only compute-related costs.

```hcl
partition "aws_cost_and_usage_report" "compute_costs" {
  filter = "product_compute_family is not null"

  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
  }
}
```

### Collect reports for all accounts in an AWS organization

For a specific AWS Organization, collect CUR data for all accounts.

```hcl
partition "aws_cost_and_usage_report" "org_cur" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.billing_account
    bucket      = "aws-cur-org-bucket"
    prefix      = "reports"
    file_layout = "%{DATA:export_name}/(?:data/%{DATA:partition}/)?(?:%{INT:from_date}-%{INT:to_date}/)?(?:%{DATA:assembly_id}/)?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.(?:zip\|gz)"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `%{DATA:export_name}/(?:data/%{DATA:partition}/)?(?:%{INT:from_date}-%{INT:to_date}/)?(?:%{DATA:assembly_id}/)?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.(?:zip\|gz)` |
