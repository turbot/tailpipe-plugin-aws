---
title: "Tailpipe Table: aws_cost_optimization_recommendation - Query AWS Cost Optimization Recommendations"
description: "AWS Cost Optimization Recommendations provide insights into potential cost-saving opportunities across your AWS resources."
---

# Table: aws_cost_optimization_recommendation - Query AWS Cost Optimization Recommendations

The `aws_cost_optimization_recommendation` table allows you to query cost optimization recommendations for your AWS resources. These recommendations identify potential savings opportunities based on usage patterns, resource configurations, and AWS pricing options.

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

### Easy to Implement Recommendations

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

### Collect Recommendations from an S3 Bucket

Collect AWS cost optimization recommendations stored in an S3 bucket.

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

### Collect Recommendations from an S3 Bucket with a Prefix

Collect AWS cost optimization recommendations stored in an S3 bucket using a prefix.

```hcl
partition "aws_cost_optimization_recommendation" "my_recommendations_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cost-optimization-recommendations-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect Recommendations from Local Files

You can also collect AWS cost optimization recommendations from local files.

```hcl
partition "aws_cost_optimization_recommendation" "local_recommendations" {
  source "file"  {
    paths       = ["/Users/myuser/aws_recommendations"]
    file_layout = "%{DATA}.csv.gz"
  }
}
```

### Filter Only High-Value Recommendations

Use the filter argument in your partition to collect only high-value recommendations.

```hcl
partition "aws_cost_optimization_recommendation" "high_value_recommendations" {
  filter = "estimated_monthly_savings_after_discount > 100"

  source "aws_s3_bucket" {
    connection = connection.aws.billing_account
    bucket     = "aws-cost-optimization-recommendations-bucket"
  }
}
```

### Collect Recommendations for All Accounts in an AWS Organization

For a specific AWS Organization, collect recommendation data for all accounts.

```hcl
partition "aws_cost_optimization_recommendation" "org_recommendations" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.billing_account
    bucket      = "aws-cost-optimization-recommendations-bucket"
    prefix      = "reports"
    file_layout = "%{DATA:export_name}/(?:data/%{DATA:partition}/)?(?:%{INT:from_date}-%{INT:to_date}/)?(?:%{DATA:assembly_id}/)?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.(?:zip|gz)"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `%{DATA:export_name}/(?:data/%{DATA:partition}/)?(?:%{INT:from_date}-%{INT:to_date}/)?(?:%{DATA:assembly_id}/)?(?:%{DATA:timestamp}-%{DATA:execution_id}/)?%{DATA:file_name}.csv.(?:zip\|gz)` |
