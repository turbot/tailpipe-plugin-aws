import re
import sys
import duckdb

# Define the schema columns for aws_s3_server_access_log
S3_LOG_COLUMNS = {
    "bucket_owner", "bucket", "requester", "request_time", "operation", "key",
    "request_uri", "http_status", "error_code", "bytes_sent", "object_size",
    "total_time", "turn_around_time", "referrer", "user_agent", "version_id",
    "host_id", "signature_version"
}

# SQL Query Evaluation Criteria
EVALUATION_CRITERIA = [
    ("2-space indentation", lambda q: all(line.startswith('  ') or line == '' for line in q.split('\n'))),
    ("Ends with semicolon", lambda q: q.strip().endswith(";")),
    ("Keywords in uppercase", lambda q: not any(kw.lower() in q for kw in ["select", "from", "group by", "order by"])),
    ("Each clause on a new line", lambda q: all(re.search(rf"\b{kw}\b", q, re.IGNORECASE) for kw in ["SELECT", "FROM", "GROUP BY", "ORDER BY"])),
    ("Valid column names", lambda q: all(col in S3_LOG_COLUMNS for col in re.findall(r"\b\w+\b", q))),
    ("DuckDB syntax valid", lambda q: is_valid_duckdb(q))
]

def is_valid_duckdb(query):
    """Run the SQL query in DuckDB to check syntax validity."""
    try:
        duckdb.sql(query)
        return True
    except Exception:
        return False

def evaluate_query(query):
    """Run the evaluation and return a formatted review."""
    results = []
    total_criteria = len(EVALUATION_CRITERIA)
    passed_criteria = 0

    for name, check in EVALUATION_CRITERIA:
        passed = check(query)
        if passed:
            passed_criteria += 1
        results.append(f"- **{name}**: {'✅ Pass' if passed else '❌ Fail'}")

    final_score = f"**Final Score: {passed_criteria}/{total_criteria} ({(passed_criteria/total_criteria)*100:.0f}%)**"
    
    # Return formatted evaluation
    return "\n".join(results) + "\n\n" + final_score

def extract_queries(md_file):
    """Extract SQL queries from a queries.md file."""
    with open(md_file, "r") as file:
        content = file.read()

    queries = re.findall(r"```sql\n(.*?)\n```", content, re.DOTALL)
    return queries

def main():
    """Main function to process queries and evaluate them."""
    if len(sys.argv) < 2:
        print("Usage: python evaluate_sql.py <queries.md>")
        sys.exit(1)

    md_file = sys.argv[1]
    queries = extract_queries(md_file)

    if not queries:
        print("No SQL queries found in the provided file.")
        sys.exit(0)

    output = []
    for i, query in enumerate(queries, 1):
        output.append(f"### Query {i} Review\n```sql\n{query.strip()}\n```\n")
        output.append(evaluate_query(query.strip()))

    print("\n".join(output))

if __name__ == "__main__":
    main()
