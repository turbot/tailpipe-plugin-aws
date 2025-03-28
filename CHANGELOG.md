## v0.8.0 [2025-03-28]

_What's new?_

- New tables added:
  - [aws_cost_and_usage_focus](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_and_usage_focus) ([#115](https://github.com/turbot/tailpipe-plugin-aws/pull/115))
  - [aws_cost_and_usage_report](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_and_usage_report) ([#115](https://github.com/turbot/tailpipe-plugin-aws/pull/115))
  - [aws_cost_optimization_recommendation](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cost_optimization_recommendation) ([#115](https://github.com/turbot/tailpipe-plugin-aws/pull/115))
  - [aws_guardduty_finding](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_guardduty_finding) ([#130](https://github.com/turbot/tailpipe-plugin-aws/pull/130))

_Dependencies_

- Bumped github.com/turbot/go-kit from 1.1.0 to 1.2.0. ([#128](https://github.com/turbot/tailpipe-plugin-aws/pull/128))

## v0.7.0 [2025-03-21]

_What's new?_

- New tables added:
  - [aws_clb_access_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_clb_access_log) ([#88](https://github.com/turbot/tailpipe-plugin-aws/pull/88))

_Dependencies_

- Bumped github.com/aws/aws-sdk-go-v2/config from 1.29.6 to 1.29.9. ([#120](https://github.com/turbot/tailpipe-plugin-aws/pull/120))
- Bumped github.com/aws/aws-sdk-go-v2/feature/s3/manager. ([#124](https://github.com/turbot/tailpipe-plugin-aws/pull/124))
- Bumped github.com/aws/aws-sdk-go-v2/service/s3 from 1.77.1 to 1.78.2. ([#125](https://github.com/turbot/tailpipe-plugin-aws/pull/125))
- Bumped github.com/containerd/containerd from 1.7.18 to 1.7.27. ([#126](https://github.com/turbot/tailpipe-plugin-aws/pull/126))
- Bumped golang.org/x/net from 0.33.0 to 0.36.0. ([#122](https://github.com/turbot/tailpipe-plugin-aws/pull/122))
- Bumped golang.org/x/sync from 0.11.0 to 0.12.0. ([#117](https://github.com/turbot/tailpipe-plugin-aws/pull/117))

## v0.6.0 [2025-03-07]

_What's new?_

- New tables added:
  - [aws_alb_access_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_alb_access_log) ([#116](https://github.com/turbot/tailpipe-plugin-aws/pull/116))
  - [aws_nlb_access_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_nlb_access_log) ([#116](https://github.com/turbot/tailpipe-plugin-aws/pull/116))
  - [aws_vpc_flow_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_vpc_flow_log) ([#116](https://github.com/turbot/tailpipe-plugin-aws/pull/116))
  - [aws_waf_traffic_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_waf_traffic_log) ([#116](https://github.com/turbot/tailpipe-plugin-aws/pull/116))

_Dependencies_

- Bumped github.com/aws/aws-sdk-go-v2 from 1.36.2 to 1.36.3. ([#110](https://github.com/turbot/tailpipe-plugin-aws/pull/110))
- Bumped github.com/aws/aws-sdk-go-v2/credentials from 1.17.59 to 1.17.61. ([#111](https://github.com/turbot/tailpipe-plugin-aws/pull/111))
- Bumped github.com/turbot/pipe-fittings/v2 from 2.1.1 to 2.2.0. ([#100](https://github.com/turbot/tailpipe-plugin-aws/pull/100))

## v0.5.0 [2025-03-03]

_Enhancements_

- Standardized all example query titles to use `Title Case` for consistency. ([#109](https://github.com/turbot/tailpipe-plugin-aws/pull/109))
- Added `folder` front matter to all queries for improved organization and discoverability in the Hub. ([#109](https://github.com/turbot/tailpipe-plugin-aws/pull/109))

_Bug fixes_

- Fixed the `display_name` in `docs/index.md` from `Amazon Web Services` to `AWS` for consistency with standard naming conventions. ([#109](https://github.com/turbot/tailpipe-plugin-aws/pull/109))

## v0.4.0 [2025-02-14]

_Breaking changes_

- The `aws_s3_server_access_log` table index is now based off of the source bucket's name instead of the destination bucket's AWS account ID. We recommend deleting your existing `aws_s3_server_access_log` partition data, e.g., `tailpipe partition delete aws_s3_server_access_log.my_partition`, and then recollecting your data. ([#89](https://github.com/turbot/tailpipe-plugin-aws/pull/89))

## v0.3.0 [2025-02-12]

_What's new?_

- New tables added:
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

- New tables added:
  - [aws_cloudtrail_log](https://hub.tailpipe.io/plugins/turbot/aws/tables/aws_cloudtrail_log)
- New sources added:
  - [aws_s3_bucket](https://hub.tailpipe.io/plugins/turbot/aws/sources/aws_s3_bucket)
