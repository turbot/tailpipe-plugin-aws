---
title: "Tailpipe Table: aws_network_firewall_log - Query AWS Network Firewall Logs"
description: "AWS Network Firewall logs capture information about traffic flowing through your AWS Network Firewall, including flow logs and alert information."
---

# Table: aws_network_firewall_log - Query AWS Network Firewall Logs

The `aws_network_firewall_log` table allows you to query data from AWS Network Firewall logs. This table provides detailed insights into the traffic flowing through your AWS Network Firewall, including source and destination IP addresses, ports, protocols, and alert information.

## Configure

Create a [partition](https://tailpipe.io/docs/manage/partition) for `aws_network_firewall_log` ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_network_firewall_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "network_firewall_logging" {
  profile = "my-network-firewall-logging"
}

partition "aws_network_firewall_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.network_firewall_logging
    bucket     = "aws-network-firewall-logs-bucket"
  }
}
```

## Collect

[Collect](https://tailpipe.io/docs/manage/collection) logs for all `aws_network_firewall_log` partitions:

```sh
tailpipe collect aws_network_firewall_log
```

Or for a single partition:

```sh
tailpipe collect aws_network_firewall_log.my_logs
```

## Query

**[Explore example queries for this table →](https://hub.tailpipe.io/plugins/turbot/aws/queries/aws_network_firewall_log)**

### Traffic by Protocol

Analyze traffic by protocol to understand the distribution of network traffic.

```sql
select
  event ->> 'proto' as protocol,
  count(*) as traffic_count
from
  aws_network_firewall_log
group by
  protocol
order by
  traffic_count desc;
```

### Destination Port Analysis

Identify the most commonly accessed destination ports to understand traffic patterns.

```sql
select
  event ->> 'dest_port' as destination_port,
  count(*) as connection_count
from
  aws_network_firewall_log
where
  event ->> 'dest_port' is not null
group by
  destination_port
order by
  connection_count desc
limit 10;
```

### Traffic by Application Protocol

Analyze traffic by application protocol to understand what types of applications are generating network traffic.

```sql
select
  event ->> 'app_proto' as app_protocol,
  count(*) as connection_count
from
  aws_network_firewall_log
group by
  app_protocol
order by
  connection_count desc;
```

## Example Configurations

### Collect logs from an S3 bucket

Collect Network Firewall logs stored in an S3 bucket.

```hcl
connection "aws" "network_firewall_logging" {
  profile = "my-network-firewall-logging"
}

partition "aws_network_firewall_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.network_firewall_logging
    bucket     = "aws-network-firewall-logs-bucket"
  }
}
```

### Collect logs from an S3 bucket with a prefix

Collect Network Firewall logs stored in an S3 bucket using a prefix.

```hcl
partition "aws_network_firewall_log" "my_logs_prefix" {
  source "aws_s3_bucket" {
    connection = connection.aws.network_firewall_logging
    bucket     = "aws-network-firewall-logs-bucket"
    prefix     = "my/prefix/"
  }
}
```

### Collect logs from local files

You can also collect Network Firewall logs from local files.

```hcl
partition "aws_network_firewall_log" "local_logs" {
  source "file"  {
    paths       = ["/Users/myuser/network_firewall_logs"]
    file_layout = "%{DATA}.json.gz"
  }
}
```

### Collect logs for all accounts in an organization

For a specific organization, collect logs for all accounts and regions.

```hcl
partition "aws_network_firewall_log" "my_logs_org" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.network_firewall_logging
    bucket      = "network-firewall-logs-bucket"
    file_layout = "AWSLogs/o-aa111bb222/%{NUMBER:account_id}/network-firewall/(?<log_type>flow|alert|tls)/%{DATA:region}/%{DATA:firewall_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{NUMBER:account_id}_network-firewall_%{DATA:log_type}_%{DATA:region}_%{DATA:firewall_name}_%{DATA:timestamp}_%{DATA:hash}.log.gz"
  }
}
```

### Collect logs for a single account

For a specific account, collect logs for all regions.

```hcl
partition "aws_network_firewall_log" "my_logs_account" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.network_firewall_logging
    bucket      = "network-firewall-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?123456789012/network-firewall/(?<log_type>flow|alert|tls)/%{DATA:region}/%{DATA:firewall_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/123456789012_network-firewall_%{DATA:log_type}_%{DATA:region}_%{DATA:firewall_name}_%{DATA:timestamp}_%{DATA:hash}.log.gz"
  }
}
```

### Collect logs for a single region

For all accounts, collect logs from `us-east-1`.

```hcl
partition "aws_network_firewall_log" "my_logs_region" {
  source "aws_s3_bucket"  {
    connection  = connection.aws.network_firewall_logging
    bucket      = "network-firewall-logs-bucket"
    file_layout = "AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/network-firewall/(?<log_type>flow|alert|tls)/us-east-1/%{DATA:firewall_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{NUMBER:account_id}_network-firewall_%{DATA:log_type}_us-east-1_%{DATA:firewall_name}_%{DATA:timestamp}_%{DATA:hash}.log.gz"
  }
}
```

### Filter for specific traffic type

Use the filter argument in your partition to focus on specific traffic types.

```hcl
partition "aws_network_firewall_log" "my_logs_filtered" {
  filter = "event->>'app_proto' = 'tls'"

  source "aws_s3_bucket" {
    connection  = connection.aws.network_firewall_logging
    bucket      = "network-firewall-logs-bucket"
  }
}
```

## Source Defaults

### aws_s3_bucket

This table sets the following defaults for the [aws_s3_bucket source](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket#arguments):

| Argument      | Default |
|--------------|---------|
| file_layout  | `AWSLogs/(%{DATA:org_id}/)?%{NUMBER:account_id}/network-firewall/(?<log_type>flow|alert|tls)/%{DATA:region}/%{DATA:firewall_name}/%{YEAR:year}/%{MONTHNUM:month}/%{MONTHDAY:day}/%{HOUR:hour}/%{NUMBER:account_id}_network-firewall_%{DATA:log_type}_%{DATA:region}_%{DATA:firewall_name}_%{DATA:timestamp}_%{DATA:hash}.log.gz` |
