---
title: "Tailpipe Table: aws_cost_and_usage_focus - Query AWS Cost and Usage Reports (FOCUS 1.0)"
description: "AWS Cost and Usage FOCUS report contains your cost and usage data formatted with FinOps Open Cost and Usage Specification (FOCUS) 1.0."
---

# Table: aws_cost_and_usage_focus - Query AWS Cost and Usage Reports (FOCUS 1.0)

The `aws_cost_and_usage_focus` table enables querying AWS Cost and Usage Report (CUR) data using the FOCUS 1.0 schema. This table provides granular insights into AWS billing, cost allocation, discounts, pricing, and resource-level usage.

Limitations and notes:
- This table currently supports collecting from `.gzip` files only.
- When determining each log's timestamp, the table uses the following order of precedence:
  - `ChargePeriodStart`
  - `ChargePeriodEnd`
  - `BillingPeriodStart`
  - `BillingPeriodEnd`
- If none of the columns above are present, then Tailpipe will be unable to collect logs from that export.
- If the export does not include the `SubAccountId` column, logs will still collected, but all rows will be indexed under `default` instead of an AWS account ID.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_cost_and_usage_focus` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_and_usage_focus#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_and_usage_focus" "my_cur_focus" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) data for all `aws_cost_and_usage_focus` partitions:

```sh
tailpipe collect aws_cost_and_usage_focus
```

Or for a single partition:

```sh
tailpipe collect aws_cost_and_usage_focus.my_cur_focus
```

## Query

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_cost_and_usage_focus)**

### Monthly Cost Breakdown

Retrieve the total cost for each month, grouped by AWS account.

```sql
select
  date_trunc('month', charge_period_start) as billing_month,
  sub_account_id as account_id,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus
group by
  billing_month,
  account_id
order by
  billing_month desc;
```

### Top 10 Most Expensive Services

List the top 10 AWS services with the highest costs.

```sql
select
  service_name,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus
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
  aws_cost_and_usage_focus
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
  aws_cost_and_usage_focus
group by
  region_name
order by
  total_cost desc;
```

## Example Configurations

### Collect for a specific export

For a specific export (`my-focus-export` in this example), collect Cost and Usage FOCUS reports.

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_and_usage_focus" "specific_cur_focus" {
  source "aws_s3_bucket"  {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-org-bucket"
    prefix     = "my/prefix/my-focus-export/"
  }
}
```

### Collect reports from an S3 bucket

Collect Cost and Usage FOCUS reports stored in an S3 bucket that use the [default log file name format](https://docs.aws.amazon.com/cur/latest/userguide/dataexports-export-delivery.html).

**Note**: We only recommend using the default log file name format if the bucket and prefix combination contains Cost and Usage FOCUS reports. If other reports, like the Cost and Usage Report 2.0, are stored in the same S3 bucket with the same prefix, Tailpipe will attempt to collect from these too, resulting in errors.

```hcl
partition "aws_cost_and_usage_focus" "my_cur_focus" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect reports from local files

You can also collect reports from local files.

```hcl
partition "aws_cost_and_usage_focus" "local_cur_focus" {
  source "file"  {
    paths       = ["/Users/myuser/aws_cur"]
    file_layout = "%{DATA}.csv.gz"
  }
}
```

### Collect only compute service costs

Use the filter argument in your partition to collect only compute service category costs.

```hcl
partition "aws_cost_and_usage_focus" "compute_costs" {
  filter = "service_category = 'Compute'"

  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    prefix     = "my/prefix/"
    bucket     = "aws-cur-billing-bucket"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default     |
| ----------- | ----------- |
| file_layout | `%{DATA:export_name}/data/%{DATA:partition}/(?:%{TIMESTAMP_ISO8601:timestamp}-%{UUID:execution_id}/)?%{DATA:filename}.csv.gz` |
