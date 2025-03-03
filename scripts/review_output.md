# Query Reviews

## Unauthenticated requests ❌

  <details><summary>Query</summary>
  ### Unauthenticated requests

  ```sql
  
select
timestamp,
bucket,
operation,
request_uri,
from
  aws_s3_server_access_log
where
  requester is null
order by
  timestamp desc;

  ```
  </details>

  <details><summary>SQL syntax checks ❌</summary>

  | Criteria      | Pass/Fail | Suggestions |
  |---------------|-----------|-------------|
  | Use 2 space indentation | ✅ |             |
  | Query should end with a semicolon | ✅ |             |
  | Keywords should be in lowercase | ✅ |             |
  | Each clause is on its own line | ✅ |             |
  | All columns exist in the schema | ✅ |             |
  | STRUCT type columns use dot notation | ✅ |             |
  | JSON type columns use `->` and `->>` operators | ✅ |             |
  | JSON type columns are wrapped in parenthesis | ✅ |             |
  | SQL query syntax uses valid DuckDB syntax | ❌ | The trailing comma before `from` should be removed. |

  </details>

  <details><summary>Query title and description checks ❌</summary>

  | Criteria | Pass/Fail | Suggestions |
  |---------------|-----------|-------------|
  | Title uses title case | ✅ |             |
  | Title accurately describes the query | ✅ |             |
  | Description explains what the query does | ❌ | The description should clarify that it lists unauthenticated requests. |
  | Description explains why a user would run the query | ❌ | The description should mention potential use cases, such as monitoring or auditing purposes. |
  | Description is concise | ✅ |             |

  </details>