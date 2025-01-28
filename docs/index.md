---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/aws.svg"
brand_color: "#FF9900"
display_name: "Amazon Web Services"
description: "Tailpipe plugin for collecting and querying various logs from AWS."
og_description: "Collect and query AWS logs with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/aws-social-graphic.png"
---

# AWS + Tailpipe

[Tailpipe](https://tailpipe.io) is an open-source CLI tool that allows you to collect logs and query them with SQL.

[AWS](https://aws.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis.

The [Tailpipe AWS plugin](https://hub.tailpipe.io/plugins/turbot/aws) for Tailpipe allows you to collect and query AWS logs using SQL to track activity, monitor trends, detect anomalies, and more!

- Documentation: [Table definitions & examples](https://hub.tailpipe.io/plugins/turbot/aws/tables)
- Community: [Join #tailpipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/tailpipe-plugin-aws/issues)

<img src="https://raw.githubusercontent.com/turbot/tailpipe-plugin-aws/main/docs/images/aws_cloudtrail_log_terminal.png" width="50%" type="thumbnail"/>
<img src="https://raw.githubusercontent.com/turbot/tailpipe-plugin-aws/main/docs/images/aws_cloudtrail_log_mitre_dashboard.png" width="50%" type="thumbnail"/>

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

Configure your [connection credentials](https://tailpipe.io/docs/reference/config-files/connection/aws), [table partition](https://tailpipe.io/docs/manage/partition), and [data source](https://tailpipe.io/docs/manage/source):

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

Download, enrich, and save logs from your S3 bucket:

```sh
tailpipe collect aws_cloudtrail_log
```

> [!NOTE]
> When running `tailpipe collect` for the first time, logs from the last 7 days are collected. Subsequent `tailpipe collect` runs will collect logs from the last collection date.
>
> For more information, please see [Managing Collection](https://tailpipe.io/docs/manage/collection).

Enter interactive query mode:

```sh
tailpipe query
```

Run a query:

```sql
select event_source, event_name, count(*) as event_count
from aws_cloudtrail_log
where not read_only
group by event_source, event_name
order by event_count desc;
```

For example:

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

Dashboards and detections in this mod and others are written as code, making them easy to customize and adapt to your specific requirements.

To get started, browse the [Powerpipe Mods for the AWS plugin](https://hub.tailpipe.io/plugins/turbot/aws/mods) and follow the instructions provided for each mod.
