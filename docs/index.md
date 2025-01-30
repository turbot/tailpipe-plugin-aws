---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/aws.svg"
brand_color: "#FF9900"
display_name: "Amazon Web Services"
description: "Tailpipe plugin for collecting and querying various logs from AWS."
og_description: "Collect AWS logs and query them instantly with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/aws-social-graphic.png"
---

# AWS + Tailpipe

[Tailpipe](https://tailpipe.io) is an open-source CLI tool that allows you to collect logs and query them with SQL.

[AWS](https://aws.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis.

The [AWS Plugin for Tailpipe](https://hub.tailpipe.io/plugins/turbot/aws) allows you to collect and query AWS logs using SQL to track activity, monitor trends, detect anomalies, and more!

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

## Connection Credentials

### Arguments

| Name                   | Type          | Required | Description                                                                                              |
|------------------------|---------------|----------|----------------------------------------------------------------------------------------------------------|
| `access_key`           | String        | No       | AWS access key used for authentication.                                                                 |
| `endpoint_url`         | String        | No       | The custom endpoint URL for AWS services (e.g., for local testing with tools like LocalStack).           |
| `max_error_retry_attempts` | Number    | No       | The maximum number of retry attempts for AWS API calls.                                                 |
| `min_error_retry_delay`    | Number    | No       | The minimum delay in milliseconds between retry attempts for AWS API calls.                             |
| `profile`              | String        | No       | The AWS CLI profile to use for credentials and configuration.                                           |
| `s3_force_path_style`  | Boolean       | No       | Forces the use of path-style URLs for S3 operations instead of the default virtual-hosted style.         |
| `secret_key`           | String        | No       | AWS secret key used for authentication.                                                                 |
| `session_token`        | String        | No       | AWS session token used for temporary credentials. This is only used if you specify `access_key` and `secret_key`. |

### AWS Profile Credentials

You may specify a named profile from an AWS credential file with the `profile` argument. A connection per profile, using named profiles is probably the most common configuration:

#### aws credential file:

```ini
[account_a]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m
region = us-west-2

[account_b]
aws_access_key_id = AKIA4YFAKEKEYJ7HS98F
aws_secret_access_key = Apf938vDKd8ThisIsNotRealzTiEUwXj9nKLWP9mg4
```

#### aws.tpc:

```hcl
connection "aws" "aws_account_a" {
  profile = "account_a"
}

connection "aws" "aws_account_b" {
  profile = "account_b"
}
```

Using named profiles allows Tailpipe to work with your existing CLI configurations, including SSO and using role assumption.

### AWS SSO Credentials

Tailpipe works with [AWS SSO](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html#sso-configure-profile-auto) via AWS profiles however:
- You must login to SSO (`aws sso login`) before starting Tailpipe
- If your credentials expire, you will need to re-authenticate outside of Tailpipe - Tailpipe currently cannot re-authenticate you.

#### aws credential file:

```ini
[account_a_with_sso]
sso_start_url = https://d-9a672b0000.awsapps.com/start
sso_region = us-east-2
sso_account_id = 000000000000
sso_role_name = SSO-ReadOnly
region = us-east-1
```

#### aws.tpc:

```hcl
connection "aws" "aws_account_a_with_sso" {
  profile = "account_a_with_sso"
}
```

### AssumeRole Credentials (No MFA)

If your aws credential file contains profiles that assume a role via the `source_profile` and `role_arn` options and MFA is not required, Tailpipe can use the profile as-is:

#### aws credential file:

```ini
# This user must have sts:AssumeRole permission for arn:aws:iam::*:role.tpc_role
[cli_user]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m

[account_a_role_without_mfa]
role_arn = arn:aws:iam::111111111111:role.tpc_role
source_profile = cli_user
external_id = xxxxx

[account_b_role_without_mfa]
role_arn = arn:aws:iam::222222222222:role.tpc_role
source_profile = cli_user
external_id = yyyyy
```

#### aws.tpc:

```hcl
connection "aws" "aws_account_a" {
  profile = "account_a_role_without_mfa"
}

connection "aws" "aws_account_b" {
  profile = "account_b_role_without_mfa"
}
```

### AssumeRole Credentials (With MFA)

Currently Tailpipe doesn't support prompting for an MFA token at run time. To overcome this problem you will need to generate an AWS profile with temporary credentials.

One way to accomplish this is to use the `credential_process` to [generate the credentials with a script or program](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html) and cache the tokens in a new profile. There is a [sample `mfa.sh` script](https://raw.githubusercontent.com/turbot/tailpipe-plugin-aws/main/scripts/mfa.sh) in the `scripts` directory of the [tailpipe-plugin-aws](https://github.com/turbot/tailpipe-plugin-aws) repo that you can use, and there are several open source projects that automate this process as well.

Note that Tailpipe cannot prompt you for your token currently, so you must authenticate before running `tailpipe collect`, and re-authenticate outside of Tailpipe whenever your credentials expire.

#### aws credential file:

```ini
[cli_user]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m
mfa_serial = arn:aws:iam::999999999999:mfa/my_role_mfa

[account_a_role_with_mfa]
credential_process = sh -c 'mfa.sh arn:aws:iam::111111111111:role/my_role arn:aws:iam::999999999999:mfa/my_role_mfa cli_user 2> $(tty)'

[account_b_role_with_mfa]
credential_process = sh -c 'mfa.sh arn:aws:iam::222222222222:role/my_role arn:aws:iam::999999999999:mfa/my_role_mfa cli_user 2> $(tty)'
```

#### aws.tpc:

```hcl
connection "aws" "aws_account_a" {
  profile = "account_a_role_with_mfa"
}

connection "aws" "aws_account_b" {
  profile = "account_b_role_with_mfa"
}
```

### AssumeRole Credentials (in ECS)

If you are using Tailpipe on AWS ECS then you need to ensure that have separated your [Task Role](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html) and [Execution Role](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html) within the Task Definition. You will also need to create a separate service role that your `Task Role` can assume.

The Task Role should have permissions to assume your service role. Additionally your service role needs a trust relationship set up, and have permissions to assume your other roles.

#### Task Role IAM Assume Role

```json
{
    "Version": "2012-10-17"
    "Statement": [
        {
            "Action": [
                "sts:AssumeRole"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:aws:iam::111111111111:role/tailpipe-service"
            ]
        }
    ]
}
```

#### Service Role

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::111111111111:role/tailpipe-ecs-task-role"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```

This will allow you to configure Tailpipe now to assume the service role.

#### aws credential file:

```ini
[default]
role_arn = arn:aws:iam::111111111111:role/tailpipe-service
credential_source = EcsContainer

[account_b]
role_arn = arn:aws:iam::222222222222:role/tailpipe_ro_role
source_profile = default
```

### AWS-Vault Credentials

Tailpipe can use profiles that use [aws-vault](https://github.com/99designs/aws-vault) via the `credential_process`. aws-vault can even be used when using AssumeRole Credentials with MFA (you must authenticate/re-authenticate outside of Tailpipe whenever your credentials expire if you are using MFA).

When authenticating with temporary credentials, like using an access key pair with aws-vault, some IAM and STS APIs may be restricted. You can avoid creating a temporary session with the `--no-session` option (e.g., `aws-vault exec my_profile --no-session -- tailpipe collect aws_cloudtrail_log"`). For more information, please see [aws-vault Temporary credentials limitations with STS, IAM
](https://github.com/99designs/aws-vault/blob/master/USAGE.md#temporary-credentials-limitations-with-sts-iam).

#### aws credential file:

```ini
[vault_user_account]
credential_process = /usr/local/bin/aws-vault exec -j vault_user_profile # vault_user_profile is the name of the profile in AWS_VAULT...

[account_a]
source_profile = vault_user_account
role_arn = arn:aws:iam::123456789012:role/my_role
mfa_serial = arn:aws:iam::123456789012:mfa/my_role_mfa
```

#### aws.tpc:

```hcl
connection "aws" "aws_account_a" {
  profile = "account_a"
}
```

### IAM Access Key Pair Credentials

The AWS plugin allows you set static credentials with the `access_key`, `secret_key`, and `session_token` arguments in your connection.

```hcl
connection "aws" "aws_account_a" {
  secret_key = "gMCYsoGqjfThisISNotARealKeyVVhh"
  access_key = "ASIA3ODZSWFYSN2PFHPJ"
}
```

### Credentials from Environment Variables

The AWS plugin will use the standard AWS environment variables to obtain credentials **only if other arguments (`profile`, `access_key`/`secret_key`) are not specified** in the connection:

```sh
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=eu-west-1
export AWS_SESSION_TOKEN=AQoDYXdzEJr...
export AWS_ROLE_SESSION_NAME=tailpipe@myaccount
```

### Credentials from an EC2 Instance Profile

If you are running Tailpipe on a AWS EC2 instance, and that instance has an [instance profile attached](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) then Tailpipe will automatically use the associated IAM role without other credentials.
