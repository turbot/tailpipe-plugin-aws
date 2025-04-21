## Activity Examples

### Daily Activity Trends

Analyze the daily distribution of Security Hub findings to identify security patterns and potential security issues over time.

```sql
select
  strftime(tp_timestamp, '%Y-%m-%d') as finding_date,
  count(*) as finding_count,
  round(avg(severity.normalized), 2) as avg_severity
from
  aws_securityhub_finding
group by
  finding_date
order by
  finding_date asc;
```

```yaml
folder: SecurityHub
```

### Recent Findings Analysis

Analyze recent security findings with detailed resource and severity information.

```sql
select
  tp_timestamp,
  title,
  types,
  severity,
  resources,
  tp_index as account_id,
  region,
  workflow_state,
  remediation.recommendation.text as remediation_text
from
  aws_securityhub_finding
where
  tp_timestamp > current_date - interval '7 days'
order by
  severity.normalized desc,
  tp_timestamp desc;
```

```yaml
folder: SecurityHub
```

### Top 10 Finding Types

Generate a ranked list of the most prevalent Security Hub finding types with severity information.

```sql
select
  types,
  count(*) as finding_count,
  round(avg(severity.normalized), 2) as avg_severity
from
  aws_securityhub_finding
group by
  types
order by
  finding_count desc
limit 10;
```

```yaml
folder: SecurityHub
```

<!-- https://docs.aws.amazon.com/securityhub/1.0/APIReference/API_Severity.html -->

### Findings by Account and Region

Analyze security findings across your AWS organization with detailed severity information.

```sql
select
  tp_index as account_id,
  region,
  count(*) as finding_count,
  round(avg(severity.normalized), 2) as avg_severity,
  sum(case when severity.normalized >= 90 then 1 else 0 end) as critical_severity_count,
  sum(case when severity.normalized >= 70 and severity.normalized < 90 then 1 else 0 end) as high_severity_count,
  sum(case when severity.normalized >= 40 and severity.normalized < 70 then 1 else 0 end) as medium_severity_count,
  sum(case when severity.normalized >= 1 and severity.normalized < 40 then 1 else 0 end) as low_severity_count,
  sum(case when severity.normalized = 0 then 1 else 0 end) as informational_severity_count
from
  aws_securityhub_finding
group by
  account_id,
  region
order by
  critical_severity_count desc;
```

```yaml
folder: SecurityHub
```

### Findings by Severity Level

Categorize Security Hub findings into severity bands with detailed counts and percentages.

```sql
select
  case
    when severity.normalized >= 90 then 'Critical (90-100)'
    when severity.normalized >= 70 then 'High (70-89)'
    when severity.normalized >= 40 then 'Medium (40-69)'
    when severity.normalized >= 1 then 'Low (1-39)'
    else 'Informational (0)'
  end as severity_level,
  count(*) as finding_count,
  round(count(*) * 100.0 / sum(count(*)) over(), 2) as percentage
from
  aws_securityhub_finding
group by
  severity_level
order by
  case severity_level
    when 'Critical (90-100)' then 1
    when 'High (70-89)' then 2
    when 'Medium (40-69)' then 3
    when 'Low (1-39)' then 4
    else 5
  end;
```

```yaml
folder: SecurityHub
```

## Compliance Examples

### Compliance Status Overview

Monitor compliance status with detailed severity information.

```sql
select
  compliance.status,
  compliance.security_control_id,
  count(*) as finding_count,
  round(avg(severity.normalized), 2) as avg_severity
from
  aws_securityhub_finding
where
  compliance is not null
group by
  compliance.status,
  compliance.security_control_id
order by
  finding_count desc;
```

```yaml
folder: SecurityHub
```

## Detection Examples

### Detect High Severity Findings with Remediation

  ```sql
select
  tp_timestamp,
  title,
  types,
  severity,
  description,
  tp_index as account_id,
  region,
  resources,
  remediation.recommendation.text as remediation_text
from
  aws_securityhub_finding
where
  severity.normalized >= 70
order by
  severity.normalized desc,
  tp_timestamp desc;
```

```yaml
folder: SecurityHub
```

### Lambda Function Security Issues

Identify security issues in Lambda functions, focusing on public access.

```sql
select
  tp_timestamp,
  title,
  severity.normalized as severity,
  json_extract(resources, '$[0].id') as function_arn,
  json_extract(resources, '$[0].details.awslambdafunction.runtime') as runtime,
  workflow_state
from
  aws_securityhub_finding
where
  json_extract(resources, '$[0].type') = '"AwsLambdaFunction"'
  and title ilike '%public access%'
  and severity.normalized >= 70
order by
  severity desc,
  tp_timestamp desc;
```

```yaml
folder: SecurityHub
```
