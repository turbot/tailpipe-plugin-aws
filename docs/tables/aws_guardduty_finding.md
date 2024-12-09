---
title: "Tailpipe Table: aws_guardduty_finding - Query AWS GuardDuty Findings"
description: "Allows users to query AWS GuardDuty findings."
---

# Table: aws_guardduty_finding - Query AWS GuardDuty Findings

*TODO*: Add description

## Table Usage Guide

The `aws_guardduty_finding` table allows you to query data from AWS GuardDuty findings. This table provides detailed information about security findings detected by AWS GuardDuty, including the finding type, severity, resource details, and more.

## Examples

### Detect High Severity Findings
Flags high severity findings detected by AWS GuardDuty, which may indicate potential security threats.

```sql
select
  resource.resource_type,
  resource.resource_id,
  severity,
  type,
  title,
  description,
  service_action.action_type,
  service_action.action_target,
  service_action.evidence,
  resource_details.instance_details.instance_id,
  resource_details.instance_details.instance_type,
  resource_details.instance_details.launch_time,
  resource_details.instance_details.availability_zone,
  resource_details.instance_details.image_id,
  resource_details.instance_details.image_description,
  resource_details.instance_details.network_interfaces,
  resource_details.instance_details.tags,
  resource_details.resource_role,
  resource_details.resource_type,
  resource_details.resource_name,
  resource_details.resource_creation_time,
  resource_details.resource_deletion_time,
  resource_details.resource_details,
  resource_details.resource_status,
  resource_details.resource_status_reason,
  resource_details.resource_status_reason_code,
  resource_details.resource_status_reason_message,
  resource_details.resource_status_reason_details,
  resource_details.resource_status_reason_details_message,
  resource_details.resource_status_reason_details_code,
  resource_details.resource_status_reason_details_status
from
  aws_guardduty_finding
where
  severity >= 7
order by
  severity desc;
```