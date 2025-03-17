---
title: "Tailpipe Table: aws_cost_and_usage_focus_1_0 - Query AWS Cost and Usage Reports (Focus 1.0)"
description: "AWS Cost and Usage Reports provide a detailed breakdown of cost, usage, and billing details for your AWS account."
---
# Table: aws_cost_and_usage_focus_1_0 - Query AWS Cost and Usage Reports (Focus 1.0)

The `aws_cost_and_usage_focus_1_0` table enables querying AWS Cost and Usage Report (CUR) data using the Focus 1.0 schema. This table provides granular insights into AWS billing, cost allocation, discounts, pricing, and resource-level usage.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_cost_and_usage_focus_1_0` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_and_usage_focus_1_0#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_and_usage_focus_1_0" "my_cur" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) data for all `aws_cost_and_usage_focus_1_0` partitions:

```sh
tailpipe collect aws_cost_and_usage_focus_1_0
```

Or for a single partition:

```sh
tailpipe collect aws_cost_and_usage_focus_1_0.my_cur
```

## Query

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_cost_and_usage_focus_1_0)**

### Monthly Cost Breakdown

Retrieve the total cost for each month, grouped by AWS account.

```sql
select
  date_trunc('month', charge_period_start) as billing_month,
  sub_account_id as account_id,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus_1_0
group by
  billing_month, account_id
order by
  billing_month desc;
```

### Top 10 Most Expensive AWS Services

List the top 10 AWS services with the highest costs.

```sql
select
  service_name,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus_1_0
group by
  service_name
order by
  total_cost desc
limit 10;
```

### High-Volume Resource Consumption

Identify resources with the highest usage quantity.

```sql
select
  resource_id,
  resource_name,
  sum(consumed_quantity) as total_usage,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus_1_0
group by
  resource_id, resource_name
order by
  total_usage desc
limit 10;
```

### Cost Breakdown by Region

Get a breakdown of cost and usage by AWS region.

```sql
select
  region_name,
  sum(billed_cost) as total_cost,
  sum(consumed_quantity) as total_usage
from
  aws_cost_and_usage_focus_1_0
group by
  region_name
order by
  total_cost desc;
```

## Example Configurations

### Collect Cost and Usage Reports from an S3 Bucket

Collect AWS CUR files stored in an S3 bucket.

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_and_usage_focus_1_0" "my_cur" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
  }
}
```

### Collect Reports from an S3 Bucket with a Prefix

Collect AWS CUR files stored in an S3 bucket using a prefix.

```hcl
partition "aws_cost_and_usage_focus_1_0" "my_cur_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect Reports from Local Files

You can also collect AWS CUR files from local files.

```hcl
partition "aws_cost_and_usage_focus_1_0" "local_cur" {
  source "file"  {
    paths       = ["/Users/myuser/aws_cur"]
    file_layout = "%{DATA}.csv.gz"
  }
}
```

### Filter Only Compute Costs

Use the filter argument in your partition to collect only compute-related costs.

```hcl
partition "aws_cost_and_usage_focus_1_0" "compute_costs" {
  filter = "service_category = 'Compute'"

  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
  }
}
```

### Collect Reports for All Accounts in an AWS Organization

For a specific AWS Organization, collect CUR data for all accounts.

```hcl
partition "aws_cost_and_usage_focus_1_0" "org_cur" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.billing_account
    bucket      = "aws-cur-org-bucket"
    file_layout = "%{DATA:prefix}/%{DATA:exportName}/%{DATA:data}/%{DATA:folderPath}/%{DATA:timestamp}/%{DATA}.csv.gz"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `%{DATA:prefix}/%{DATA:exportName}/%{DATA:data}/%{DATA:folderPath}/%{DATA:timestamp}/%{DATA}.csv.gz` |