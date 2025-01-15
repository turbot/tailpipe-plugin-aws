---
organization: Turbot
category: ["public cloud"]
icon_url: "/images/plugins/turbot/aws.svg"
brand_color: "#FF9900"
display_name: "Amazon Web Services"
short_name: "aws"
description: "Tailpipe plugin for collecting and querying various logs from AWS."
og_description: "Query AWS logs with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/aws-social-graphic.png"
engines: ["tailpipe"]
---

# AWS + Tailpipe

[Tailpipe](https://tailpipe.io) is an open-source CLI tool that allows you to collect logs and query them with SQL.

[AWS](https://aws.amazon.com/) provides on-demand cloud computing platforms and APIs to authenticated customers on a metered pay-as-you-go basis.

For example:

```sql
select
  event_time,
  event_source,
  event_name,
  user_identity.arn
from
  aws_cloudtrail_log;
```

```sh
+---------------------+-------------------+------------+--------------------------------------------------------------+
| event_time          | event_source      | event_name | arn                                                          |
+---------------------+-------------------+------------+--------------------------------------------------------------+
| 2024-09-30 23:56:11 | iam.amazonaws.com | CreateUser | arn:aws:sts::123456789012:assumed-role/admin_role/pam        |
| 2024-09-30 23:56:11 | iam.amazonaws.com | CreateRole | arn:aws:sts::123456789012:role/warehouse                     |
| 2024-09-30 23:56:13 | ec2.amazonaws.com | CopyImage  | arn:aws:sts::123456789012:assumed-role/assistant_role/dwight |
| 2024-09-30 23:56:13 | sts.amazonaws.com | AssumeRole | arn:aws:sts::123456789012:assumed-role/qa_role/creed         |
| 2024-09-30 23:57:14 | sts.amazonaws.com | AssumeRole | arn:aws:sts::123456789012:assumed-role/qa_role/creed         |
+---------------------+-------------------+------------+--------------------------------------------------------------+
```

## Documentation

- **[Table definitions →](/plugins/turbot/aws/tables)**
- **[Table queries →](/plugins/turbot/aws/queries)**
- **[Source definitions →](/plugins/turbot/aws/sources)**

## Get Started

Install the plugin with [Tailpipe](https://tailpipe.io):

```shell
tailpipe plugin install aws
```

Configure your log source:

```shell
vi ~/.tailpipe/config/aws.tpc
```

```terraform
connection "aws" "dev" {
  profile = "dev"
}

partition "aws_cloudtrail_log" "dev" {
  source "aws_s3_bucket" {
    bucket = "aws-cloudtrail-logs-bucket"
  }
}
```

Collect logs:

```shell
tailpipe collect aws_cloudtrail_log.dev
```

Run a query:

```sql
select event_source, event_name, recipient_account_id, count(*) as event_count from aws_cloudtrail_log where not read_only group by event_source, event_name, recipient_account_id order by event_count desc;
```

## Credentials

| Item | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| - |----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Credentials | Specify a named profile from an AWS credential file with the `profile` argument.                                                                                                                                                                                                                                                                                                                                                                               |
| Permissions | Grant the `ReadOnlyAccess` policy to your user or role.                                                                                                                                                                                                                                                                                                                                                                                                        |
| Radius | Each connection represents a single AWS account.                                                                                                                                                                                                                                                                                                                                                                                                               |
| Resolution | 1. Credentials explicitly set in a Tailpipe config file (`~/.tailpipe/config/aws.tpc`).<br />2. Credentials specified in environment variables, e.g., `AWS_ACCESS_KEY_ID`.<br />3. Credentials in the credential file (`~/.aws/credentials`) for the profile specified in the `AWS_PROFILE` environment variable.<br />4. Credentials for the default profile from the credential file.<br />5. EC2 instance role credentials (if running on an EC2 instance). |

### Profiles

You may specify a named profile from an AWS credential file with the `profile` argument. A connection per profile, using named profiles is probably the most common configuration:

```sh
vi ~/.aws/credentials
```

```ini
[account_a]
aws_access_key_id     = AKIA4YFAKEKEYXTDS...
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5...
region                = us-west-2

[account_b]
aws_access_key_id     = AKIA4YFAKEKEYJ7H...
aws_secret_access_key = Apf938vDKd8ThisIsNotRealzTiEUwXj9nKLWP9...
```

```sh
vi ~/.tailpipe/config/aws.tpc
```

```hcl
connection "aws" "account_a" {
  profile = "account_a"
}

connection "aws" "account_b" {
  profile = "account_b"
}
```

Using named profiles allows Tailpipe to work with your existing CLI configurations, including SSO and using role assumption.

### AWS SSO Credentials

Tailpipe works with [AWS SSO](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html#sso-configure-profile-auto) via AWS profiles however:
- You must login to SSO (`aws sso login`) before starting Tailpipe.
- If your credentials expire, you will need to re-authenticate outside of Tailpipe - Tailpipe currently cannot re-authenticate you.

#### aws credential file:
```ini
[account_a_with_sso]
sso_start_url  = https://d-9a672b0000.awsapps.com/start
sso_region     = us-east-2
sso_account_id = 000000000000
sso_role_name  = SSO-ReadOnly
region         = us-east-1
```

#### aws.tpc:
```hcl
connection "aws" "account_a_with_sso" {
  profile = "account_a_with_sso"
}
```

### AssumeRole Credentials (No MFA)

If your aws credential file contains profiles that assume a role via the `source_profile` and `role_arn` options and MFA is not required, Tailpipe can use the profile as-is:

#### aws credential file:

```ini
# This user must have sts:AssumeRole permission for arn:aws:iam::*:role/tpc_role
[cli_user]
aws_access_key_id = AKIA4YFAKEKEYXTDS252
aws_secret_access_key = SH42YMW5p3EThisIsNotRealzTiEUwXN8BOIOF5J8m

[account_a_role_without_mfa]
role_arn = arn:aws:iam::111111111111:role/tpc_role
source_profile = cli_user
external_id = xxxxx
region = us-east-1

[account_b_role_without_mfa]
role_arn = arn:aws:iam::222222222222:role/tpc_role
source_profile = cli_user
external_id = yyyyy
region = us-east-2
```

#### aws.tpc:

```hcl
connection "aws" "account_a" {
  profile = "account_a_role_without_mfa"
}

connection "aws" "account_b" {
  profile = "account_b_role_without_mfa"
}
```

### AssumeRole Credentials (With MFA)

Currently, Tailpipe doesn't support prompting for an MFA token at run time. To overcome this problem you will need to generate an AWS profile with temporary credentials.

One way to accomplish this is to use the `credential_process` to [generate the credentials with a script or program](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html) and cache the tokens in a new profile. There is a [sample `mfa.sh` script](https://raw.githubusercontent.com/turbot/steampipe-plugin-aws/main/scripts/mfa.sh) in the `scripts` directory of the [steampipe-plugin-aws](https://github.com/turbot/steampipe-plugin-aws) repo that you can use, and there are several open source projects that automate this process as well.

Note that Tailpipe cannot prompt you for your token currently, so you must authenticate before starting Tailpipe, and re-authenticate outside of Tailpipe whenever your credentials expire.

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
connection "aws" "account_a" {
  profile        = "account_a_role_with_mfa"
  default_region = "us-east-1"
}

connection "aws" "account_b" {
  profile = "account_b_role_with_mfa"
  regions = ["us-east-2"]
}
```

### AssumeRole Credentials (in ECS)

If you are using Tailpipe on AWS ECS then you need to ensure that have separated your [Task Role](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html) and [Execution Role](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html) within the Task Definition. You will also need to create a separate service role that your `Task Role` can assume.

The Task Role should have permissions to assume your service role. Additionally, your service role needs a trust relationship set up, and have permissions to assume your other roles.

#### Task Role IAM Assume Role
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "sts:AssumeRole"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:aws:iam::111111111111:role/steampipe-service"
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
                "AWS": "arn:aws:iam::111111111111:role/steampipe-ecs-task-role"
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

#### aws.tpc:
```hcl
connection "aws" "account_b" {
  profile = "account_b"
}
```

### AWS-Vault Credentials

Tailpipe can use profiles that use [aws-vault](https://github.com/99designs/aws-vault) via the `credential_process`. aws-vault can even be used when using AssumeRole Credentials with MFA (you must authenticate/re-authenticate outside of Tailpipe whenever your credentials expire if you are using MFA).

When authenticating with temporary credentials, like using an access key pair with aws-vault, some IAM and STS APIs may be restricted. You can avoid creating a temporary session with the `--no-session` option (e.g., `aws-vault exec my_profile --no-session -- tailpipe collect aws_cloudtrail_log.account_a"`). For more information, please see [aws-vault Temporary credentials limitations with STS, IAM](https://github.com/99designs/aws-vault/blob/master/USAGE.md#temporary-credentials-limitations-with-sts-iam).

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
connection "aws" "account_a" {
  profile = "account_a"
  regions = ["*"]
}
```

### IAM Access Key Pair Credentials

The AWS plugin allows you set static credentials with the `access_key`, `secret_key`, and `session_token` arguments in your connection.

```hcl
connection "aws_account_a" {
  access_key     = "ASIA3ODZSWFYSN2PFHPJ"
  secret_key     = "gMCYsoGqjfThisISNotARealKeyVVhh"
  session_token  = "FwoGZXIvYXdzEJv//////////wEaDINJ"
  default_region = "us-east-1"
}
```

### Credentials from Environment Variables

The AWS plugin will use the standard AWS environment variables to obtain credentials **only if other arguments (`profile`, `access_key`/`secret_key`, `regions`) are not specified / no connection is passed**:

```sh
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=eu-west-1
export AWS_SESSION_TOKEN=AQoDYXdzEJr...
export AWS_ROLE_SESSION_NAME=steampipe@myaccount
```

### Credentials from an EC2 Instance Profile

If you are running Tailpipe on a AWS EC2 instance, and that instance has an [instance profile attached](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) then Tailpipe will automatically use the associated IAM role without the need for making or passing a connection.

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Tailpipe](https://tailpipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #tailpipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Tailpipe](https://github.com/turbot/tailpipe/labels/help%20wanted)
- [AWS Plugin](https://github.com/turbot/tailpipe-plugin-aws/labels/help%20wanted)
