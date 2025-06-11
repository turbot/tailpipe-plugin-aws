## Activity Examples

### Top 10 Cost-Saving Recommendations

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

```yaml
folder: Cost Optimization Hub
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

```yaml
folder: Cost Optimization Hub
```

### Account-Level Optimization Opportunities

Identify which accounts have the most optimization potential.

```sql
select
  account_id,
  count(*) as recommendation_count,
  round(sum(estimated_monthly_savings_after_discount), 2) as total_potential_savings,
  round(avg(estimated_savings_percentage_after_discount), 2) as avg_savings_percentage
from
  aws_cost_optimization_recommendation
group by
  account_id
order by
  total_potential_savings desc;
```

```yaml
folder: Cost Optimization Hub
```

## Detection Examples

### Recommendations by Source

Analyze which recommendation sources provide the most opportunities.

```sql
select
  recommendation_source,
  count(*) as recommendation_count,
  round(sum(estimated_monthly_savings_after_discount), 2) as total_potential_savings,
  round(avg(estimated_monthly_savings_after_discount), 2) as avg_savings_per_recommendation
from
  aws_cost_optimization_recommendation
group by
  recommendation_source
order by
  total_potential_savings desc;
```

```yaml
folder: Cost Optimization Hub
```

### Stale Recommendations

Identify recommendations that haven't been refreshed recently.

```sql
select
  tp_timestamp,
  recommendation_id,
  account_id,
  current_resource_type,
  last_refresh_timestamp,
  estimated_monthly_savings_after_discount as savings_amount
from
  aws_cost_optimization_recommendation
where
  last_refresh_timestamp::timestamp < now()::timestamp - interval '7' day
order by
  tp_timestamp,
  last_refresh_timestamp desc;
```

```yaml
folder: Cost Optimization Hub
```

## Operational Examples

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

```yaml
folder: Cost Optimization Hub
```

### Easy to Implement Recommendations

Find recommendations that are easy to implement and don't require restarts.

```sql
select
  tp_timestamp,
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

```yaml
folder: Cost Optimization Hub
```

## Volume Examples

### Top 10 High Value Recommendations by Resource Type

Identify which resource types offer the highest saving potential.

```sql
select
  current_resource_type,
  count(*) as recommendation_count,
  round(sum(estimated_monthly_savings_after_discount), 2) as total_savings,
  round(avg(estimated_monthly_savings_after_discount), 2) as avg_savings_per_recommendation
from
  aws_cost_optimization_recommendation
where
  estimated_monthly_savings_after_discount > 100
group by
  current_resource_type
order by
  total_savings desc
limit 10;
```

```yaml
folder: Cost Optimization Hub
```

## Baseline Examples

### Projected Annual Savings

Calculate potential annual savings from implementing all recommendations.

```sql
select
  round(sum(estimated_monthly_savings_after_discount) * 12, 2) as annual_savings_potential,
  count(*) as total_recommendations,
  round(avg(estimated_savings_percentage_after_discount), 2) as avg_savings_percentage
from
  aws_cost_optimization_recommendation;
```

```yaml
folder: Cost Optimization Hub
```
