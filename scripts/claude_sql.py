import os
import sys
import re
import time
import random
import anthropic
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get API key from environment
CLAUDE_API_KEY = os.getenv("ANTHROPIC_API_KEY")
if not CLAUDE_API_KEY:
    print(
        "❌ Claude API key is missing. Set the ANTHROPIC_API_KEY environment variable."
    )
    sys.exit(1)

# Initialize Claude client
client = anthropic.Anthropic(api_key=CLAUDE_API_KEY)

# Claude model to use
MODEL = "claude-3-5-sonnet-20240620"

# Get table name from environment or use default
TABLE_NAME = os.getenv("TABLE_NAME", "aws_s3_server_access_log")


def extract_queries(file_path):
    """Extract SQL queries along with their titles and descriptions from a queries.md file."""
    with open(file_path, "r") as f:
        content = f.read()

    queries = re.findall(r"###\s*(.*?)\s*\n(.*?)```sql\n(.*?)```", content, re.DOTALL)
    return [
        {"title": title.strip(), "description": desc.strip(), "query": query.strip()}
        for title, desc, query in queries
    ]


def generate_evaluation_prompt(query_data, schema):
    return f"""
For the Tailpipe table `{TABLE_NAME}`, with the schema:

```go
{schema}
```

Evaluate this SQL query against these criteria:

```markdown
### {query_data['title']}

{query_data['description']}

```sql
{query_data['query']}
```
```

Evaluate the query against each of these specific criteria sets:

# SQL syntax checks
1. Use 2 space indentation.
2. Query should end with a semicolon.
3. Keywords should be in lowercase.
4. Each clause (SELECT, FROM, WHERE, GROUP BY, ORDER BY) is on its own line.
5. All columns should exist in the schema.
6. STRUCT type columns should use dot notation when accessing properties.
7. JSON type columns should use `->` and `->>` operators when accessing properties.
8. JSON type columns using `->` and `->>` operators should be wrapped in parenthesis to avoid operator precedence issues.
9. There should be a space before and after each `->` and `->>`.
10. SQL query syntax uses valid DuckDB syntax.

# Title and description checks
1. The query's title should use title case.
2. The query's title should accurately describe what the query does.
3. If the query contains `limit X`, that number should reflected in the title, e.g., `Top 10 Expensive Services` with `limit 10`.
4. The first sentence of the query description should explain what the query does.
5. The second sentence of the query description should explain why a user would want to run the query.
6. Each sentence in the query description should be concise.

# Query relevance to logs checks
1. The query should provide useful insights for the specific log type it analyzes.
2. The query should be relevant to security, operational, or performance monitoring use cases.

# Column selection checks
1. For aggregated queries, timestamp related columns should not be included.
2. For non-aggregated queries, the first column in the `SELECT` statement should be the event timestamp.
3. For non-aggregated queries, if the log contains information on where resources exist (e.g., `account_id` and `region` for AWS CloudTrail logs, `subscription_id` and `resource_group_name` for Azure activity logs), include those columns.
4. For non-aggregated queries, the columns related to where the resources exist should be the last columns in the `SELECT` statement.
5. If the query's `WHERE` clause contains a specific value lookup, e.g., `where elb_status_code = 502`, do not include that column in the `SELECT` statement to avoid including redundant information.

# Sorting strategy checks
1. For non-aggregated queries, the default ordering should be `<event timestamp column> desc` so the most recent log data is returned first.
   - However, if another ordering gives more relevant results, e.g., `select timestamp, bucket, key, bytes_sent, operation from aws_s3_server_access_log where bytes_sent > 50000 order by bytes_sent desc;`, that is acceptable.
2. For aggregated queries, the default ordering should be the main count in descending order, e.g., `select client_ip, count(*) as request_count from aws_alb_access_log group by client_ip order by response_count desc limit 10`.
   - However, if the query is looking at trends over dates or times, the query can be ordered by the relevant time column in ascending order.

You MUST provide ONLY the following output format, with no additional text:

```markdown
## {query_data['title']} [OVERALL_MARK]

<details><summary>Query</summary>

### {query_data['title']}

{query_data['description']}

```sql
{query_data['query']}
```
</details>

<details><summary>SQL syntax checks [SQL_CHECKS_MARK]</summary>

| Criteria | Pass/Fail | Suggestions |
|----------|-----------|-------------|
| Use 2 space indentation | [MARK] | [SUGGESTION] |
| Query should end with a semicolon | [MARK] | [SUGGESTION] |
| Keywords should be in lowercase | [MARK] | [SUGGESTION] |
| Each clause is on its own line | [MARK] | [SUGGESTION] |
| All columns exist in the schema | [MARK] | [SUGGESTION] |
| STRUCT type columns use dot notation | [MARK] | [SUGGESTION] |
| JSON type columns use `->` and `->>` operators | [MARK] | [SUGGESTION] |
| JSON type columns are wrapped in parenthesis | [MARK] | [SUGGESTION] |
| Space before and after each `->` and `->>` | [MARK] | [SUGGESTION] |
| SQL query syntax uses valid DuckDB syntax | [MARK] | [SUGGESTION] |

</details>

<details><summary>Title and description checks [TITLE_CHECKS_MARK]</summary>

| Criteria | Pass/Fail | Suggestions |
|----------|-----------|-------------|
| Title uses title case | [MARK] | [SUGGESTION] |
| Title accurately describes the query | [MARK] | [SUGGESTION] |
| Title contains limit value if in query | [MARK] | [SUGGESTION] |
| Description explains what the query does | [MARK] | [SUGGESTION] |
| Description explains why a user would run the query | [MARK] | [SUGGESTION] |
| Description is concise | [MARK] | [SUGGESTION] |

</details>

<details><summary>Query relevance checks [RELEVANCE_CHECKS_MARK]</summary>

| Criteria | Pass/Fail | Suggestions |
|----------|-----------|-------------|
| Provides useful insights for this log type | [MARK] | [SUGGESTION] |
| Relevant to security, operational, or performance monitoring | [MARK] | [SUGGESTION] |

</details>

<details><summary>Column selection checks [COLUMN_CHECKS_MARK]</summary>

| Criteria | Pass/Fail | Suggestions |
|----------|-----------|-------------|
| Aggregated queries should not include timestamp columns | [MARK] | [SUGGESTION] |
| Non-aggregated queries should have the log's event timestamp as the first column | [MARK] | [SUGGESTION] |
| Non-aggregated queries should include columns related to where the resources exist | [MARK] | [SUGGESTION] |
| Non-aggregated queries should place columns related to where the resources exist last | [MARK] | [SUGGESTION] |
| Avoid selecting columns with fixed values in WHERE clause | [MARK] | [SUGGESTION] |

</details>

<details><summary>Sorting strategy checks [SORTING_CHECKS_MARK]</summary>

| Criteria | Pass/Fail | Suggestions |
|----------|-----------|-------------|
| Non-aggregated queries default to `<event timestamp column> desc` | [MARK] | [SUGGESTION] |
| Aggregated queries ordered by count desc or timestamp asc | [MARK] | [SUGGESTION] |

</details>
```

Where:
- Replace [OVERALL_MARK] with ✅ if all checks pass, or ❌ if any checks fail
- Replace [SQL_CHECKS_MARK], [TITLE_CHECKS_MARK], [RELEVANCE_CHECKS_MARK], [COLUMN_CHECKS_MARK], and [SORTING_CHECKS_MARK] with ✅ if all checks in that section pass, or ❌ if any fail
- For each criteria row, replace [MARK] with either ✅ for pass or ❌ for fail
- Replace [SUGGESTION] with specific improvement suggestions if the check fails, or leave it blank if it passes

Note: If a specific check is not applicable (N/A) for this query, mark it as ✅ and note "N/A" in the suggestions.
"""


def evaluate_query_with_retries(query_data, schema, max_retries=5):
    """Calls Claude API to evaluate the SQL query with retries and exponential backoff."""
    retries = 0
    base_delay = 2  # Start with a 2-second delay
    max_delay = 60  # Cap at 1 minute

    while retries <= max_retries:
        try:
            prompt = generate_evaluation_prompt(query_data, schema)

            response = client.messages.create(
                model=MODEL,
                max_tokens=4000,
                temperature=0,
                system="You are an expert SQL validator who strictly follows the output format provided. You will evaluate SQL queries against specific technical criteria and provide feedback in EXACTLY the format requested, with no deviations or additional text.",
                messages=[{"role": "user", "content": prompt}],
            )

            if response.content:
                evaluation_result = response.content[0].text

                # Extract just the markdown output (in case Claude adds extra text)
                match = re.search(
                    r"```markdown\s*(## .*?)\s*```", evaluation_result, re.DOTALL
                )
                if match:
                    evaluation_result = match.group(1)
                else:
                    # If no markdown block is found, try to extract just the review section
                    match = re.search(r"(## .*)", evaluation_result, re.DOTALL)
                    if match:
                        evaluation_result = match.group(1)

                return evaluation_result
            else:
                return (
                    "❌ Claude response was empty or malformed. No insights generated."
                )

        except anthropic.APIError as e:
            retries += 1

            # Check if it's an overloaded error (429/529)
            error_code = getattr(e, "status_code", 0)
            is_overloaded = error_code in [429, 529] or (
                hasattr(e, "error")
                and getattr(e.error, "type", "") == "overloaded_error"
            )

            if is_overloaded and retries <= max_retries:
                # Calculate backoff delay with jitter
                delay = min(max_delay, base_delay * (2 ** (retries - 1)))
                jitter = random.uniform(0, 0.1 * delay)  # Add 0-10% jitter
                wait_time = delay + jitter

                print(
                    f"⚠️ Claude API overloaded. Retry {retries}/{max_retries} after {wait_time:.2f}s..."
                )
                time.sleep(wait_time)
                continue

            # If we've exhausted retries or it's not an overloaded error
            if retries > max_retries:
                return f"❌ Max retries exceeded ({max_retries}). Claude API remains overloaded."
            else:
                return f"❌ Error evaluating query: {str(e)}"

        except Exception as e:
            return f"❌ Error evaluating query: {str(e)}"


def main():
    if len(sys.argv) < 2:
        print("❌ No file provided for evaluation.")
        print("Usage: python claude_sql.py input_file.md")
        sys.exit(1)

    file_path = sys.argv[1]

    # Use schema from environment variable if available, otherwise use default
    schema_content = os.getenv("SCHEMA_CONTENT")

    # Define the schema for S3 Server Access Log (default)
    schema = (
        schema_content
        if schema_content
        else """type S3ServerAccessLog struct {
    schema.CommonFields
    AccessPointArn     *string   `json:"access_point_arn,omitempty"`
    AclRequired        *bool     `json:"acl_required,omitempty"`
    AuthenticationType *string   `json:"authentication_type,omitempty"`
    Bucket             string    `json:"bucket"`
    BucketOwner        string    `json:"bucket_owner"`
    BytesSent          *int64    `json:"bytes_sent,omitempty"`
    CipherSuite        *string   `json:"cipher_suite,omitempty"`
    ErrorCode          *string   `json:"error_code,omitempty"`
    HTTPStatus         *int      `json:"http_status"`
    HostHeader         *string   `json:"host_header,omitempty"`
    HostID             *string   `json:"host_id,omitempty"`
    Key                *string   `json:"key,omitempty"`
    ObjectSize         *int64    `json:"object_size,omitempty"`
    Operation          string    `json:"operation"`
    Referer            *string   `json:"referer,omitempty"`
    RemoteIP           string    `json:"remote_ip"`
    RequestID          string    `json:"request_id"`
    RequestURI         *string   `json:"request_uri"`
    Requester          string    `json:"requester,omitempty"`
    SignatureVersion   *string   `json:"signature_version,omitempty"`
    TLSVersion         *string   `json:"tls_version,omitempty"`
    Timestamp          time.Time `json:"timestamp"`
    TotalTime          *int      `json:"total_time"`
    TurnAroundTime     *int      `json:"turn_around_time,omitempty"`
    UserAgent          *string   `json:"user_agent,omitempty"`
    VersionID          *string   `json:"version_id,omitempty"`
}"""
    )

    try:
        # Check if file exists
        if not os.path.isfile(file_path):
            print(f"❌ File not found: {file_path}")
            sys.exit(1)

        print(f"📂 Reading queries from {file_path}...")
        queries = extract_queries(file_path)

        if not queries:
            print("❌ No SQL queries found in the provided file.")
            sys.exit(1)

        print(f"✅ Found {len(queries)} queries to evaluate.")
        print(f"🏷️ Using table name: {TABLE_NAME}")
        all_reviews = []

        for i, query in enumerate(queries):
            print(f"🔍 Evaluating ({i+1}/{len(queries)}): {query['title']} ...")

            # Use the new function with retries
            review = evaluate_query_with_retries(query, schema)
            all_reviews.append(review)

            # Add a small delay to avoid hitting API rate limits
            if i < len(queries) - 1:
                time.sleep(2)  # Increased delay to reduce rate limiting

        combined_reviews = "\n\n".join(all_reviews)

        output_file = "review_output.md"
        with open(output_file, "w") as f:
            f.write(combined_reviews)

        print(f"✅ Query evaluation complete. Results saved in {output_file}")

    except Exception as e:
        print(f"❌ Error processing file: {str(e)}")
        sys.exit(1)


if __name__ == "__main__":
    main()
