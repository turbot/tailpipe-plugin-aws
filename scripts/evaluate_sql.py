import os
import sys
import re
import openai
import time

# Ensure OpenAI API key is set as an environment variable
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
if not OPENAI_API_KEY:
    print("❌ OpenAI API key is missing. Set the OPENAI_API_KEY environment variable.")
    sys.exit(1)

client = openai.Client(api_key=OPENAI_API_KEY)  # Corrected API client initialization

# OpenAI Model
MODEL = "gpt-4o-mini"

# SQL Query Evaluation Prompt Template
PROMPT_TEMPLATE = """For the Tailpipe table `{table_name}`, with the schema:

```go
{schema}
```

Can you please evaluate this SQL query:

```sql
{query}
```

Using these exact evaluation criteria and output format:

# Query Reviews

## {query_title} ✅/❌

  <details><summary>Query</summary>
  ### {query_title}

  {query_description}

  ```sql
  {query}
  ```
  </details>

  <details><summary>SQL syntax checks ✅/❌</summary>

  | Criteria      | Pass/Fail | Suggestions |
  |---------------|-----------|-------------|
  | Use 2 space indentation | ✅/❌ |             |
  | Query should end with a semicolon | ✅/❌ |             |
  | Keywords should be in lowercase | ✅/❌ |             |
  | Each clause is on its own line | ✅/❌ |             |
  | All columns exist in the schema | ✅/❌ |             |
  | STRUCT type columns use dot notation | ✅/❌ |             |
  | JSON type columns use `->` and `->>` operators | ✅/❌ |             |
  | JSON type columns are wrapped in parenthesis | ✅/❌ |             |
  | SQL query syntax uses valid DuckDB syntax | ✅/❌ |             |

  </details>

  <details><summary>Query title and description checks ✅/❌</summary>

  | Criteria | Pass/Fail | Suggestions |
  |---------------|-----------|-------------|
  | Title uses title case | ✅/❌ | Suggestion |
  | Title accurately describes the query | ✅/❌ | Suggestion |
  | Description explains what the query does | ✅/❌ | Suggestion |
  | Description explains why a user would run the query | ✅/❌ | Suggestion |
  | Description is concise | ✅/❌ | Suggestion |

  </details>
"""

# S3 Server Access Log Schema
SCHEMA = """type S3ServerAccessLog struct {
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

def extract_queries(file_path):
    """Extract SQL queries along with their descriptions from a queries.md file."""
    with open(file_path, "r") as f:
        content = f.read()
    
    queries = re.findall(r"###\s*(.*?)\s*\n(.*?)\n+```sql\n(.*?)\n```", content, re.DOTALL)
    return [{"title": title.strip().title(), "description": desc.strip(), "query": query.strip()} for title, desc, query in queries]

def evaluate_query(query_data):
    """Calls OpenAI's GPT API to evaluate the SQL query."""
    table_name = "aws_s3_server_access_log"
    prompt = PROMPT_TEMPLATE.format(
        table_name=table_name,
        schema=SCHEMA,
        query=query_data["query"],
        query_title=query_data["title"],
        query_description=query_data["description"]
    )

    try:
        response = client.chat.completions.create(
            model=MODEL,
            messages=[
                {"role": "system", "content": "You are an expert in SQL query validation and DuckDB."},
                {"role": "user", "content": prompt}
            ],
            temperature=0.5
        )

        if response.choices:
            evaluation_result = response.choices[0].message.content.strip()
        else:
            return "❌ OpenAI response was empty or malformed. No insights generated."

        return evaluation_result
    except Exception as e:
        return f"❌ Error evaluating query: {str(e)}"


def main():
    if len(sys.argv) < 2:
        print("❌ No file provided for evaluation.")
        sys.exit(1)

    file_path = sys.argv[1]
    queries = extract_queries(file_path)

    if not queries:
        print("❌ No SQL queries found in the provided file.")
        sys.exit(1)

    reviews = []
    for query in queries:
        print(f"🔍 Evaluating: {query['title']} ...")
        review = evaluate_query(query)
        reviews.append(review)

    with open("review_output.md", "w") as f:
        f.write("\n".join(reviews))
    
    print("✅ Query evaluation complete. Results saved in `review_output.md`.")

if __name__ == "__main__":
    main()
