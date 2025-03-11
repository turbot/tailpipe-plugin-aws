## Activity Examples

### Daily Cost Trends

Track daily AWS spending trends.

```sql
select
  date_trunc('day', line_item_usage_start_date) as usage_date,
  sum(cast(line_item_unblended_cost as double)) as daily_cost
from
  aws_cost_and_usage_report
where
  line_item_usage_start_date >= current_date - interval '30 days'
group by
  usage_date
order by
  usage_date;
```

```yaml
folder: Account
```

### Top 10 Costly Services

Identify the most expensive AWS services over the past month.

```sql
select
  product_service_code as service,
  sum(cast(line_item_unblended_cost as double)) as total_cost
from
  aws_cost_and_usage_report
where
  line_item_usage_start_date >= current_date - interval '1 month'
group by
  service
order by
  total_cost desc
limit 10;
```

```yaml
folder: Account
```

### Top Spending Accounts

Determine which AWS accounts are generating the highest costs.

```sql
select
  line_item_usage_account_id as account_id,
  sum(cast(line_item_unblended_cost as double)) as total_cost
from
  aws_cost_and_usage_report
group by
  account_id
order by
  total_cost desc
limit 10;
```

```yaml
folder: Account
```

## Detection Examples

### Unusual Cost Spikes

Detect services with a sudden increase in costs.

```sql
with monthly_costs as (
  select
    date_trunc('month', line_item_usage_start_date) as month,
    line_item_product_code as service,
    sum(cast(line_item_unblended_cost as double)) as total_cost
  from
    aws_cost_and_usage_report
  where
    line_item_usage_start_date >= current_date - interval '2 months'
  group by
    month, service
)
select
  current_month.service,
  previous_month.total_cost as previous_cost,
  current_month.total_cost as current_cost,
  round(((current_month.total_cost - previous_month.total_cost) / previous_month.total_cost) * 100, 2) as percentage_increase
from
  monthly_costs current_month
join
  monthly_costs previous_month
  on current_month.service = previous_month.service
  and current_month.month = date_trunc('month', current_date)
  and previous_month.month = date_trunc('month', current_date) - interval '1 month'
where
  previous_month.total_cost > 0
  and ((current_month.total_cost - previous_month.total_cost) / previous_month.total_cost) > 0.2
order by
  percentage_increase desc;
```

```yaml
folder: Account
```

### High Data Transfer Usage

Find accounts with high outbound data transfer usage.

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

```yaml
folder: Account
```

## Operational Examples

### EC2 Cost Breakdown by Instance Type

Analyze EC2 costs based on instance types.

```sql
select
  product_instance_type as instance_type,
  count(distinct line_item_resource_id) as instance_count,
  sum(cast(line_item_unblended_cost as double)) as total_cost
from
  aws_cost_and_usage_report
where
  line_item_usage_start_date >= current_date - interval '30 days'
  and line_item_product_code = 'AmazonEC2'
  and line_item_usage_type like '%BoxUsage%'
group by
  instance_type
order by
  total_cost desc;
```

```yaml
folder: EC2
```

### EBS Volumes with High Costs

Identify expensive Amazon EBS volumes.

```sql
select
  line_item_resource_id,
  line_item_product_code,
  product_volume_api_name,
  product_product_name,
  sum(cast(line_item_unblended_cost as double)) as total_cost
from
  aws_cost_and_usage_report
where
  line_item_usage_start_date >= current_date - interval '30 days'
  and line_item_product_code = 'AmazonEC2'
  and line_item_usage_type like '%EBS:VolumeUsage%'
group by
  line_item_resource_id, 
  product_volume_api_name, 
  line_item_product_code, 
  product_product_name
order by
  total_cost desc
limit 10;
```

```yaml
folder: EBS
```

## Volume Examples

### Reserved Instance Utilization

Calculate the percentage of EC2 usage covered by Reserved Instances.

```sql
with ri_eligible as (
  select
    sum(cast(line_item_unblended_cost as double)) as total_on_demand_cost
  from
    aws_cost_and_usage_report
  where
    line_item_usage_start_date >= current_date - interval '30 days'
    and line_item_product_code = 'AmazonEC2'
    and line_item_line_item_type = 'Usage'
    and pricing_term = 'OnDemand'
),
ri_covered as (
  select
    sum(cast(line_item_unblended_cost as double)) as total_ri_cost
  from
    aws_cost_and_usage_report
  where
    line_item_usage_start_date >= current_date - interval '30 days'
    and line_item_product_code = 'AmazonEC2'
    and (pricing_term = 'Reserved' or line_item_line_item_type = 'DiscountedUsage')
)
select
  e.total_on_demand_cost,
  c.total_ri_cost,
  round((c.total_ri_cost / (e.total_on_demand_cost + c.total_ri_cost)) * 100, 2) as ri_coverage_percent
from
  ri_eligible e,
  ri_covered c;
```

```yaml
folder: EC2
```

### High-Volume API Calls

Detect AWS services generating a high volume of API calls.

```sql
select
  product_service_code as service,
  count(*) as api_call_count
from
  aws_cost_and_usage_report
where
  line_item_usage_start_date >= current_date - interval '30 days'
group by
  service
order by
  api_call_count desc
limit 10;
```

```yaml
folder: Account
```

## Baseline Examples

### Cost Breakdown by AWS Region

Compare spending across AWS regions.

```sql
select
  product_region as region,
  sum(cast(line_item_unblended_cost as double)) as total_cost
from
  aws_cost_and_usage_report
group by
  region
order by
  total_cost desc;
```

```yaml
folder: Account
```

### Services with Unexpected Costs

Identify services that usually have low costs but show unexpected spending increases.

```sql
with avg_service_cost as (
  select
    product_service_code as service,
    avg(cast(line_item_unblended_cost as double)) as avg_monthly_cost
  from
    aws_cost_and_usage_report
  where
    line_item_usage_start_date >= current_date - interval '6 months'
  group by
    service
)
select
  c.service,
  c.total_cost,
  a.avg_monthly_cost,
  round(((c.total_cost - a.avg_monthly_cost) / a.avg_monthly_cost) * 100, 2) as percentage_increase
from (
  select
    product_service_code as service,
    sum(cast(line_item_unblended_cost as double)) as total_cost
  from
    aws_cost_and_usage_report
  where
    line_item_usage_start_date >= current_date - interval '1 month'
  group by
    service
) c
join avg_service_cost a on c.service = a.service
where
  ((c.total_cost - a.avg_monthly_cost) / a.avg_monthly_cost) > 0.2
order by
  percentage_increase desc;
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
