## v0.3.0 [2025-02-12]

_What's new?_

- New tables added
  - [aws_s3_server_access_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_s3_server_access_log) ([#75](https://github.com/turbot/tailpipe-plugin-aws/pull/75))

_Enhancements_

- Added `Type` column in `aws_s3_bucket` source arguments table.

_Dependencies_

- Bumped github.com/aws/aws-sdk-go-v2/config from 1.28.11 to 1.29.6. ([#79](https://github.com/turbot/tailpipe-plugin-aws/pull/79))
- Bumped github.com/aws/aws-sdk-go-v2/credentials from 1.17.52 to 1.17.57. ([#68](https://github.com/turbot/tailpipe-plugin-aws/pull/68))
- Bumped github.com/aws/aws-sdk-go-v2/feature/s3/manager from 1.17.49 to 1.17.60. ([#81](https://github.com/turbot/tailpipe-plugin-aws/pull/81))
- Bumped github.com/aws/aws-sdk-go-v2/service/s3 from 1.72.3 to 1.76.0. ([#78](https://github.com/turbot/tailpipe-plugin-aws/pull/78))
- Bumped github.com/hashicorp/hcl/v2 from 2.20.1 to 2.23.0. ([#84](https://github.com/turbot/tailpipe-plugin-aws/pull/84))
- Bumped github.com/rs/xid from 1.5.0 to 1.6.0. ([#67](https://github.com/turbot/tailpipe-plugin-aws/pull/67))
- Bumped github.com/turbot/tailpipe-plugin-sdk from 0.1.0 to 0.1.1. ([#75](https://github.com/turbot/tailpipe-plugin-aws/pull/75))
- Bumped golang.org/x/sync from 0.10.0 to 0.11.0. ([#82](https://github.com/turbot/tailpipe-plugin-aws/pull/82))

## v0.2.0 [2025-02-06]

_Enhancements_

- Updated documentation formatting and enhanced argument descriptions for `aws_s3_bucket` source. ([#76](https://github.com/turbot/tailpipe-plugin-aws/pull/76))

## v0.1.0 [2025-01-30]

_What's new?_

- New tables added
  - [aws_cloudtrail_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cloudtrail_log)
- New sources added
  - [aws_s3_bucket](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket)
