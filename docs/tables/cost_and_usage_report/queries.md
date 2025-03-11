# Queries: aws_cost_and_usage_report

## Activity Examples

### Total cost by service

Find the total cost by AWS service for the current billing period.

```sql
select
  line_item_product_code as service,
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

```yaml
folder: Account
```

### Cost by account

Break down costs by AWS account to understand spending across your organization.

```sql
select
  line_item_usage_account_id as account_id,
  line_item_usage_account_name as account_name,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  line_item_usage_account_id,
  line_item_usage_account_name,
  line_item_currency_code
order by
  total_cost desc;
```

```yaml
folder: Account
```

### Cost by operation

Analyze costs by operation to understand which specific API calls or actions are driving your costs.

```sql
select
  line_item_operation as operation,
  line_item_product_code as service,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  line_item_operation is not null
  and bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  line_item_operation,
  line_item_product_code,
  line_item_currency_code
order by
  total_cost desc
limit 20;
```

```yaml
folder: Account
```

## Detection Examples

### Top 10 most expensive resources

Identify the most expensive resources in your AWS environment.

```sql
select
  line_item_resource_id as resource_id,
  line_item_product_code as service,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  line_item_resource_id is not null
  and bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  line_item_resource_id,
  line_item_product_code,
  line_item_currency_code
order by
  total_cost desc
limit 10;
```

```yaml
folder: Account
```

### Sudden cost increases

Detect sudden increases in daily costs compared to the previous day.

```sql
with daily_costs as (
  select
    date_trunc('day', line_item_usage_start_date) as usage_date,
    sum(line_item_unblended_cost) as daily_cost,
    line_item_currency_code as currency
  from
    aws_cost_and_usage_report
  where
    line_item_usage_start_date >= current_date - interval '30 days'
  group by
    usage_date,
    line_item_currency_code
)
select
  current_day.usage_date,
  current_day.daily_cost,
  previous_day.daily_cost as previous_day_cost,
  (current_day.daily_cost - previous_day.daily_cost) as cost_increase,
  ((current_day.daily_cost - previous_day.daily_cost) / previous_day.daily_cost * 100) as percentage_increase,
  current_day.currency
from
  daily_costs as current_day
  join daily_costs as previous_day on (
    previous_day.usage_date = current_day.usage_date - interval '1 day'
    and previous_day.currency = current_day.currency
  )
where
  current_day.daily_cost > previous_day.daily_cost * 1.2  -- 20% increase threshold
order by
  percentage_increase desc;
```

```yaml
folder: Account
```

### Underutilized reserved instances

Identify reserved instances that are not being fully utilized.

```sql
select
  reservation_reservation_arn as reservation_arn,
  reservation_unused_quantity as unused_quantity,
  reservation_total_reserved_units as total_reserved_units,
  (reservation_unused_quantity / nullif(cast(reservation_total_reserved_units as double), 0)) * 100 as unused_percentage,
  reservation_unused_recurring_fee as unused_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  reservation_reservation_arn is not null
  and reservation_unused_quantity > 0
  and bill_billing_period_start_date >= date_trunc('month', current_date)
order by
  unused_cost desc;
```

```yaml
folder: Account
```

## Operational Examples

### EC2 instance costs by instance type

Analyze EC2 costs by instance type to optimize your compute resources.

```sql
select
  product_instance_type as instance_type,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  line_item_product_code = 'AmazonEC2'
  and product_instance_type is not null
  and bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  product_instance_type,
  line_item_currency_code
order by
  total_cost desc;
```

```yaml
folder: Account
```

### Cost by region

Analyze costs by AWS region to understand geographical distribution of your spending.

```sql
select
  product_region_code as region,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  product_region_code is not null
  and bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  product_region_code,
  line_item_currency_code
order by
  total_cost desc;
```

```yaml
folder: Account
```

### Cost by tag

Analyze costs by resource tags to understand spending by project, environment, or other dimensions.

```sql
select
  json_extract_path_text(resource_tags, 'user:Project') as project_tag,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  resource_tags is not null
  and json_extract_path_text(resource_tags, 'user:Project') is not null
  and bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  project_tag,
  line_item_currency_code
order by
  total_cost desc;
```

```yaml
folder: Account
```

## Volume Examples

### Daily cost trend

Analyze daily cost trends over the last 30 days.

```sql
select
  date_trunc('day', line_item_usage_start_date) as usage_date,
  sum(line_item_unblended_cost) as daily_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  line_item_usage_start_date >= current_date - interval '30 days'
group by
  usage_date,
  line_item_currency_code
order by
  usage_date;
```

```yaml
folder: Account
```

### Monthly cost trend

Analyze monthly cost trends over the past year.

```sql
select
  date_trunc('month', line_item_usage_start_date) as month,
  sum(line_item_unblended_cost) as monthly_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  line_item_usage_start_date >= current_date - interval '12 months'
group by
  month,
  line_item_currency_code
order by
  month;
```

```yaml
folder: Account
```

### Cost by usage type

Break down costs by usage type to understand what specific resources or operations are driving your costs.

```sql
select
  line_item_usage_type as usage_type,
  line_item_product_code as service,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  line_item_usage_type,
  line_item_product_code,
  line_item_currency_code
order by
  total_cost desc
limit 20;
```

```yaml
folder: Account
```

### Cost by line item type

Analyze costs by line item type to understand the breakdown of charges, credits, and refunds.

```sql
select
  line_item_line_item_type as line_item_type,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  line_item_line_item_type,
  line_item_currency_code
order by
  total_cost desc;
```

```yaml
folder: Account
```

## Baseline Examples

### Savings from reserved instances

Calculate savings from reserved instances compared to on-demand pricing.

```sql
select
  reservation_reservation_arn as reservation_arn,
  sum(pricing_public_on_demand_cost) as on_demand_cost,
  sum(line_item_unblended_cost) as actual_cost,
  sum(pricing_public_on_demand_cost - line_item_unblended_cost) as savings,
  (sum(pricing_public_on_demand_cost - line_item_unblended_cost) / sum(pricing_public_on_demand_cost)) * 100 as savings_percentage,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  reservation_reservation_arn is not null
  and pricing_public_on_demand_cost > 0
  and bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  reservation_reservation_arn,
  line_item_currency_code
order by
  savings desc;
```

```yaml
folder: Account
```

### Savings plan utilization

Analyze savings plan utilization to optimize your commitment-based discounts.

```sql
select
  savings_plan_savings_plan_arn as savings_plan_arn,
  sum(savings_plan_savings_plan_effective_cost) as covered_cost,
  sum(pricing_public_on_demand_cost) as on_demand_cost,
  (sum(savings_plan_savings_plan_effective_cost) / sum(pricing_public_on_demand_cost)) * 100 as utilization_percentage,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  savings_plan_savings_plan_arn is not null
  and pricing_public_on_demand_cost > 0
  and bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  savings_plan_savings_plan_arn,
  line_item_currency_code
order by
  covered_cost desc;
```

```yaml
folder: Account
```

### Cost by purchase option

Compare costs across different purchase options (on-demand, reserved, spot, etc.).

```sql
select
  pricing_purchase_option as purchase_option,
  sum(line_item_unblended_cost) as total_cost,
  line_item_currency_code as currency
from
  aws_cost_and_usage_report
where
  pricing_purchase_option is not null
  and bill_billing_period_start_date >= date_trunc('month', current_date)
group by
  pricing_purchase_option,
  line_item_currency_code
order by
  total_cost desc;
```

```yaml
folder: Account
```

### Cost comparison across billing periods

Compare costs between the current and previous billing period.

```sql
with current_period as (
  select
    line_item_product_code as service,
    sum(line_item_unblended_cost) as cost,
    line_item_currency_code as currency
  from
    aws_cost_and_usage_report
  where
    bill_billing_period_start_date = (
      select max(bill_billing_period_start_date)
      from aws_cost_and_usage_report
    )
  group by
    line_item_product_code,
    line_item_currency_code
),
previous_period as (
  select
    line_item_product_code as service,
    sum(line_item_unblended_cost) as cost,
    line_item_currency_code as currency
  from
    aws_cost_and_usage_report
  where
    bill_billing_period_start_date = (
      select max(bill_billing_period_start_date)
      from aws_cost_and_usage_report
      where bill_billing_period_start_date < (
        select max(bill_billing_period_start_date)
        from aws_cost_and_usage_report
      )
    )
  group by
    line_item_product_code,
    line_item_currency_code
)

select
  current_period.service,
  current_period.cost as current_cost,
  previous_period.cost as previous_cost,
  (current_period.cost - coalesce(previous_period.cost, 0)) as cost_difference,
  case
    when previous_period.cost > 0 then
      ((current_period.cost - previous_period.cost) / previous_period.cost) * 100
    else
      null
  end as percentage_change,
  current_period.currency
from
  current_period
  left join previous_period on (
    current_period.service = previous_period.service
    and current_period.currency = previous_period.currency
  )
order by
  current_cost desc;
```

```yaml
folder: Account
```
