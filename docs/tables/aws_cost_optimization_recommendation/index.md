---
title: "Tailpipe Table: aws_cost_optimization_recommendation - Query AWS Cost Optimization Recommendations"
description: "AWS Cost Optimization Recommendations provide insights into potential cost-saving opportunities across your AWS resources."
---

# Table: aws_cost_optimization_recommendation - Query AWS Cost Optimization Recommendations

The `aws_cost_optimization_recommendation` table allows you to query [AWS cost optimization recommendations](https://docs.aws.amazon.com/cur/latest/userguide/table-dictionary-cor.html) for your AWS resources. These recommendations identify potential savings opportunities based on usage patterns, resource configurations, and AWS pricing options.

Limitations and notes:
- This table currently supports collecting from `.gzip` files only.
- If the export does not include the `last_refresh_timestamp` column, logs will **not** be collected since this is the only timestamp related column in the report.
- If the export does not include the `account_ID` column, logs will be collected, but all rows will be indexed under `default` instead of an AWS account ID.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_cost_optimization_recommendation` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_optimization_recommendation#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_optimization_recommendation" "my_recommendations" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cost-optimization-recommendations-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) data for all `aws_cost_optimization_recommendation` partitions:

```sh
tailpipe collect aws_cost_optimization_recommendation
```

Or for a single partition:

```sh
tailpipe collect aws_cost_optimization_recommendation.my_recommendations
```

## Query

**[Explore 10+ example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_cost_optimization_recommendation)**

### Top Cost-Saving Recommendations

Identify the recommendations with the highest potential monthly savings.

```sql
select
  recommendation_id,
  account_id,
  current_resource_type,
  recommended_resource_type,
  region,
  estimated_monthly_savings_after_discount as savings_amount,
  estimated_savings_percentage_after_discount as savings_percentage,
  currency_code
from
  aws_cost_optimization_recommendation
order by
  savings_amount desc
limit 10;
```

### Low Effort Recommendations

Find recommendations that are easy to implement and don't require restarts.

```sql
select
  recommendation_id,
  account_id,
  current_resource_type,
  recommended_resource_type,
  implementation_effort,
  restart_needed,
  estimated_monthly_savings_after_discount as savings_amount
from
  aws_cost_optimization_recommendation
where
  implementation_effort = 'Low'
  and restart_needed = false
order by
  savings_amount desc;
```

### Recommendations by Service Type

Group recommendations by resource type to focus optimization efforts.

```sql
select
  current_resource_type,
  count(*) as recommendation_count,
  round(sum(estimated_monthly_savings_after_discount), 2) as total_potential_savings,
  round(avg(estimated_savings_percentage_after_discount), 2) as avg_savings_percentage
from
  aws_cost_optimization_recommendation
group by
  current_resource_type
order by
  total_potential_savings desc;
```

### Regional Optimization Opportunities

Analyze cost-saving opportunities by AWS region.

```sql
select
  region,
  count(*) as recommendation_count,
  round(sum(estimated_monthly_savings_after_discount), 2) as total_potential_savings,
  round(avg(estimated_savings_percentage_after_discount), 2) as avg_savings_percentage
from
  aws_cost_optimization_recommendation
group by
  region
order by
  total_potential_savings desc;
```

## Example Configurations

### Collect for a specific export

For a specific export (`my-recommendations-export` in this example), collect cost optimization recommendations.

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_optimization_recommendation" "specific_recommendations" {
  source "aws_s3_bucket" {
    connection  = connection.aws.billing_account
    bucket      = "aws-cost-optimization-recommendations-bucket"
    prefix      = "my/prefix/"
    file_layout = "my-recommendations-export/data/%{DATA:partition}/(?:%{TIMESTAMP_ISO8601:timestamp}-%{UUID:execution_id}/)?%{DATA:filename}.csv.gz"
  }
}
```

### Collect recommendations from an S3 bucket

Collect cost optimization recommendations stored in an S3 bucket that use the [default log file name format](https://docs.aws.amazon.com/cur/latest/userguide/dataexports-export-delivery.html#export-summary).

**Note**: We only recommend using the default log file name format if the bucket and prefix combination contains cost optimization recommendations. If other reports, like the Cost and Usage FOCUS report, are stored in the same S3 bucket with the same prefix, Tailpipe will attempt to collect from these too, resulting in errors.

```hcl
connection "aws" "billing_account" {
  profile = "my-billing-account"
}

partition "aws_cost_optimization_recommendation" "my_recommendations" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cost-optimization-recommendations-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect recommendations from local files

You can also collect AWS cost optimization recommendations from local files.

```hcl
partition "aws_cost_optimization_recommendation" "local_recommendations" {
  source "file"  {
    paths       = ["/Users/myuser/aws_recommendations"]
    file_layout = "%{DATA}.csv.gz"
  }
}
```

### Collect only high-value recommendations

Use the filter argument in your partition to collect only high-value recommendations.

```hcl
partition "aws_cost_optimization_recommendation" "high_value_recommendations" {
  filter = "estimated_monthly_savings_after_discount > 100"

  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cost-optimization-recommendations-bucket"
    prefix     = "my/prefix/"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `%{DATA:export_name}/data/%{DATA:partition}/(?:%{TIMESTAMP_ISO8601:timestamp}-%{UUID:execution_id}/)?%{DATA:filename}.csv.gz`|
