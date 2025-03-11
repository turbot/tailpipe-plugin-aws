---
title: "Tailpipe Table: aws_cost_and_usage_report - Query AWS Cost and Usage Reports"
description: "AWS Cost and Usage Reports (CUR) provide detailed information about your AWS costs and usage, including itemized usage data and cost allocation."
---

# Table: aws_cost_and_usage_report - Query AWS Cost and Usage Reports

The `aws_cost_and_usage_report` table allows you to query data from AWS Cost and Usage Reports. This table provides comprehensive information about your AWS spending, including detailed breakdowns by service, usage type, and resource tags.

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
    bucket     = "aws-cur-bucket"
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

### Total cost by service for the current billing period

```sql
select
  line_item_product_code,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  line_item_product_code,
  line_item_currency_code
order by
  total_cost desc;
```

### Resources with highest costs

```sql
select
  line_item_resource_id,
  line_item_product_code,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  line_item_resource_id is not null
group by
  line_item_resource_id,
  line_item_product_code,
  line_item_currency_code
order by
  total_cost desc
limit 10;
```

## Example Configurations

### Collect CUR from an S3 bucket

Collect Cost and Usage Reports stored in an S3 bucket using the default file format.

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_and_usage_report" "my_cur" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-bucket"
  }
}
```

### Collect CUR with a specific prefix

Collect Cost and Usage Reports stored in an S3 bucket using a specific prefix.

```hcl
partition "aws_cost_and_usage_report" "my_cur_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-bucket"
    prefix     = "my-report-prefix/"
  }
}
```

### Collect logs from local files

You can also collect CloudTrail logs from local files, like the [flaws.cloud public dataset](https://summitroute.com/blog/2020/10/09/public_dataset_of_cloudtrail_logs_from_flaws_cloud/).

```hcl
partition "aws_cloudtrail_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/cur"]
    file_layout = "%{DATA}.csv.gz"
  }
}
```

### Collect CUR for a specific service

Collect Cost and Usage Reports for a specific service using service-based filtering.

```hcl
partition "aws_cost_and_usage_report" "my_cur_filtered" {
  filter = "line_item_product_code like '%EC2%'"

  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cur-bucket"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument    | Default                                                                   |
| ----------- | ------------------------------------------------------------------------- |
| file_layout | `%{YEAR:year}/%{MONTHNUM:month}/%{DATA:report_name}/%{DATA}.csv.(zip|gz)` |
