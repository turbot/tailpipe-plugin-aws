package guardduty_finding

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

// GuardDutyMapper is a mapper that receives GuardDutyBatch objects and extracts GuardDutyFinding records from them
type GuardDutyMapper struct {
}

func (g *GuardDutyMapper) Identifier() string {
	return "guardduty_mapper"
}

// Map casts the data item as a GuardDutyBatch and returns the GuardDutyFinding records
func (g *GuardDutyMapper) Map(_ context.Context, a any, _ ...mappers.MapOption[*GuardDutyFinding]) (*GuardDutyFinding, error) {
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
	row := &GuardDutyFinding{
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
		row.Service = &Service{
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
			ServiceName:          finding.Service.ServiceName,
			UserFeedback:         finding.Service.UserFeedback,
		}

		// service.action
		if finding.Service.Action != nil {
			row.Service.Action = &Action{
				ActionType: finding.Service.Action.ActionType,
			}

			// TODO: Temporarily removed err handling from convertToMap functions to fix linting, we should handle these errors instead though
			var details map[string]interface{}
			switch {
			case finding.Service.Action.AwsApiCallAction != nil:
				details, _ = convertToMap(finding.Service.Action.AwsApiCallAction)
			case finding.Service.Action.DnsRequestAction != nil:
				details, _ = convertToMap(finding.Service.Action.DnsRequestAction)
			case finding.Service.Action.NetworkConnectionAction != nil:
				details, _ = convertToMap(finding.Service.Action.NetworkConnectionAction)
			case finding.Service.Action.PortProbeAction != nil:
				details, _ = convertToMap(finding.Service.Action.PortProbeAction)
			case finding.Service.Action.KubernetesApiCallAction != nil:
				details, _ = convertToMap(finding.Service.Action.KubernetesApiCallAction)
			case finding.Service.Action.KubernetesPermissionCheckedDetails != nil:
				details, _ = convertToMap(finding.Service.Action.KubernetesPermissionCheckedDetails)
			case finding.Service.Action.KubernetesRoleDetails != nil:
				details, _ = convertToMap(finding.Service.Action.KubernetesRoleDetails)
			case finding.Service.Action.KubernetesRoleBindingDetails != nil:
				details, _ = convertToMap(finding.Service.Action.KubernetesRoleBindingDetails)
			case finding.Service.Action.RdsLoginAttemptAction != nil:
				details, _ = convertToMap(finding.Service.Action.RdsLoginAttemptAction)
			}

			row.Service.Action.ActionDetails = details
		}

		// service.additionInfo
		if finding.Service.AdditionalInfo != nil {
			row.Service.AdditionalInfo = &ServiceAdditionInfo{
				Type:  finding.Service.AdditionalInfo.Type,
				Value: finding.Service.AdditionalInfo.Value,
			}
		}

		// service.runtimeDetails
		if finding.Service.RuntimeDetails != nil {
			var runtimeDetails RuntimeDetails

			// service.runtimeDetails.context
			if finding.Service.RuntimeDetails.Context != nil {
				runtimeDetails.Context = &RuntimeDetailsContext{
					AddressFamily:      finding.Service.RuntimeDetails.Context.AddressFamily,
					CommandLineExample: finding.Service.RuntimeDetails.Context.CommandLineExample,
					FileSystemType:     finding.Service.RuntimeDetails.Context.FileSystemType,
					Flags:              finding.Service.RuntimeDetails.Context.Flags,
					IanaProtocolNumber: finding.Service.RuntimeDetails.Context.IanaProtocolNumber,
					LdPreloadValue:     finding.Service.RuntimeDetails.Context.LdPreloadValue,
					LibraryPath:        finding.Service.RuntimeDetails.Context.LibraryPath,
					MemoryRegions:      finding.Service.RuntimeDetails.Context.MemoryRegions,
					ModifiedAt:         finding.Service.RuntimeDetails.Context.ModifiedAt,
					ModuleFilePath:     finding.Service.RuntimeDetails.Context.ModuleFilePath,
					ModuleName:         finding.Service.RuntimeDetails.Context.ModuleName,
					ModuleSha256:       finding.Service.RuntimeDetails.Context.ModuleSha256,
					MountSource:        finding.Service.RuntimeDetails.Context.MountSource,
					MountTarget:        finding.Service.RuntimeDetails.Context.MountTarget,
					ReleaseAgentPath:   finding.Service.RuntimeDetails.Context.ReleaseAgentPath,
					RuncBinaryPath:     finding.Service.RuntimeDetails.Context.RuncBinaryPath,
					ScriptPath:         finding.Service.RuntimeDetails.Context.ScriptPath,
					ServiceName:        finding.Service.RuntimeDetails.Context.ServiceName,
					ShellHistoryPath:   finding.Service.RuntimeDetails.Context.ShellHistoryFilePath,
					SocketPath:         finding.Service.RuntimeDetails.Context.SocketPath,
					ThreatFilePath:     finding.Service.RuntimeDetails.Context.ThreatFilePath,
					ToolCategory:       finding.Service.RuntimeDetails.Context.ToolCategory,
					ToolName:           finding.Service.RuntimeDetails.Context.ToolName,
				}

				// service.runtimeDetails.context.modifyingProcess
				if finding.Service.RuntimeDetails.Context.ModifyingProcess != nil {
					runtimeDetails.Context.ModifyingProcess = &ProcessDetails{
						Euid:             finding.Service.RuntimeDetails.Context.ModifyingProcess.Euid,
						ExecutablePath:   finding.Service.RuntimeDetails.Context.ModifyingProcess.ExecutablePath,
						ExecutableSha256: finding.Service.RuntimeDetails.Context.ModifyingProcess.ExecutableSha256,
						Name:             finding.Service.RuntimeDetails.Context.ModifyingProcess.Name,
						NamespacePid:     finding.Service.RuntimeDetails.Context.ModifyingProcess.NamespacePid,
						ParentUuid:       finding.Service.RuntimeDetails.Context.ModifyingProcess.ParentUuid,
						Pid:              finding.Service.RuntimeDetails.Context.ModifyingProcess.Pid,
						Pwd:              finding.Service.RuntimeDetails.Context.ModifyingProcess.Pwd,
						StartTime:        finding.Service.RuntimeDetails.Context.ModifyingProcess.StartTime,
						User:             finding.Service.RuntimeDetails.Context.ModifyingProcess.User,
						UserId:           finding.Service.RuntimeDetails.Context.ModifyingProcess.UserId,
						Uuid:             finding.Service.RuntimeDetails.Context.ModifyingProcess.Uuid,
					}
				}

				// service.runtimeDetails.context.targetProcess
				if finding.Service.RuntimeDetails.Context.TargetProcess != nil {
					runtimeDetails.Context.TargetProcess = &ProcessDetails{
						Euid:             finding.Service.RuntimeDetails.Context.TargetProcess.Euid,
						ExecutablePath:   finding.Service.RuntimeDetails.Context.TargetProcess.ExecutablePath,
						ExecutableSha256: finding.Service.RuntimeDetails.Context.TargetProcess.ExecutableSha256,
						Name:             finding.Service.RuntimeDetails.Context.TargetProcess.Name,
						NamespacePid:     finding.Service.RuntimeDetails.Context.TargetProcess.NamespacePid,
						ParentUuid:       finding.Service.RuntimeDetails.Context.TargetProcess.ParentUuid,
						Pid:              finding.Service.RuntimeDetails.Context.TargetProcess.Pid,
						Pwd:              finding.Service.RuntimeDetails.Context.TargetProcess.Pwd,
						StartTime:        finding.Service.RuntimeDetails.Context.TargetProcess.StartTime,
						User:             finding.Service.RuntimeDetails.Context.TargetProcess.User,
						UserId:           finding.Service.RuntimeDetails.Context.TargetProcess.UserId,
						Uuid:             finding.Service.RuntimeDetails.Context.TargetProcess.Uuid,
					}
				}
			}

			// service.runtimeDetails.process
			if finding.Service.RuntimeDetails.Process != nil {
				runtimeDetails.Process = &ProcessDetails{
					Euid:             finding.Service.RuntimeDetails.Process.Euid,
					ExecutablePath:   finding.Service.RuntimeDetails.Process.ExecutablePath,
					ExecutableSha256: finding.Service.RuntimeDetails.Process.ExecutableSha256,
					Name:             finding.Service.RuntimeDetails.Process.Name,
					NamespacePid:     finding.Service.RuntimeDetails.Process.NamespacePid,
					ParentUuid:       finding.Service.RuntimeDetails.Process.ParentUuid,
					Pid:              finding.Service.RuntimeDetails.Process.Pid,
					Pwd:              finding.Service.RuntimeDetails.Process.Pwd,
					StartTime:        finding.Service.RuntimeDetails.Process.StartTime,
					User:             finding.Service.RuntimeDetails.Process.User,
					UserId:           finding.Service.RuntimeDetails.Process.UserId,
					Uuid:             finding.Service.RuntimeDetails.Process.Uuid,
				}

				row.Service.RuntimeDetails = &runtimeDetails
			}
		}
	}

	// resource
	if finding.Resource != nil {
		row.Resource = &Resource{
			ResourceType: finding.Resource.ResourceType,
		}

		if finding.Resource.AccessKeyDetails != nil {
			row.Resource.AccessKeyDetails = &AccessKeyDetails{
				AccessKeyId: finding.Resource.AccessKeyDetails.AccessKeyId,
				UserName:    finding.Resource.AccessKeyDetails.UserName,
				UserType:    finding.Resource.AccessKeyDetails.UserType,
				PrincipalId: finding.Resource.AccessKeyDetails.PrincipalId,
			}
		}

		// TODO: Temporarily removed err handling from convertToMap functions to fix linting, we should handle these errors instead though
		var details map[string]interface{}
		switch {
		case finding.Resource.ContainerDetails != nil:
			details, _ = convertToMap(finding.Resource.ContainerDetails)
		case finding.Resource.EbsVolumeDetails != nil:
			details, _ = convertToMap(finding.Resource.EbsVolumeDetails)
		case finding.Resource.EcsClusterDetails != nil:
			details, _ = convertToMap(finding.Resource.EcsClusterDetails)
		case finding.Resource.EksClusterDetails != nil:
			details, _ = convertToMap(finding.Resource.EksClusterDetails)
		case finding.Resource.InstanceDetails != nil:
			details, _ = convertToMap(finding.Resource.InstanceDetails)
		case finding.Resource.KubernetesDetails != nil:
			details, _ = convertToMap(finding.Resource.KubernetesDetails)
		case finding.Resource.LambdaDetails != nil:
			details, _ = convertToMap(finding.Resource.LambdaDetails)
		case finding.Resource.RdsDbInstanceDetails != nil:
			details, _ = convertToMap(finding.Resource.RdsDbInstanceDetails)
		case finding.Resource.RdsDbUserDetails != nil:
			details, _ = convertToMap(finding.Resource.RdsDbUserDetails)
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
