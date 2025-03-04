# Query Reviews

## Daily Access Trends ✅

<details><summary>Query</summary>
### Daily Access Trends

Count access log entries per day to identify trends over time.

```sql
select
  strftime(timestamp, '%Y-%m-%d') as access_date,
  count(*) AS requests
from
  aws_s3_server_access_log
group by
  access_date
order by
  access_date asc;
```
</details>

<details><summary>SQL syntax checks ✅</summary>

| Criteria      | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Use 2 space indentation | ✅ |  |
| Query should end with a semicolon | ✅ |  |
| Keywords should be in lowercase | ✅ |  |
| Each clause is on its own line | ✅ |  |
| All columns exist in the schema | ✅ |  |
| STRUCT type columns use dot notation | ✅ |  |
| JSON type columns use `->` and `->>` operators | ✅ |  |
| JSON type columns are wrapped in parenthesis | ✅ |  |
| SQL query syntax uses valid DuckDB syntax | ✅ |  |

</details>

<details><summary>Query title and description checks ✅</summary>

| Criteria | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Title uses title case | ✅ |  |
| Title accurately describes the query | ✅ |  |
| Description explains what the query does | ✅ |  |
| Description explains why a user would run the query | ✅ |  |
| Description is concise | ✅ |  |

</details>

# Query Reviews

## Top 10 Accessed Objects ❌

<details><summary>Query</summary>
### Top 10 Accessed Objects

List the 10 most frequently accessed IAM objects.

```sql
select
bucket,
key,
  count(*) as requests
from
  aws_s3_server_access_log
where
  key is not null
group by
  bucket,
  key
order by
  requests desc
limit 30;
```
</details>

<details><summary>SQL syntax checks ❌</summary>

| Criteria      | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Use 2 space indentation | ❌ | Inconsistent indentation. Use 2 spaces for all lines. |
| Query should end with a semicolon | ✅ |  |
| Keywords should be in lowercase | ✅ |  |
| Each clause is on its own line | ✅ |  |
| All columns exist in the schema | ✅ |  |
| STRUCT type columns use dot notation | ✅ |  |
| JSON type columns use `->` and `->>` operators | ✅ |  |
| JSON type columns are wrapped in parenthesis | ✅ |  |
| SQL query syntax uses valid DuckDB syntax | ✅ |  |

</details>

<details><summary>Query title and description checks ❌</summary>

| Criteria | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Title uses title case | ✅ |  |
| Title accurately describes the query | ❌ | The query returns the top 30 objects, not 10 as stated in the title. |
| Description explains what the query does | ❌ | The description should mention it returns the top 30 most frequently accessed S3 objects, not IAM objects. |
| Description explains why a user would run the query | ❌ | Add why a user would want to see the most frequently accessed objects. |
| Description is concise | ✅ |  |

</details>

# Query Reviews

## Top 10 Requester IP Addresses ✅

<details><summary>Query</summary>
### Top 10 Requester IP Addresses

List the top 10 requester IP addresses.

```sql
select
  remote_ip,
  count(*) as request_count,
from
  aws_s3_server_access_log
group by
  remote_ip
order by
  request_count desc
limit 10;
```
</details>

<details><summary>SQL syntax checks ❌</summary>

| Criteria      | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Use 2 space indentation | ✅ |  |
| Query should end with a semicolon | ✅ |  |
| Keywords should be in lowercase | ✅ |  |
| Each clause is on its own line | ✅ |  |
| All columns exist in the schema | ✅ |  |
| STRUCT type columns use dot notation | ✅ |  |
| JSON type columns use `->` and `->>` operators | ✅ |  |
| JSON type columns are wrapped in parenthesis | ✅ |  |
| SQL query syntax uses valid DuckDB syntax | ❌ | Remove the comma after `count(*) as request_count` |

</details>

<details><summary>Query title and description checks ✅</summary>

| Criteria | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Title uses title case | ✅ |  |
| Title accurately describes the query | ✅ |  |
| Description explains what the query does | ✅ |  |
| Description explains why a user would run the query | ✅ |  |
| Description is concise | ✅ |  |

</details>

# Query Reviews

## Top Error Codes ✅

<details><summary>Query</summary>
### Top Error Codes

Identify the most frequent error codes.

```sql
select
  http_status,
  error_code,
  count(*) as error_count
from
  aws_s3_server_access_log
where
  error_code is not null
group by
  http_status,
  error_code
order by
  error_count desc;
```
</details>

<details><summary>SQL syntax checks ✅</summary>

| Criteria      | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Use 2 space indentation | ✅ |  |
| Query should end with a semicolon | ✅ |  |
| Keywords should be in lowercase | ✅ |  |
| Each clause is on its own line | ✅ |  |
| All columns exist in the schema | ✅ |  |
| STRUCT type columns use dot notation | ✅ |  |
| JSON type columns use `->` and `->>` operators | ✅ |  |
| JSON type columns are wrapped in parenthesis | ✅ |  |
| SQL query syntax uses valid DuckDB syntax | ✅ |  |

</details>

<details><summary>Query title and description checks ✅</summary>

| Criteria | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Title uses title case | ✅ |  |
| Title accurately describes the query | ✅ |  |
| Description explains what the query does | ✅ |  |
| Description explains why a user would run the query | ✅ |  |
| Description is concise | ✅ |  |

</details>

# Query Reviews

## Unusually Large File Downloads ✅

<details><summary>Query</summary>
### Unusually Large File Downloads

Detect unusually large downloads based on file size.

```sql
select
  timestamp,
  bucket,
  key,
  bytes_sent,
  operation,
  request_uri,
  requester,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  bytes_sent > 50000000 -- 50MB
  and http_status = 200
order by
  bytes_sent desc;
```
</details>

<details><summary>SQL syntax checks ✅</summary>

| Criteria      | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Use 2 space indentation | ✅ |  |
| Query should end with a semicolon | ✅ |  |
| Keywords should be in lowercase | ✅ |  |
| Each clause is on its own line | ✅ |  |
| All columns exist in the schema | ✅ |  |
| STRUCT type columns use dot notation | ✅ |  |
| JSON type columns use `->` and `->>` operators | ✅ |  |
| JSON type columns are wrapped in parenthesis | ✅ |  |
| SQL query syntax uses valid DuckDB syntax | ✅ |  |

</details>

<details><summary>Query title and description checks ✅</summary>

| Criteria | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Title uses title case | ✅ |  |
| Title accurately describes the query | ✅ |  |
| Description explains what the query does | ✅ |  |
| Description explains why a user would run the query | ✅ |  |
| Description is concise | ✅ |  |

</details>

# Query Reviews

## Requests from Unapproved IAM Roles and Users ✅

<details><summary>Query</summary>
### Requests from Unapproved IAM Roles and Users

Flag requests from IAM roles and users outside an approved list (by AWS account ID in this example).

```sql
select
  timestamp,
  bucket,
  operation,
  requester,
  remote_ip,
  user_agent
from
  aws_s3_server_access_log
where
  requester is not null -- Exclude unauthenticated requests
  and requester not like 'arn:%:%:%:123456789012:%'
order by
  timestamp desc;
```
</details>

<details><summary>SQL syntax checks ✅</summary>

| Criteria      | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Use 2 space indentation | ✅ |  |
| Query should end with a semicolon | ✅ |  |
| Keywords should be in lowercase | ✅ |  |
| Each clause is on its own line | ✅ |  |
| All columns exist in the schema | ✅ |  |
| STRUCT type columns use dot notation | ✅ |  |
| JSON type columns use `->` and `->>` operators | ✅ |  |
| JSON type columns are wrapped in parenthesis | ✅ |  |
| SQL query syntax uses valid DuckDB syntax | ✅ |  |

</details>

<details><summary>Query title and description checks ✅</summary>

| Criteria | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Title uses title case | ✅ |  |
| Title accurately describes the query | ✅ |  |
| Description explains what the query does | ✅ |  |
| Description explains why a user would run the query | ✅ |  |
| Description is concise | ✅ |  |

</details>