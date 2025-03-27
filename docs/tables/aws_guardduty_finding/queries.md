## Activity Examples

### Daily Activity Trends

Analyze the daily distribution of GuardDuty findings to identify security patterns and potential attack campaigns over time. This temporal analysis helps establish baseline activity and detect anomalous spikes in security events.

```sql
select
  strftime(created_at, '%Y-%m-%d') as finding_date,
  count(*) as finding_count
from
  aws_guardduty_finding
group by
  finding_date
order by
  finding_date asc;
```

```yaml
folder: GuardDuty
```

### Top 10 Finding Types

Generate a ranked list of the most prevalent GuardDuty finding types, helping security teams focus on the most common security issues affecting their infrastructure. This insight drives prioritization of security controls and remediation efforts.

```sql
select
  type,
  count(*) as finding_count
from
  aws_guardduty_finding
group by
  type
order by
  finding_count desc
limit 10;
```

```yaml
folder: GuardDuty
```

### Findings by Account and Region

Perform a comprehensive analysis of security findings across your AWS organization, breaking down incidents by account and region while calculating average severity scores. This helps identify high-risk areas and potential security gaps in your multi-account infrastructure.

```sql
select
  account_id,
  region,
  count(*) as finding_count,
  round(avg(severity), 2) as avg_severity
from
  aws_guardduty_finding
group by
  account_id,
  region
order by
  finding_count desc;
```

```yaml
folder: GuardDuty
```

### Findings by Severity Level

Strategically categorize GuardDuty findings into High (7.0-8.9), Medium (4.0-6.9), and Low (0.1-3.9) severity bands to enable risk-based prioritization and resource allocation for incident response teams.

```sql
select
  case
    when severity >= 7 then 'High (7.0-8.9)'
    when severity >= 4 then 'Medium (4.0-6.9)'
    else 'Low (0.1-3.9)'
  end as severity_level,
  count(*) as finding_count
from
  aws_guardduty_finding
group by
  severity_level
order by
  finding_count desc;
```

```yaml
folder: GuardDuty
```

## Detection Examples

### Privilege Escalation Attempts

Monitor and detect potential privilege escalation activities that could indicate lateral movement attempts or compromised credentials. This query helps identify attackers trying to gain elevated access within your AWS environment.

```sql
select
  tp_timestamp,
  title,
  type,
  severity,
  description,
  account_id,
  region
from
  aws_guardduty_finding
where
  type like 'PrivilegeEscalation:%'
order by
  tp_timestamp desc;
```

```yaml
folder: GuardDuty
```

### Cryptocurrency Mining Activity

Detect potential cryptojacking attempts and unauthorized cryptocurrency mining operations that could indicate compromised resources. This query correlates various indicators including specific finding types and descriptive patterns to identify mining activities.

```sql
select
  tp_timestamp,
  title,
  type,
  severity,
  description,
  resource.resource_type as resource_type,
  account_id,
  region
from
  aws_guardduty_finding
where
  type like 'CryptoCurrency:%'
  or title like '%crypto%'
  or title like '%mining%'
order by
  severity desc,
  tp_timestamp desc;
```

```yaml
folder: GuardDuty
```

### Suspicious API Calls

Track unauthorized, discovery-oriented, and stealthy API calls that may indicate reconnaissance activities or attempted breaches. This comprehensive detection helps identify potential attackers probing your AWS environment.

```sql
select
  tp_timestamp,
  title,
  type,
  severity,
  description,
  account_id,
  region
from
  aws_guardduty_finding
where
  type like 'UnauthorizedAccess:IAMUser%'
  or type like 'Discovery:IAMUser%'
  or type like 'Stealth:IAMUser%'
order by
  tp_timestamp desc;
```

```yaml
folder: GuardDuty
```

### Malware Detection on EC2 Instances

Monitor EC2 instances for malware detections using GuardDuty's Malware Protection feature. This query helps identify compromised instances and potential malware infections in your compute environment.

```sql
select
  tp_timestamp,
  title,
  type,
  severity,
  service.feature_name as feature_name,
  resource.resource_type as resource_type,
  resource.resource_details as resource_details,
  account_id,
  region
from
  aws_guardduty_finding
where
  service.feature_name = 'MalwareProtection'
  and resource.resource_type = 'Instance'
order by
  tp_timestamp desc;
```

```yaml
folder: GuardDuty
```

### Abnormal Resource Behavior

Identify high-severity (>=5) behavioral anomalies, backdoors, and trojan activities that could indicate compromised AWS resources. This query focuses on critical threats requiring immediate investigation.

```sql
select
  tp_timestamp,
  title,
  type,
  severity,
  description,
  resource.resource_type as resource_type,
  account_id,
  region
from
  aws_guardduty_finding
where
  (type like 'Behavior:%' or type like 'Backdoor:%' or type like 'Trojan:%')
  and severity >= 5
order by
  severity desc,
  tp_timestamp desc;
```

```yaml
folder: GuardDuty
```

## Operational Examples

### Runtime Monitoring Findings

Investigate security events related to runtime behavior in containerized (EKS) and virtual machine (EC2) environments. This analysis helps detect suspicious processes, unauthorized activities, and potential security violations during program execution.

```sql
select
  tp_timestamp,
  title,
  type,
  severity,
  description,
  resource.resource_type as resource_type,
  account_id,
  region
from
  aws_guardduty_finding
where
  service.feature_name = 'RuntimeMonitoring'
order by
  tp_timestamp desc;
```

```yaml
folder: GuardDuty
```

### IAM Credential Misuse

Track potential compromises or misuse of IAM credentials by analyzing AccessKey-related findings. This query helps identify unauthorized access attempts and potential credential leaks across your AWS environment.

```sql
select
  tp_timestamp,
  title,
  type,
  severity,
  resource.access_key_details as access_key_details,
  account_id,
  region
from
  aws_guardduty_finding
where
  resource.resource_type = 'AccessKey'
order by
  severity desc,
  tp_timestamp desc;
```

```yaml
folder: GuardDuty
```

### Process Analysis for Runtime Threats

Perform detailed analysis of process-level activities detected by GuardDuty, including process names, paths, and command lines. This deep inspection helps security teams understand attack techniques and malicious process behaviors.

```sql
select
  tp_timestamp,
  title,
  severity,
  service.runtime_details.process.name as process_name,
  service.runtime_details.process.executable_path as executable_path,
  service.runtime_details.process.command_line_example as command_line,
  account_id,
  region
from
  aws_guardduty_finding
where
  service.runtime_details.process is not null
order by
  severity desc,
  tp_timestamp desc;
```

```yaml
folder: GuardDuty
```

## Volume Examples

### High Volume of Findings

Identify AWS accounts experiencing unusually high numbers of security findings (>50 per day), which may indicate targeted attacks, misconfiguration, or systematic security issues requiring immediate attention.

```sql
select
  account_id,
  count(*) as finding_count,
  date_trunc('day', created_at) as date
from
  aws_guardduty_finding
group by
  account_id,
  date
having
  count(*) > 50
order by
  finding_count desc;
```

```yaml
folder: GuardDuty
```

### Repeat Findings

Track recurring security issues by identifying findings that appear multiple times (>5 occurrences), helping teams focus on persistent security problems that require systematic resolution rather than one-off fixes.

```sql
select
  type,
  title,
  count(*) as occurrences
from
  aws_guardduty_finding
group by
  type,
  title
having
  count(*) > 5
order by
  occurrences desc;
```

```yaml
folder: GuardDuty
```

## Baseline Examples

### New Finding Types

Monitor the emergence of new security finding types within the last 7 days, helping security teams quickly identify and respond to novel threats or attack patterns in your AWS environment.

```sql
select
  type,
  min(created_at) as first_seen,
  count(*) as occurrence_count
from
  aws_guardduty_finding
group by
  type
having
  min(created_at) > current_date - interval '7 days'
order by
  first_seen desc;
```

```yaml
folder: GuardDuty
```
