package mappers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// GuardDutyMapper is a mapper that receives GuardDutyBatch objects and extracts GuardDutyFinding records from them
type GuardDutyMapper struct {
}

// NewGuardDutyMapper creates a new GuardDutyMapper
func NewGuardDutyMapper() table.Mapper[rows.GuardDutyFinding] {
	return &GuardDutyMapper{}
}

func (g *GuardDutyMapper) Identifier() string {
	return "guardduty_mapper"
}

// Map casts the data item as a GuardDutyBatch and returns the GuardDutyFinding records
func (g *GuardDutyMapper) Map(_ context.Context, a any) ([]rows.GuardDutyFinding, error) {
	var jsonBytes []byte
	// The expected input type is a JSON byte[] deserializable to GuardDutyBatch
	switch v := a.(type) {

	case []byte:
		jsonBytes = v
	case string:
		jsonBytes = []byte(v)
	default:
		return nil, fmt.Errorf("expected byte[] or string, got %T", a)
	}
	// Decode JSON into a slice of AWS SDK `types.Finding`
	var finding types.Finding
	err := json.Unmarshal(jsonBytes, &finding)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON to findings: %w", err)
	}

	// Populate `GuardDutyFinding` instance with values from `types.Finding`
	row := rows.GuardDutyFinding{
		AccountId:     finding.AccountId,
		Arn:           finding.Arn,
		Description:   finding.Description,
		Id:            finding.Id,
		Partition:     finding.Partition,
		Region:        finding.Region,
		ResourceType:  finding.Resource.ResourceType,
		SchemaVersion: finding.SchemaVersion,
		Archived:      finding.Service.Archived,
		Count:         finding.Service.Count,
		DetectorId:    finding.Service.DetectorId,
		ResourceRole:  finding.Service.ResourceRole,
		ServiceName:   finding.Service.ServiceName,
		Severity:      finding.Severity,
		Title:         finding.Title,
		Type:          finding.Type,
		UpdatedAt:     finding.UpdatedAt,
	}
	if finding.Service.Action != nil {
		if finding.Service.Action.ActionType != nil {
			row.ActionType = finding.Service.Action.ActionType
		}
		if finding.Service.Action.AwsApiCallAction != nil {
			if finding.Service.Action.AwsApiCallAction.Api != nil {
				row.Api = finding.Service.Action.AwsApiCallAction.Api
			}
			if finding.Service.Action.AwsApiCallAction.CallerType != nil {
				row.CallerType = finding.Service.Action.AwsApiCallAction.CallerType
			}
			if finding.Service.Action.AwsApiCallAction.ErrorCode != nil {
				row.ErrorCode = finding.Service.Action.AwsApiCallAction.ErrorCode
			}
			if finding.Service.Action.AwsApiCallAction.RemoteIpDetails != nil {
				if finding.Service.Action.AwsApiCallAction.RemoteIpDetails.IpAddressV4 != nil {
					row.IpAddressV4 = finding.Service.Action.AwsApiCallAction.RemoteIpDetails.IpAddressV4
				}
				if finding.Service.Action.AwsApiCallAction.RemoteIpDetails.IpAddressV6 != nil {
					row.IpAddressV6 = finding.Service.Action.AwsApiCallAction.RemoteIpDetails.IpAddressV6
				}
			}
		}
	}

	// Parse `CreatedAt` if it's provided
	if finding.CreatedAt != nil {
		createdAt, err := time.Parse(time.RFC3339, *finding.CreatedAt)
		if err == nil {
			row.CreatedAt = createdAt
		} else {
			slog.Warn("GuardDutyMapper", "error parsing CreatedAt", err)
		}
	}

	// Parse `EventFirstSeen` and `EventLastSeen`
	if finding.Service.EventFirstSeen != nil {
		row.EventFirstSeen = finding.Service.EventFirstSeen
	}
	if finding.Service.EventLastSeen != nil {
		row.EventLastSeen = finding.Service.EventLastSeen
	}

	// Extract other nested fields if available
	if finding.Resource.AccessKeyDetails != nil {
		row.AccessKeyId = finding.Resource.AccessKeyDetails.AccessKeyId
		row.PrincipalId = finding.Resource.AccessKeyDetails.PrincipalId
		row.UserName = finding.Resource.AccessKeyDetails.UserName
		row.UserType = finding.Resource.AccessKeyDetails.UserType
	}

	if finding.Resource.InstanceDetails != nil {
		row.AvailabilityZone = finding.Resource.InstanceDetails.AvailabilityZone
		row.InstanceArn = finding.Resource.InstanceDetails.IamInstanceProfile.Arn
		row.ImageDescription = finding.Resource.InstanceDetails.ImageDescription
		row.ImageId = finding.Resource.InstanceDetails.ImageId
		row.InstanceId = finding.Resource.InstanceDetails.InstanceId
		row.InstanceState = finding.Resource.InstanceDetails.InstanceState
		row.InstanceType = finding.Resource.InstanceDetails.InstanceType
		row.OutpostArn = finding.Resource.InstanceDetails.OutpostArn
		row.LaunchTime = finding.Resource.InstanceDetails.LaunchTime

		// Network details
		if len(finding.Resource.InstanceDetails.NetworkInterfaces) > 0 {
			ni := finding.Resource.InstanceDetails.NetworkInterfaces[0]
			row.NetworkInterfaceId = *ni.NetworkInterfaceId
			row.PrivateDnsName = *ni.PrivateDnsName
			row.PrivateIpAddress = *ni.PrivateIpAddress
			row.PublicDnsName = *ni.PublicDnsName
			row.PublicIp = *ni.PublicIp
			row.VpcId = *ni.VpcId
			row.SubnetId = *ni.SubnetId

			if len(ni.SecurityGroups) > 0 {
				row.GroupId = *ni.SecurityGroups[0].GroupId
				row.GroupName = *ni.SecurityGroups[0].GroupName
			}
		}
	}
	return []rows.GuardDutyFinding{row}, nil
}
