# AWS Plugin for Tailpipe

Collect and query AWS logs using SQL to track activity, monitor trends, and detect anomalies.

- **[Get started →](https://hub.tailpipe.io/plugins/turbot/aws)**
- Documentation: [Table queries](https://hub.tailpipe.io/plugins/turbot/aws/queries)
- Community: [Join #tailpipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/tailpipe-plugin-aws/issues)

## Quick Start

Install the plugin with [Tailpipe](https://tailpipe.io):

```sh
tailpipe plugin install aws
```

Configure your log source:

```sh
vi ~/.tailpipe/config/aws.tpc
```

```terraform
connection "aws" "dev" {
  profile = "dev"
}

partition "aws_cloudtrail_log" "dev" {
  source "aws_s3_bucket" {
    connection = connection.aws.dev
    bucket     = "aws-cloudtrail-logs-bucket"
  }
}
```

Collect logs:

```sh
tailpipe collect aws_cloudtrail_log.dev
```

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

## Advanced Configuration

The AWS plugin has the power to:
* Collect logs from various sources, including AWS CloudWatch log groups, S3 buckets, and more
* Use many different methods for credentials (roles, SSO, etc.)

- **[Detailed configuration guide →](https://hub.tailpipe.io/plugins/turbot/aws#get-started)**

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
> .inspect
```

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Tailpipe](https://tailpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #tailpipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Tailpipe](https://github.com/turbot/tailpipe/labels/help%20wanted)
- [AWS Plugin](https://github.com/turbot/tailpipe-plugin-aws/labels/help%20wanted)
