import os
import sys
import re
import time
import anthropic
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get API key from environment
CLAUDE_API_KEY = os.getenv("ANTHROPIC_API_KEY")
if not CLAUDE_API_KEY:
    print("❌ Claude API key is missing. Set the ANTHROPIC_API_KEY environment variable.")
    sys.exit(1)

# Initialize Claude client
client = anthropic.Anthropic(api_key=CLAUDE_API_KEY)

# Claude model to use
MODEL = "claude-3-opus-20240229"

def extract_queries(file_path):
    """Extract SQL queries along with their titles and descriptions from a queries.md file."""
    with open(file_path, "r") as f:
        content = f.read()
    
    queries = re.findall(r"###\s*(.*?)\s*\n(.*?)```sql\n(.*?)```", content, re.DOTALL)
    return [{"title": title.strip(), "description": desc.strip(), "query": query.strip()} for title, desc, query in queries]

def generate_evaluation_prompt(query_data, schema):
    return f"""
For the Tailpipe table `aws_s3_server_access_log`, with the schema:

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

You MUST provide ONLY the following output format, with no additional text:

```markdown
# Query Reviews

## {query_data['title']} [OVERALL_MARK]

<details><summary>Query</summary>
### {query_data['title']}

{query_data['description']}

```sql
{query_data['query']}
```
</details>

<details><summary>SQL syntax checks [SQL_CHECKS_MARK]</summary>

| Criteria      | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Use 2 space indentation | [MARK] | [SUGGESTION] |
| Query should end with a semicolon | [MARK] | [SUGGESTION] |
| Keywords should be in lowercase | [MARK] | [SUGGESTION] |
| Each clause is on its own line | [MARK] | [SUGGESTION] |
| All columns exist in the schema | [MARK] | [SUGGESTION] |
| STRUCT type columns use dot notation | [MARK] | [SUGGESTION] |
| JSON type columns use `->` and `->>` operators | [MARK] | [SUGGESTION] |
| JSON type columns are wrapped in parenthesis | [MARK] | [SUGGESTION] |
| SQL query syntax uses valid DuckDB syntax | [MARK] | [SUGGESTION] |

</details>

<details><summary>Query title and description checks [TITLE_CHECKS_MARK]</summary>

| Criteria | Pass/Fail | Suggestions |
|---------------|-----------|-------------|
| Title uses title case | [MARK] | [SUGGESTION] |
| Title accurately describes the query | [MARK] | [SUGGESTION] |
| Description explains what the query does | [MARK] | [SUGGESTION] |
| Description explains why a user would run the query | [MARK] | [SUGGESTION] |
| Description is concise | [MARK] | [SUGGESTION] |

</details>
```

Where:
- Replace [OVERALL_MARK] with ✅ if all checks pass, or ❌ if any checks fail
- Replace [SQL_CHECKS_MARK] with ✅ if all SQL syntax checks pass, or ❌ if any fail
- Replace [TITLE_CHECKS_MARK] with ✅ if all title/description checks pass, or ❌ if any fail
- For each criteria row, replace [MARK] with either ✅ for pass or ❌ for fail
- Replace [SUGGESTION] with specific improvement suggestions if the check fails, or leave it blank if it passes

Give special attention to these checks:
1. Check if the title matches what the query actually does (count, LIMIT value, etc.)
2. Look for SQL syntax errors like extra commas between columns
3. Check for inconsistent indentation
4. Verify all columns actually exist in the schema
5. Ensure each clause (SELECT, FROM, WHERE, etc.) is on its own line

The goal is to provide a technical review of the query that checks both SQL syntax correctness and the quality/accuracy of the title and description.
"""

def evaluate_query(query_data, schema):
    """Calls Claude API to evaluate the SQL query."""
    prompt = generate_evaluation_prompt(query_data, schema)
    
    try:
        response = client.messages.create(
            model=MODEL,
            max_tokens=4000,
            temperature=0,
            system="You are an expert SQL validator who strictly follows the output format provided. You will evaluate SQL queries against specific technical criteria and provide feedback in EXACTLY the format requested, with no deviations or additional text.",
            messages=[
                {"role": "user", "content": prompt}
            ]
        )
        
        if response.content:
            evaluation_result = response.content[0].text
            
            # Extract just the markdown output (in case Claude adds extra text)
            match = re.search(r'```markdown\s*(# Query Reviews.*?)\s*```', evaluation_result, re.DOTALL)
            if match:
                evaluation_result = match.group(1)
            else:
                # If no markdown block is found, try to extract just the review section
                match = re.search(r'(# Query Reviews.*)', evaluation_result, re.DOTALL)
                if match:
                    evaluation_result = match.group(1)
            
            return evaluation_result
        else:
            return "❌ Claude response was empty or malformed. No insights generated."
    except Exception as e:
        return f"❌ Error evaluating query: {str(e)}"

def main():
    if len(sys.argv) < 2:
        print("❌ No file provided for evaluation.")
        print("Usage: python claude_query_evaluator.py input_file.md")
        sys.exit(1)

    file_path = sys.argv[1]
    
    # Define the schema for S3 Server Access Log
    schema = """type S3ServerAccessLog struct {
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
        all_reviews = []
        
        for i, query in enumerate(queries):
            print(f"🔍 Evaluating ({i+1}/{len(queries)}): {query['title']} ...")
            review = evaluate_query(query, schema)
            all_reviews.append(review)
            
            # Add a small delay to avoid hitting API rate limits
            if i < len(queries) - 1:
                time.sleep(1)
        
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