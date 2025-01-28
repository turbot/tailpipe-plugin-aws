# AWS Plugin for Tailpipe

[Tailpipe](https://tailpipe.io) is an open-source CLI tool that allows you to collect logs and query them with SQL.

[AWS](https://aws.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis.

The [AWS Plugin for Tailpipe](https://hub.tailpipe.io/plugins/turbot/aws) allows you to collect and query AWS logs using SQL to track activity, monitor trends, detect anomalies, and more!

- **[Get started →](https://hub.tailpipe.io/plugins/turbot/aws)**
- Documentation: [Table definitions & examples](https://hub.tailpipe.io/plugins/turbot/aws/tables)
- Community: [Join #tailpipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/tailpipe-plugin-aws/issues)

Collect and query logs:
![image](docs/images/aws_cloudtrail_log_terminal.png)

## Getting Started

Install Tailpipe from the [downloads](https://tailpipe.io/downloads) page:

```sh
# MacOS
brew install turbot/tap/tailpipe
```

```sh
# Linux or Windows (WSL)
sudo /bin/sh -c "$(curl -fsSL https://tailpipe.io/install/tailpipe.sh)"
```

Install the plugin:

```sh
tailpipe plugin install aws
```

Configure your [connection credentials](https://hub.tailpipe.io/plugins/turbot/aws#connection-credentials), table partition, and data source ([examples](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cloudtrail_log#example-configurations)):

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "logging_account" {
  profile = "my-logging-account"
}

partition "aws_cloudtrail_log" "my_logs" {
  source "aws_s3_bucket" {
    connection = connection.aws.logging_account
    bucket     = "aws-cloudtrail-logs-bucket"
  }
}
```

Download, enrich, and save logs from your source ([examples](https://tailpipe.io/docs/reference/cli/collect)):

```sh
tailpipe collect aws_cloudtrail_log
```

Enter interactive query mode:

```sh
tailpipe query
```

Run a query:

```sql
select
  event_source,
  event_name,
  count(*) as event_count
from
  aws_cloudtrail_log
where
  not read_only
group by
  event_source,
  event_name
order by
  event_count desc;
```

```sh
+----------------------+-----------------------+-------------+
| event_source         | event_name            | event_count |
+----------------------+-----------------------+-------------+
| logs.amazonaws.com   | CreateLogStream       | 793845      |
| ecs.amazonaws.com    | RunTask               | 350836      |
| ecs.amazonaws.com    | SubmitTaskStateChange | 190185      |
| s3.amazonaws.com     | PutObject             | 60842       |
| sns.amazonaws.com    | TagResource           | 25499       |
| lambda.amazonaws.com | TagResource           | 20673       |
+----------------------+-----------------------+-------------+
```

## Detections as Code with Powerpipe

Pre-built dashboards and detections for the AWS plugin are available in [Powerpipe](https://powerpipe.io) mods, helping you monitor and analyze activity across your AWS accounts.

For example, the [AWS CloudTrail Logs Detections mod](https://hub.powerpipe.io/mods/turbot/tailpipe-mod-aws-cloudtrail-log-detections) scans your CloudTrail logs for anomalies, such as an S3 bucket being made public or a change in your VPC network infrastructure.

Dashboards and detections are [open source](https://github.com/topics/tailpipe-mod), allowing easy customization and collaboration.

To get started, choose a mod from the [Powerpipe Hub](https://hub.powerpipe.io/?engines=tailpipe&q=aws).

![image](docs/images/aws_cloudtrail_log_mitre_dashboard.png)

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

Re-collect your data:

```
tailpipe collect aws_cloudtrail_log
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
