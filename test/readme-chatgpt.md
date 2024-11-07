
# `aws_alb_access_log` Table - Tailpipe Plugin

The `aws_alb_access_log` table is a component of the Tailpipe plugin that enables SQL-based analysis of AWS ALB (Application Load Balancer) logs. Built on Tailpipe’s Plugin SDK, it uses the SDK's services to streamline data collection, parsing, and querying. This document covers the plugin's architecture, usage, and a detailed breakdown of the implementation files.

## Overview

The `aws_alb_access_log` table’s implementation is modular, with each file addressing a specific responsibility. This structure makes it easier to maintain, extend, and integrate new data sources.

### Key Components

1. **Data Collection and Parsing**: Extracts, parses, and structures log data into a format suitable for SQL querying.
2. **Configuration and Partitioning**: Allows customization of data sources and logical grouping of logs.
3. **SQL Query Interface**: Provides SQL-like syntax for analyzing data, enabling users to categorize, filter, and aggregate log entries efficiently.

## Files and Their Roles

The `aws_alb_access_log` table’s functionality is distributed across several files, each responsible for a distinct part of the implementation.

### 1. `alb_generate.py` (Test Data Generator)

   **Purpose**: This script is used for generating synthetic test data in the format expected by the `aws_alb_access_log` table. This log data can be used to validate the table’s parsing logic, schema, and query capabilities.

   **Key Functions**:
   - Generates log entries with simulated timestamps, IP addresses, HTTP methods, response codes, and user-agent strings.
   - Ensures that test data matches the schema defined in the `aws_alb_access_log` table, supporting reliable testing and debugging.

### 2. `aws_alb_access_log_config.hcl` (Configuration File)

   **Purpose**: Defines the plugin settings and log file paths for the `aws_alb_access_log` table. This HCL (HashiCorp Configuration Language) file specifies data source locations, file types, and partition names, allowing Tailpipe to locate and parse log files effectively.

   **Key Components**:
   - **Partition Section**: Specifies partitions (or logical data groupings) for each set of logs. In this case, each `partition` block points to a file path for logs (e.g., `"alb_test"` partition pointing to a `.log` file).
   - **Source Definition**: Indicates that the log data source is the filesystem, with file extensions restricted to `.log`.
   - **Plugin Setting**: Associates the `aws_alb_access_log` table with the AWS plugin, enabling it to process and parse ALB log files specifically.

### 3. `aws_alb_access_log_table.go` (Table Definition)

   **Purpose**: This Go file defines the schema and core logic for the `aws_alb_access_log` table. It specifies each column, data type, and parsing rule to structure raw log data into a standardized format.

   **Key Components**:
   - **Schema Definition**: Lists each column (e.g., `user_agent`, `timestamp`, `client_ip`) along with its data type (e.g., `string`, `timestamp`). This schema is critical for querying since it maps each log field to a structured format.
   - **Parsing Logic**: Implements functions that parse each log entry, extracting relevant fields and converting them into columns. For instance, it may parse `user_agent` strings to categorize the client type as `Mobile`, `Chrome`, `Firefox`, etc.
   - **Type Mapping**: Maps each field from the raw logs to a specific SQL-compatible type, ensuring consistency and compatibility in queries.

### 4. `aws_alb_access_log_parser.go` (Log Parsing Logic)

   **Purpose**: This file contains the core logic for parsing raw ALB log entries, breaking each entry into fields according to the AWS ALB log format.

   **Key Components**:
   - **Field Extraction**: Implements functions that recognize and extract fields (e.g., `client_ip`, `method`, `path`, `status_code`) from each log line.
   - **Pattern Matching**: Uses regular expressions or delimiter-based parsing to separate fields within each log entry.
   - **Error Handling**: Accounts for incomplete or corrupted log lines, applying default values or skipping entries as necessary.
   - **Normalization**: Normalizes data to standard formats (e.g., `timestamp` formats, IP addresses) for consistency across log entries.

### 5. `aws_alb_access_log_test.go` (Unit Tests)

   **Purpose**: This file includes unit tests for the `aws_alb_access_log` table. It validates the functionality of data parsing, schema conformance, and ensures that each function behaves as expected.

   **Key Components**:
   - **Schema Validation**: Verifies that log entries are parsed into the correct schema (e.g., all columns have the expected types).
   - **Parsing Accuracy**: Checks the correctness of field extraction, ensuring each log entry’s data maps accurately to columns.
   - **Error Scenarios**: Tests for edge cases like malformed log entries, missing fields, and unexpected data types to ensure the parser handles them gracefully.

## Tailpipe Workflow: Collection and Query Phases

Tailpipe operates in two main phases for processing log data: **collection** and **query**.

### Collection Phase

The collection phase is where Tailpipe extracts raw log data, converts it into a structured format, and stores it as Parquet files for efficient querying. Here’s how it works for the `aws_alb_access_log` table:

1. **Source Definition**: The configuration file (`aws_alb_access_log_config.hcl`) defines the log source path, file extension, and partition metadata.
2. **Data Parsing**: `aws_alb_access_log_parser.go` parses each log entry according to the schema, converting fields like `user_agent` into SQL-compatible columns.
3. **Parquet Conversion**: Once parsed, data is stored in Parquet format, allowing for quick retrieval and reduced storage overhead.

   Example command to initiate data collection:
   
   ```bash
   TAILPIPE_LOG_LEVEL=debug ./tailpipe collect aws_alb_access_log.alb_test
   ```

### Query Phase

The query phase enables SQL-based access to the data. Tailpipe’s SQL engine supports a variety of operations, including filtering, aggregation, and ordering. For instance, the following query categorizes requests by client type and counts them:

```sql
SELECT
    CASE
        WHEN user_agent LIKE '%Mobile%' THEN 'Mobile'
        WHEN user_agent LIKE '%Chrome%' THEN 'Chrome'
        WHEN user_agent LIKE '%Firefox%' THEN 'Firefox'
        WHEN user_agent LIKE '%Safari%' THEN 'Safari'
        WHEN user_agent LIKE '%bot%' OR user_agent LIKE '%Bot%' THEN 'Bot'
        ELSE 'Other'
    END as client_type,
    COUNT(*) as request_count
FROM aws_alb_access_log
GROUP BY client_type
ORDER BY request_count DESC;
```

### Example Output

```plaintext
┌─────────────┬───────────────┐
│ client_type │ request_count │
├─────────────┼───────────────┤
│ Other       │          6337 │
│ Mobile      │          1842 │
│ Chrome      │          1821 │
└─────────────┴───────────────┘
```

This output provides insights into the distribution of requests across client types, useful for identifying usage patterns.

## Summary

This `aws_alb_access_log` implementation is modular and efficient, providing flexible log analysis via SQL queries. Each file serves a specific role, from configuration and parsing to testing and schema management. This structure not only makes the table easy to use but also allows for quick adaptation to different data sources or query requirements.
