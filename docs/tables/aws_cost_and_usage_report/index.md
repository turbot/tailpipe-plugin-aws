---
title: "Tailpipe Table: aws_cost_and_usage_report - Query AWS Cost and Usage Reports"
description: "AWS Cost and Usage Reports contain the most comprehensive set of cost and usage data available for your AWS account."
---
# Table: aws_cost_and_usage_report - Query AWS Cost and Usage Reports

The `aws_cost_and_usage_report` table allows you to query [AWS Cost and Usage Report (CUR)](https://docs.aws.amazon.com/cur/latest/userguide/table-dictionary-cur2.html) data. This table provides insights into your AWS billing, usage, cost categories, and discounts.

Limitations and notes:
- This table currently supports collecting from `.gzip` files only.
- [CUR 2.0](https://docs.aws.amazon.com/cur/latest/userguide/table-dictionary-cur2.html) and [Legacy CUR](https://docs.aws.amazon.com/cur/latest/userguide/what-is-cur.html) data can be collected by this table.
- When determining each log's timestamp, the table uses the following order of precedence:
  - `line_item_usage_start_date`
  - `line_item_usage_end_date`
  - `billing_period_start`
  - `billing_period_end`
  - If none of the columns above are present, logs will not be collected, and Tailpipe will report these as errors.
- When determining each log's index, the table uses the following order of precedence:
  - `line_item_usage_account_id` (CUR 2.0)/`lineItem/UsageAccountId` (Legacy CUR)
  - `line_item_resource_id` (CUR 2.0)/`lineItem/ResourceId` (Legacy CUR)
    - This column is used if the resource ID is an ARN and is only included in reports if **Include resource IDs** was selected during export creation.
  - If none of the columns above are present, the log will use `default` as the index instead of an AWS account ID.

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

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_cost_and_usage_report)**

### Monthly Cost Breakdown

Retrieve the total cost for each month, grouped by AWS account.

```sql
select
  date_trunc('month', bill_billing_period_start_date) as billing_month,
  line_item_usage_account_id as account_id,
  sum(line_item_unblended_cost) as total_cost
from
  aws_cost_and_usage_report
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

### High-Volume Data Transfer Usage

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

### Usage by Region

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

### Collect for a specific CUR 2.0 export

For a specific export (`my-cur-2-0-export` in this example), collect Cost and Usage 2.0 reports.

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_and_usage_report" "specific_cur_2_0" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.billing_account
    bucket      = "aws-cur-billing-bucket"
    prefix      = "my/prefix/"
    file_layout = "my-cur-2-0-export/data/%{DATA:partition}/?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.gz"
  }
}
```

### Collect for a specific Legacy CUR export

For a specific export (`my-cur-legacy-export` in this example), collect Legacy Cost and Usage reports.

```hcl
partition "aws_cost_and_usage_report" "specific_cur_legacy" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.billing_account
    bucket      = "aws-cur-billing-bucket"
    prefix      = "my/prefix/"
    file_layout = "my-cur-legacy-export/%{INT:from_date}-%{INT:to_date}/(?:%{DATA:assembly_id}/)?%{DATA:file_name}.csv.zip"
  }
}
```

### Collect reports from an S3 bucket

Collect Cost and Usage reports stored in an S3 bucket that use the [default log file name format](https://docs.aws.amazon.com/cur/latest/userguide/dataexports-export-delivery.html#export-summary).

**Note**: We only recommend using the default log file name format if the bucket and prefix combination contains Cost and Usage reports. If other reports, like the Cost and Usage FOCUS report, are stored in the same S3 bucket with the same prefix, Tailpipe will attempt to collect from these too, resulting in errors.

```hcl
partition "aws_cost_and_usage_report" "my_cur" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect reports from local files

You can also reports from local files.

```hcl
partition "aws_cost_and_usage_report" "local_cur" {
  source "file"  {
    paths       = ["/Users/myuser/aws_cur"]
    file_layout = "%{DATA}.csv.gz"
  }
}
```

### Collect only compute service costs

Use the filter argument in your partition to collect only compute product family costs.

```hcl
partition "aws_cost_and_usage_report" "compute_costs" {
  filter = "product_compute_family is not null"

  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-billing-bucket"
    prefix     = "my/prefix/"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `%{DATA:export_name}/(?:data/%{DATA:partition}/)?(?:%{INT:from_date}-%{INT:to_date}/)?(?:%{DATA:assembly_id}/)?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.(?:zip\|gz)` |
