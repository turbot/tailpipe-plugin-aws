# AWS Plugin for Tailpipe

Collect and query AWS logs using SQL to track activity, monitor trends, detect anomalies, and more!

## Table of Contents

- [Documentation](#documentation)
- [Getting Started](#getting-started)
- [Credentials](#credentials)
- [Developing](#developing)
- [Open Source & Contributing](#open-source-contributing)
- [Get Involved](#get-involved)

## Documentation

- **[Table configuration and definitions →](https://hub.tailpipe.io/plugins/turbot/aws/tables)**
- **[Table queries →](https://hub.tailpipe.io/plugins/turbot/aws/queries)**
- **[Source definitions →](https://hub.tailpipe.io/plugins/turbot/aws/sources)**

## Getting Started

### Installation

Download and install Tailpipe (https://tailpipe.io/downloads). Or use Brew:

```sh
brew tap turbot/tap
brew install tailpipe
```

Install the plugin:

```sh
tailpipe plugin install aws
```

### Configuration

Configure your log source:

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "aws_profile" {
  profile = "my-profile"
}

partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.aws_profile
    bucket     = "aws-cloudtrail-logs-bucket"
  }
}
```

### Log Collection

Collect logs:

```sh
tailpipe collect aws_cloudtrail_log
```

When running `tailpipe collect` for the first time, logs from the last 7 days are collected. Subsequent `tailpipe collect` runs will collect logs from the last collection date.

You can override the default behaviour by specifying `--from`:

```sh
tailpipe collect aws_cloudtrail_log --from 2025-01-01
```

You can also use relative times. For instance, to collect logs from the last 60 days:

```sh
tailpipe collect aws_cloudtrail_log --from T-60d
```

Please note that if you specify a date in `--from`, Tailpipe will delete any collected data for that partition starting from that date to help avoid gaps in the data.

For additional examples on using `tailpipe collect`, please see [tailpipe collect](https://tailpipe.io/docs/reference/cli/collect) reference documentation.

### Query

Run a query:

```sql
select event_source, event_name, recipient_account_id, count(*) as event_count from aws_cloudtrail_log where not read_only group by event_source, event_name, recipient_account_id order by event_count desc;
```

For example:

```sh
+----------------------+-----------------------+----------------------+-------------+
| event_source         | event_name            | recipient_account_id | event_count |
+----------------------+-----------------------+----------------------+-------------+
| logs.amazonaws.com   | CreateLogStream       | 123456789012         | 793845      |
| ecs.amazonaws.com    | RunTask               | 456789012345         | 350836      |
| ecs.amazonaws.com    | SubmitTaskStateChange | 456789012345         | 190185      |
| s3.amazonaws.com     | PutObject             | 789012345678         | 60842       |
| sns.amazonaws.com    | TagResource           | 456789012345         | 25499       |
| lambda.amazonaws.com | TagResource           | 123456789012         | 20673       |
+----------------------+-----------------------+----------------------+-------------+
```

## Credentials

By default, the following environment variables will be used for authentication:

- `AWS_PROFILE`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

You can also create `connection` resources in configuration files:

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "aws_profile" {
  profile = "my-profile"
}

connection "aws" "aws_access_key_pair" {
  access_key = "AKIA..."
  secret_key = "dP+C+J..."
}

connection "aws" "aws_session_token" {
  access_key    = "AKIA..."
  secret_key    = "dP+C+J..."
  session_token = "AQoDX..."
}
```

For more information on AWS connections in Tailpipe, please see [Managing AWS Connections](https://tailpipe.io/docs/reference/config-files/connection/aws).

## Developing

Prerequisites:

- [Tailpipe](https://tailpipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/tailpipe-plugin-aws.git
cd tailpipe-plugin-aws
```

After making your local changes, build the plugin, which automatically installs the new version to your `~/.tailpipe/plugins` directory:

```
make
```

Try it!

```
tailpipe query
> .inspect aws_cloudtrail_log
```

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Tailpipe](https://tailpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #tailpipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Tailpipe](https://github.com/turbot/tailpipe/labels/help%20wanted)
- [AWS Plugin](https://github.com/turbot/tailpipe-plugin-aws/labels/help%20wanted)
