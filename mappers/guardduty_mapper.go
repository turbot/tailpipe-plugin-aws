package mappers

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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

	// pre-process JSON for compatibility issue(s) with AWS SDK https://github.com/aws/aws-sdk-go-v2/issues/2145
	jsonBytes = preprocessJSON(jsonBytes)

	// Decode JSON into AWS SDK `types.Finding`
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
			ResourceType: finding.Resource.ResourceType,
		}

		if finding.Resource.AccessKeyDetails != nil {
			row.Resource.AccessKeyDetails = &rows.AccessKeyDetails{
				AccessKeyId: finding.Resource.AccessKeyDetails.AccessKeyId,
				UserName:    finding.Resource.AccessKeyDetails.UserName,
				UserType:    finding.Resource.AccessKeyDetails.UserType,
				PrincipalId: finding.Resource.AccessKeyDetails.PrincipalId,
			}
		}

		var details map[string]interface{}
		switch {
		case finding.Resource.ContainerDetails != nil:
			details, err = convertToMap(finding.Resource.ContainerDetails)
		case finding.Resource.EbsVolumeDetails != nil:
			details, err = convertToMap(finding.Resource.EbsVolumeDetails)
		case finding.Resource.EcsClusterDetails != nil:
			details, err = convertToMap(finding.Resource.EcsClusterDetails)
		case finding.Resource.EksClusterDetails != nil:
			details, err = convertToMap(finding.Resource.EksClusterDetails)
		case finding.Resource.InstanceDetails != nil:
			details, err = convertToMap(finding.Resource.InstanceDetails)
		case finding.Resource.KubernetesDetails != nil:
			details, err = convertToMap(finding.Resource.KubernetesDetails)
		case finding.Resource.LambdaDetails != nil:
			details, err = convertToMap(finding.Resource.LambdaDetails)
		case finding.Resource.RdsDbInstanceDetails != nil:
			details, err = convertToMap(finding.Resource.RdsDbInstanceDetails)
		case finding.Resource.RdsDbUserDetails != nil:
			details, err = convertToMap(finding.Resource.RdsDbUserDetails)
		case finding.Resource.S3BucketDetails != nil:
			details = map[string]interface{}{
				"s3_bucket_details": finding.Resource.S3BucketDetails,
			}
		}

		row.Resource.ResourceDetails = &details
	}

	return row, nil
}

func preprocessJSON(input []byte) []byte {
	re := regexp.MustCompile(`"([a-zA-Z]+(?:At|Time|Seen))":(\d+(\.\d+)?E[+-]?\d+)`)
	return re.ReplaceAllFunc(input, func(match []byte) []byte {
		// Extract the field name and numeric value
		groups := re.FindSubmatch(match)
		fieldName := groups[1]        // "createdAt" or "updatedAt"
		numericTimestamp := groups[2] // e.g., 1.636625755218E9

		// Parse the numeric timestamp into an integer
		seconds, err := strconv.ParseFloat(string(numericTimestamp), 64)
		if err != nil {
			// If parsing fails, return the original match unchanged
			return match
		}

		// Convert seconds to ISO 8601 format
		timestamp := time.Unix(int64(seconds), 0).UTC().Format(time.RFC3339)
		return []byte(fmt.Sprintf(`"%s":"%s"`, fieldName, timestamp))
	})
}

func convertToMap(input any) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
