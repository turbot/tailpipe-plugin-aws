![image](https://hub.tailpipe.io/images/plugins/turbot/aws-social-graphic.png)

# AWS Plugin for Tailpipe

Use SQL to collect and query logs including CloudTrail logs, ELB access logs, S3 server access logs and more from AWS.

- **[Get started →](https://hub.tailpipe.io/plugins/turbot/aws)**
- Documentation: [Table definitions & examples](https://hub.tailpipe.io/plugins/turbot/aws/tables)
- Community: [Join #tailpipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/tailpipe-plugin-aws/issues)

## Quick start

Install the plugin with [Tailpipe](https://tailpipe.io):

```shell
tailpipe plugin install aws
```

Run a query:

```sql
select tp_timestamp, event_source, event_name, user_identity from aws_cloudtrail_log;
```

## Advanced configuration

The AWS plugin has the power to:
* Query multiple accounts
* Query multiple regions
* Use many different methods for credentials (roles, SSO, etc)

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

Build, which automatically installs the new version to your `~/.tailpipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.tailpipe/config
vi ~/.tailpipe/config/aws.tpc
```

Try it!

```
tailpipe query
> .inspect aws
```

Further reading:

- [Writing plugins](https://tailpipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://tailpipe.io/docs/develop/writing-your-first-table)

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Tailpipe](https://tailpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #tailpipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Tailpipe](https://github.com/turbot/tailpipe/labels/help%20wanted)
- [AWS Plugin](https://github.com/turbot/tailpipe-plugin-aws/labels/help%20wanted)
