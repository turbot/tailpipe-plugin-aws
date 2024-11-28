package mappers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"

	"github.com/turbot/tailpipe-plugin-aws/rows"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

// GuardDutyMapper is a mapper that receives GuardDutyBatch objects and extracts GuardDutyFinding records from them
type GuardDutyMapper struct {
}

// NewGuardDutyMapper creates a new GuardDutyMapper
func NewGuardDutyMapper() table.Mapper[*rows.GuardDutyFinding] {
	return &GuardDutyMapper{}
}

func (g *GuardDutyMapper) Identifier() string {
	return "guardduty_mapper"
}

// Map casts the data item as a GuardDutyBatch and returns the GuardDutyFinding records
func (g *GuardDutyMapper) Map(_ context.Context, a any) (*rows.GuardDutyFinding, error) {
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
	row := &rows.GuardDutyFinding{
		AccountId:     finding.AccountId,
		Arn:           finding.Arn,
		Description:   finding.Description,
		Id:            finding.Id,
		Partition:     finding.Partition,
		Region:        finding.Region,
		SchemaVersion: finding.SchemaVersion,
		Severity:      finding.Severity,
		Title:         finding.Title,
		Type:          finding.Type,
	}

	if finding.CreatedAt != nil {
		row.CreatedAt, err = time.Parse(time.RFC3339, *finding.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing CreatedAt: %w", err)
		}
	}

	if finding.UpdatedAt != nil {
		var updatedAt time.Time
		updatedAt, err = time.Parse(time.RFC3339, *finding.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing UpdatedAt: %w", err)
		}
		row.UpdatedAt = &updatedAt
	}

	// service
	if finding.Service != nil {
		row.Service = &rows.GuardDutyFindingService{
			Action:               finding.Service.Action,
			Archived:             finding.Service.Archived,
			Count:                finding.Service.Count,
			Detection:            finding.Service.Detection,
			DetectorId:           finding.Service.DetectorId,
			EbsVolumeScanDetails: finding.Service.EbsVolumeScanDetails,
			EventFirstSeen:       finding.Service.EventFirstSeen,
			EventLastSeen:        finding.Service.EventLastSeen,
			Evidence:             finding.Service.Evidence,
			FeatureName:          finding.Service.FeatureName,
			MalwareScanDetails:   finding.Service.MalwareScanDetails,
			ResourceRole:         finding.Service.ResourceRole,
			RuntimeDetails:       finding.Service.RuntimeDetails,
			ServiceName:          finding.Service.ServiceName,
			UserFeedback:         finding.Service.UserFeedback,
		}
	}

	// resource
	if finding.Resource != nil {
		row.Resource = &rows.GuardDutyFindingResource{
			ResourceType:         finding.Resource.ResourceType,
			AccessKeyDetails:     finding.Resource.AccessKeyDetails,
			ContainerDetails:     finding.Resource.ContainerDetails,
			EbsVolumeDetails:     finding.Resource.EbsVolumeDetails,
			EcsClusterDetails:    finding.Resource.EcsClusterDetails,
			EksClusterDetails:    finding.Resource.EksClusterDetails,
			InstanceDetails:      finding.Resource.InstanceDetails,
			KubernetesDetails:    finding.Resource.KubernetesDetails,
			LambdaDetails:        finding.Resource.LambdaDetails,
			RdsDbInstanceDetails: finding.Resource.RdsDbInstanceDetails,
			RdsDbUserDetails:     finding.Resource.RdsDbUserDetails,
			S3BucketDetails:      finding.Resource.S3BucketDetails,
		}
	}

	return row, nil
}
