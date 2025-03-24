## Activity Examples

### Daily Cost Trends By Account

Track daily AWS cost trends from your detailed billing reports, segregated by account.

```sql
select
  cast(usage_start_date as date) as usage_date,
  linked_account_id as account_id,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
group by
  usage_date, account_id
order by
  usage_date, account_id;
```

```yaml
folder: Detailed Billing Report
```

### Top Services By Account

Identify the most expensive AWS services for each account.

```sql
select
  linked_account_id as account_id,
  product_name,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
where
  product_name is not null
group by
  account_id, product_name
order by
  account_id, total_cost desc;
```

```yaml
folder: Detailed Billing Report
```

### Top Account Spending

Determine which AWS accounts have the highest spending.

```sql
select
  coalesce(linked_account_id, payer_account_id) as account_id,
  coalesce(linked_account_name, payer_account_name) as account_name,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
group by
  account_id,
  account_name
order by
  total_cost desc
limit 10;
```

```yaml
folder: Detailed Billing Report
```

## Detection Examples

### Cost Variations By Account

Compare costs between billing periods by service and account.

```sql
with monthly_costs as (
  select
    linked_account_id as account_id,
    date_trunc('month', usage_start_date) as month,
    product_name,
    sum(total_cost) as total_cost
  from
    aws_detailed_billing_report
  where
    product_name is not null
  group by
    account_id, month, product_name
)
select
  account_id,
  product_name,
  sum(case when month = date_trunc('month', current_date - interval '1 month') then total_cost else 0 end) as "last_month_cost",
  sum(case when month = date_trunc('month', current_date) then total_cost else 0 end) as "current_month_cost"
from
  monthly_costs
group by
  account_id, product_name
having
  sum(case when month = date_trunc('month', current_date - interval '1 month') then total_cost else 0 end) > 0
  or sum(case when month = date_trunc('month', current_date) then total_cost else 0 end) > 0
order by
  account_id, "last_month_cost" desc;
```

```yaml
folder: Detailed Billing Report
```

### Tax Related Items By Account

Identify items that have associated tax amounts for each account.

```sql
select
  coalesce(linked_account_id, payer_account_id) as account_id,
  coalesce(linked_account_name, payer_account_name) as account_name,
  product_name,
  tax_type,
  sum(tax_amount) as tax_total,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
where
  tax_type is not null
  and trim(tax_type) != ''
  and tax_amount is not null
group by
  account_id, account_name, product_name, tax_type
order by
  account_id, tax_total desc;

```

```yaml
folder: Detailed Billing Report
```

## Operational Examples

### EC2 Costs By Operation And Account

Analyze EC2 costs by different operations for each account.

```sql
select
  coalesce(linked_account_id, payer_account_id) as account_id,
  operation,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
where
  product_name = 'Amazon Elastic Compute Cloud'
  and operation is not null
group by
  account_id, operation
order by
  account_id, total_cost desc;
```

```yaml
folder: EC2
```

### S3 Storage and Data Transfer Costs By Account

Break down S3 costs between operations for each account.

```sql
select
  linked_account_id as account_id,
  operation,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
where
  product_name = 'Amazon Simple Storage Service'
group by
  account_id, operation
order by
  account_id, total_cost desc;
```

```yaml
folder: S3
```

## Volume Examples

### Data Transfer Costs By Account

Find cost patterns for data transfer operations by account.

```sql
select
  coalesce(linked_account_id, payer_account_id) as account_id,
  product_name,
  operation,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
where
  operation like '%DataTransfer%'
  or operation like '%Transfer%'
group by
  account_id, product_name, operation
order by
  account_id, total_cost desc;
```

```yaml
folder: Detailed Billing Report
```

### Product Operation Cost Analysis By Account

Analyze which operations contribute the most to costs for each account.

```sql
select
  linked_account_id as account_id,
  product_name,
  operation,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
where
  total_cost > 0
  and operation is not null
  and operation != ''
group by
  account_id, product_name, operation
order by
  account_id, total_cost desc;
```

```yaml
folder: Detailed Billing Report
```

## Baseline Examples

### Cost Distribution By Service And Account

Compare costs across AWS services for each account.

```sql
select
  linked_account_id as account_id,
  product_name,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
where
  product_name is not null
group by
  account_id, product_name
order by
  account_id, total_cost desc;
```

```yaml
folder: Detailed Billing Report
```

### Account Type Cost Distribution

Compare cost metrics between different types of AWS accounts.

```sql
select
  case
    when linked_account_id = payer_account_id then 'Management Account'
    when linked_account_id != '' then 'Member Account'
    else 'Consolidated Billing'
  end as account_type,
  count(distinct linked_account_id) as account_count,
  sum(total_cost) as total_cost
from
  aws_detailed_billing_report
group by
  account_type
order by
  total_cost desc;
```

```yaml
folder: Detailed Billing Report
```