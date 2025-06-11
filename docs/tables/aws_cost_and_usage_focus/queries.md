## Activity Examples

### Daily Cost Trends

Track AWS spending trends on a daily basis.

```sql
select
  date_trunc('day', charge_period_start) as usage_date,
  sum(billed_cost) as daily_cost
from
  aws_cost_and_usage_focus
where
  charge_period_start >= current_date - interval '30' day
group by
  usage_date
order by
  usage_date;
```

```yaml
folder: Cost and Usage Report
```

### Top 10 Costly Services

Identify the highest-cost AWS services over the last month.

```sql
select
  service_name,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus
where
  billing_period_start >= current_date - interval '1' month
group by
  service_name
order by
  total_cost desc
limit 10;
```

```yaml
folder: Cost and Usage Report
```

### Month-over-Month Cost Growth

Track how your costs grow over time by analyzing month-over-month changes by service.

```sql
with monthly_costs as (
  select
    date_trunc('month', billing_period_start) as month,
    service_name,
    sum(billed_cost) as monthly_cost
  from
    aws_cost_and_usage_focus
  where
    billing_period_start >= date_trunc('month', current_date) - interval '6' month
  group by
    month, service_name
)
select
  strftime(month, '%Y-%m') as "Month",
  service_name as "Service",
  monthly_cost as "Monthly Cost",
  lag(monthly_cost) over (partition by service_name order by month) as "Previous Month Cost",
  case
    when lag(monthly_cost) over (partition by service_name order by month) > 0 then
      round(
        (monthly_cost - lag(monthly_cost) over (partition by service_name order by month)) /
        lag(monthly_cost) over (partition by service_name order by month) * 100,
        2
      )
    else null
  end as "Growth (%)"
from
  monthly_costs
where
  monthly_cost > 10
order by
  month desc,
  monthly_cost desc;
```

```yaml
folder: Cost and Usage Report
```

### Top 10 Spending Accounts

Find AWS accounts that have the highest spending.

```sql
select
  sub_account_id as account_id,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus
group by
  account_id
order by
  total_cost desc
limit 10;
```

```yaml
folder: Cost and Usage Report
```

## Detection Examples

### Top 10 Services With Daily Cost Variations

Detect services with significant daily cost variations over the past week.

```sql
with daily_costs as (
  select
    date_trunc('day', charge_period_start) as day,
    service_name,
    sum(billed_cost) as total_cost
  from
    aws_cost_and_usage_focus
  where
    charge_period_start >= current_date - interval '7' day
  group by
    day, service_name
)
select
  service_name,
  min(total_cost) as min_daily_cost,
  max(total_cost) as max_daily_cost,
  round(avg(total_cost), 2) as avg_daily_cost,
  round((max(total_cost) - min(total_cost)) / nullif(min(total_cost), 0) * 100, 2) as cost_variation_pct
from
  daily_costs
group by
  service_name
having
  min(total_cost) > 0
  and count(*) > 1
order by
  cost_variation_pct desc
limit 10;
```

```yaml
folder: Cost and Usage Report
```

## Operational Examples

### EC2 Cost Breakdown by Resource Type

Analyze EC2 costs by resource types.

```sql
select
  resource_type,
  count(distinct resource_id) as instance_count,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus
where
  billing_period_start >= current_date - interval '30' day
  and service_name = 'Amazon Elastic Compute Cloud'
group by
  resource_type
order by
  total_cost desc;
```

```yaml
folder: EC2
```

### Top 10 EBS Volumes by Cost

Identify expensive Amazon EBS volumes.

```sql
select
  resource_id,
  service_name,
  charge_description,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus
where
  billing_period_start >= current_date - interval '30' day
  and resource_type = 'volume'
group by
  resource_id,
  service_name,
  charge_description
order by
  total_cost desc
limit 10;
```

```yaml
folder: EBS
```

## Volume Examples

### Top 10 Accounts by Data Transfer Usage

Find accounts with the highest outbound data transfer usage.

```sql
select
  sub_account_id as account_id,
  sum(consumed_quantity) as total_data_transfer_gb
from
  aws_cost_and_usage_focus
where
  x_usage_type like 'DataTransfer-Out%'
group by
  account_id
order by
  total_data_transfer_gb desc
limit 10;
```

```yaml
folder: Cost and Usage Report
```

### Top 10 Costs by API

Identify the most expensive API operations by cost.

```sql
select
  service_name,
  x_operation as api_operation,
  sum(billed_cost) as total_cost,
  count(*) as api_call_count
from
  aws_cost_and_usage_focus
where
  billing_period_start >= current_date - interval '30' day
  and x_operation is not null
group by
  service_name, api_operation
order by
  total_cost desc
limit 10;
```

```yaml
folder: Cost and Usage Report
```

## Baseline Examples

### Cost Breakdown by AWS Region

Compare spending across AWS regions.

```sql
select
  region_name as region,
  sum(billed_cost) as total_cost
from
  aws_cost_and_usage_focus
group by
  region
order by
  total_cost desc;
```

```yaml
folder: Cost and Usage Report
```

### Cost Comparison Across Billing Periods

Compare costs between the current and previous billing periods.

```sql
select
  service_name,
  sum(case when billing_period_start = (select max(billing_period_start) from aws_cost_and_usage_focus) then billed_cost else 0 end) as current_period_cost,
  sum(case when billing_period_start = (select max(billing_period_start) - interval '1' month from aws_cost_and_usage_focus) then billed_cost else 0 end) as previous_period_cost
from
  aws_cost_and_usage_focus
group by
  service_name
order by
  current_period_cost desc;
```

```yaml
folder: Cost and Usage Report
```

### Cost Breakdown by Service Category

Analyze AWS costs by service category to identify spending patterns.

```sql
select
  service_category,
  sum(billed_cost) as "Total Cost",
  count(distinct service_name) as "Number of Services"
from
  aws_cost_and_usage_focus
where
  billing_period_start >= current_date - interval '1' month
group by
  service_category
order by
  "Total Cost" desc;
```

```yaml
folder: Cost and Usage Report
```
