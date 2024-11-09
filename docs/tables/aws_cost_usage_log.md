# AWS Cost and Usage Log

## 1. Total Monthly Cost by Account
Calculates the total monthly cost across linked AWS accounts to identify which accounts are incurring the most cost.

```sql
select
  linked_account_id,
  linked_account_name,
  sum(total_cost) as monthly_total_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  linked_account_id, linked_account_name
order by
  monthly_total_cost desc;
```

## 2. Top Costly AWS Services
Displays the top AWS services by total cost, helping to identify where the majority of spending occurs.

```sql
select
  product_name,
  sum(total_cost) as total_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  product_name
order by
  total_cost desc
limit 10;
```

## 3. Daily Cost Trend for the Current Month
Tracks the daily expenditure trend to identify unusual cost spikes.

```sql
select
  tp_date,
  sum(total_cost) as daily_total_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  tp_date
order by
  tp_date;
```

## 4. Cost Breakdown by Usage Type
Breaks down costs by usage type, providing insights into specific usage patterns that may drive costs.

```sql
select
  usage_type,
  sum(total_cost) as total_usage_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  usage_type
order by
  total_usage_cost desc;
```

## 5. High-Cost Operations
Lists operations that have incurred the highest costs, which can indicate potential cost optimization areas.

```sql
select
  operation,
  sum(total_cost) as total_operation_cost
from
  aws_cost_and_usage_log
where
  operation is not null
  and tp_date >= date_trunc('month', current_date)
group by
  operation
order by
  total_operation_cost desc
limit 10;
```

## 6. Average Daily Cost for the Current Month
Calculates the average daily cost for the current month, providing a baseline for typical daily expenditure.

```sql
select
  avg(total_cost) as average_daily_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date);
```

## 7. Top Accounts by Usage Quantity
Identifies the accounts with the highest usage quantity, helping to associate high usage with specific accounts.

```sql
select
  linked_account_id,
  linked_account_name,
  sum(usage_quantity) as total_usage_quantity
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  linked_account_id, linked_account_name
order by
  total_usage_quantity desc
limit 10;
```

## 8. High Tax Amount by Product
Displays products with the highest tax amount, useful for identifying tax-heavy services.

```sql
select
  product_name,
  sum(tax_amount) as total_tax_amount
from
  aws_cost_and_usage_log
where
  tax_amount > 0
  and tp_date >= date_trunc('month', current_date)
group by
  product_name
order by
  total_tax_amount desc;
```

## 9. Total Credits by Account
Shows the total credits applied to each account, helping to understand any discounts or rebates.

```sql
select
  linked_account_id,
  linked_account_name,
  sum(credits) as total_credits
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  linked_account_id, linked_account_name
order by
  total_credits desc;
```

## 10. Cost Allocation by Product Code
Provides a cost breakdown by product code, aiding in identifying high-cost products.

```sql
select
  product_code,
  sum(total_cost) as total_product_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  product_code
order by
  total_product_cost desc;
```

## 11. Highest Cost by Resource Item Description
Identifies the resources or items with the highest cost based on item description to locate costly services or operations.

```sql
select
  item_description,
  sum(total_cost) as total_item_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  item_description
order by
  total_item_cost desc
limit 10;
```

## 12. Cost Trends by Payer Account
Shows the monthly cost trend for each payer account, useful for tracking cost changes across accounts.

```sql
select
  payer_account_id,
  payer_account_name,
  date_trunc('month', tp_date) as month,
  sum(total_cost) as monthly_cost
from
  aws_cost_and_usage_log
group by
  payer_account_id, payer_account_name, month
order by
  payer_account_id, month;
```

## 13. Monthly Tax Amount by Linked Account
Displays the monthly tax amount by linked account, helping to understand tax-related expenses.

```sql
select
  linked_account_id,
  linked_account_name,
  sum(tax_amount) as monthly_tax_amount
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  linked_account_id, linked_account_name
order by
  monthly_tax_amount desc;
```

## 14. Top Billing Periods by Cost
Ranks billing periods by total cost to highlight the periods with the highest spending.

```sql
select
  billing_period_start_date,
  billing_period_end_date,
  sum(total_cost) as billing_period_total_cost
from
  aws_cost_and_usage_log
group by
  billing_period_start_date, billing_period_end_date
order by
  billing_period_total_cost desc
limit 10;
```

## 15. Cost by Currency Code
Groups costs by currency code, which is useful in multi-currency scenarios to monitor costs in different currencies.

```sql
select
  currency_code,
  sum(total_cost) as total_currency_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  currency_code
order by
  total_currency_cost desc;
```

## 16. Top 10 Product Codes by Usage Quantity
Identifies product codes with the highest usage quantities, offering insights into frequently used services.

```sql
select
  product_code,
  sum(usage_quantity) as total_usage_quantity
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  product_code
order by
  total_usage_quantity desc
limit 10;
```

## 17. Cost by Seller of Record
Breaks down costs by seller to understand charges from specific sellers of AWS products.

```sql
select
  seller_of_record,
  sum(total_cost) as total_seller_cost
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  seller_of_record
order by
  total_seller_cost desc;
```

## 18. High Blended Rate Products
Lists products with a high blended rate, helping to identify products with significant effective rates.

```sql
select
  product_name,
  avg(blended_rate) as average_blended_rate
from
  aws_cost_and_usage_log
where
  blended_rate is not null
  and tp_date >= date_trunc('month', current_date)
group by
  product_name
order by
  average_blended_rate desc
limit 10;
```

## 19. Invoices with Highest Cost Before Tax
Displays invoices with the highest cost before tax, useful for analyzing specific large charges.

```sql
select
  invoice_id,
  sum(cost_before_tax) as total_cost_before_tax
from
  aws_cost_and_usage_log
where
  tp_date >= date_trunc('month', current_date)
group by
  invoice_id
order by
  total_cost_before_tax desc
limit 10;
```

## 20. Monthly Cost Change Percentage by Product
Calculates the month-over-month cost change percentage for each product, providing insights into spending increases or decreases.

```sql
with monthly_costs as (
  select
    product_name,
    date_trunc('month', tp_date) as month,
    sum(total_cost) as monthly_cost
  from
    aws_cost_and_usage_log
  group by
    product_name, month
),
monthly_changes as (
  select
    product_name,
    month,
    monthly_cost,
    lag(monthly_cost) over (partition by product_name order by month) as previous_month_cost
  from
    monthly_costs
)
select
  product_name,
  month,
  (monthly_cost - previous_month_cost) / previous_month_cost * 100 as cost_change_percentage
from
  monthly_changes
where
  previous_month_cost is not null
order by
  product_name, month;
```