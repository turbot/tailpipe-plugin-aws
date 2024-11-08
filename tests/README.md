# Test data

The `generator` folder contains a script, `generate.py`, which writes ~500,000 records to a parquet file. To run the sample queries in [docs/tables](../docs/tables), cd into `generator` and run:

```bash
$ python generate.py
Generated 500000 records
```

The output is `alb_access_log.parquet`.

In `generator`, run DuckDB.

```bash
$ duckdb
v1.1.3 19864453f7
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
D CREATE VIEW alb_access_log AS SELECT * FROM read_parquet('alb_access_log.parquet');
```

You can copy queries from the table docs and paste them here.

```sql
WITH suspicious_paths AS (
          SELECT
              client_ip,
              user_agent,
              alb_name,
              request as sample_request,
              timestamp,
              ROW_NUMBER() OVER (PARTITION BY client_ip, alb_name ORDER BY timestamp) AS row_num
          FROM alb_access_log
          WHERE
              request LIKE '%actuator%' OR
              request LIKE '%metrics%' OR
              request LIKE '%phpinfo%' OR
              request LIKE '%server-status%' OR
              request LIKE '%jndi:ldap%' OR
              request LIKE '%class.module.classLoader%' OR
              request LIKE '%.env%' OR
              request LIKE '%wp-config%' OR
              request LIKE '%/debug%'
      )
      SELECT
          client_ip,
          user_agent,
          COUNT(DISTINCT alb_name) as albs_targeted,
          COUNT(*) as total_probes,
          STRING_AGG(DISTINCT alb_name, ', ') as targeted_albs,
          STRING_AGG(sample_request, ' | ') as sample_requests,
          MIN(timestamp) as first_seen,
          MAX(timestamp) as last_seen,
          EXTRACT(MINUTES FROM MAX(timestamp) - MIN(timestamp)) as campaign_duration_mins
      FROM suspicious_paths
      WHERE row_num <= 3
      GROUP BY client_ip, user_agent
      HAVING
          COUNT(DISTINCT alb_name) > 1 AND  -- Targeting multiple ALBs
          COUNT(*) >= 3                     -- At least 3 probe attempts
      ORDER BY albs_targeted DESC, total_probes DESC
      LIMIT 10;
```

```
┌─────────────────┬──────────────────────┬───┬──────────────────────┬──────────────────────┐
│    client_ip    │      user_agent      │ … │      last_seen       │ campaign_duration_…  │
│     varchar     │       varchar        │   │     timestamp_ns     │        int64         │
├─────────────────┼──────────────────────┼───┼──────────────────────┼──────────────────────┤
│ 185.181.233.81  │ Nmap Scripting Eng…  │ … │ 2024-11-05 17:27:2…  │                   57 │
│ 185.181.92.9    │ Nuclei/2.9.1 (http…  │ … │ 2024-11-06 08:20:2…  │                   35 │
│ 185.181.132.138 │ gobuster/3.5         │ … │ 2024-11-01 01:52:4…  │                   40 │
│ 193.27.228.154  │ Acunetix-Agent       │ … │ 2024-11-02 21:40:2…  │                   52 │
│ 185.181.171.121 │ gobuster/3.5         │ … │ 2024-11-06 17:00:0…  │                    2 │
│ 185.181.41.10   │ Qualys SSL Assessm…  │ … │ 2024-11-02 20:05:4…  │                   21 │
│ 193.27.228.243  │ dirbuster/1.0-RC1    │ … │ 2024-11-02 19:08:1…  │                   26 │
│ 45.155.205.93   │ WhatWeb/0.5.5        │ … │ 2024-11-01 00:03:4…  │                    2 │
│ 193.27.228.200  │                      │ … │ 2024-11-01 17:52:0…  │                    7 │
│ 45.155.205.16   │ zgrab/0.x            │ … │ 2024-11-02 15:19:0…  │                   17 │
├─────────────────┴──────────────────────┴───┴──────────────────────┴──────────────────────┤
│ 10 rows                                                              9 columns (4 shown) │
└──────────────────────────────────────────────────────────────────────────────────────────┘
```
