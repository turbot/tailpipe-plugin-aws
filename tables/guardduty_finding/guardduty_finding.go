package guardduty_finding

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type GuardDutyFinding struct {
	schema.CommonFields

	AccountId     *string    `json:"account_id"`
	Arn           *string    `json:"arn"`
	Description   *string    `json:"description"`
	Id            *string    `json:"id"`
	Partition     *string    `json:"partition"`
	Region        *string    `json:"region"`
	SchemaVersion *string    `json:"schema_version"`
	Severity      *float64   `json:"severity"`
	Title         *string    `json:"title"`
	Type          *string    `json:"type"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Service       *Service   `json:"service,omitempty"`
	Resource      *Resource  `json:"resource,omitempty"`
}

type Service struct {
	Action               *Action                     `json:"action,omitempty"`
	AdditionalInfo       *ServiceAdditionInfo        `json:"additional_info,omitempty"`
	Archived             *bool                       `json:"archived,omitempty"`
	Count                *int32                      `json:"count,omitempty"`
	Detection            *types.Detection            `json:"detection,omitempty" parquet:"type=JSON"` // contains maps
	DetectorId           *string                     `json:"detector_id,omitempty"`
	EbsVolumeScanDetails *types.EbsVolumeScanDetails `json:"ebs_volume_scan_details,omitempty" parquet:"type=JSON"` // contains []struct
	EventFirstSeen       *string                     `json:"event_first_seen,omitempty"`
	EventLastSeen        *string                     `json:"event_last_seen,omitempty"`
	Evidence             *types.Evidence             `json:"evidence,omitempty" parquet:"type=JSON"` // contains []struct
	FeatureName          *string                     `json:"feature_name,omitempty"`
	MalwareScanDetails   *types.MalwareScanDetails   `json:"malware_scan_details,omitempty" parquet:"type=JSON"` // contains []struct
	ResourceRole         *string                     `json:"resource_role,omitempty"`
	RuntimeDetails       *RuntimeDetails             `json:"runtime_details,omitempty"`
	ServiceName          *string                     `json:"service_name,omitempty"`
	UserFeedback         *string                     `json:"user_feedback,omitempty"`
}

type ServiceAdditionInfo struct {
	Type  *string `json:"type"`
	Value *string `json:"value"`
}

type RuntimeDetails struct {
	Context *RuntimeDetailsContext `json:"context,omitempty"`
	Process *ProcessDetails        `json:"process,omitempty"`
}

type RuntimeDetailsContext struct {
	AddressFamily      *string         `json:"address_family,omitempty"`
	CommandLineExample *string         `json:"command_line_example,omitempty"`
	FileSystemType     *string         `json:"file_system_type,omitempty"`
	Flags              []string        `json:"flags,omitempty"`
	IanaProtocolNumber *int32          `json:"iana_protocol_number,omitempty"`
	LdPreloadValue     *string         `json:"ld_preload_value,omitempty"`
	LibraryPath        *string         `json:"library_path,omitempty"`
	MemoryRegions      []string        `json:"memory_regions,omitempty"`
	ModifiedAt         *time.Time      `json:"modified_at,omitempty"`
	ModifyingProcess   *ProcessDetails `json:"modifying_process,omitempty"`
	ModuleFilePath     *string         `json:"module_file_path,omitempty"`
	ModuleName         *string         `json:"module_name,omitempty"`
	ModuleSha256       *string         `json:"module_sha256,omitempty"`
	MountSource        *string         `json:"mount_source,omitempty"`
	MountTarget        *string         `json:"mount_target,omitempty"`
	ReleaseAgentPath   *string         `json:"release_agent_path,omitempty"`
	RuncBinaryPath     *string         `json:"runc_binary_path,omitempty"`
	ScriptPath         *string         `json:"script_path,omitempty"`
	ServiceName        *string         `json:"service_name,omitempty"`
	ShellHistoryPath   *string         `json:"shell_history_path,omitempty"`
	SocketPath         *string         `json:"socket_path,omitempty"`
	TargetProcess      *ProcessDetails `json:"target_process,omitempty"`
	ThreatFilePath     *string         `json:"threat_file_path,omitempty"`
	ToolCategory       *string         `json:"tool_category,omitempty"`
	ToolName           *string         `json:"tool_name,omitempty"`
}

type ProcessDetails struct {
	Euid             *int32     `json:"euid,omitempty"`
	ExecutablePath   *string    `json:"executable_path,omitempty"`
	ExecutableSha256 *string    `json:"executable_sha256,omitempty"`
	Name             *string    `json:"name,omitempty"`
	NamespacePid     *int32     `json:"namespace_pid,omitempty"`
	ParentUuid       *string    `json:"parent_uuid,omitempty"`
	Pid              *int32     `json:"pid,omitempty"`
	Pwd              *string    `json:"pwd,omitempty"`
	StartTime        *time.Time `json:"start_time,omitempty"`
	User             *string    `json:"user,omitempty"`
	UserId           *int32     `json:"user_id,omitempty"`
	Uuid             *string    `json:"uuid,omitempty"`
}

type Resource struct {
	ResourceType     *string                 `json:"resource_type,omitempty"`
	AccessKeyDetails *AccessKeyDetails       `json:"access_key_details,omitempty"`
	ResourceDetails  *map[string]interface{} `json:"resource_details,omitempty" parquet:"type=JSON"`
}

type Action struct {
	ActionType    *string                `json:"action_type,omitempty"`
	ActionDetails map[string]interface{} `json:"action_details,omitempty" parquet:"type=JSON"`
}

type AccessKeyDetails struct {
	AccessKeyId *string `json:"access_key_id,omitempty"`
	PrincipalId *string `json:"principal_id,omitempty"`
	UserName    *string `json:"user_name,omitempty"`
	UserType    *string `json:"user_type,omitempty"`
}

func (c *GuardDutyFinding) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"account_id":     "The AWS account ID where the finding was generated.",
		"arn":            "The Amazon Resource Name (ARN) of the finding.",
		"description":    "A detailed description of the security finding, including what was detected and its potential impact.",
		"id":             "The unique identifier for the finding.",
		"partition":      "The AWS partition where the finding was generated (e.g., aws, aws-cn, aws-gov).",
		"region":         "The AWS region where the finding was generated.",
		"schema_version": "The version of the GuardDuty finding format schema.",
		"severity":       "The severity level of the finding, ranging from 0.1 to 8.9, where higher values indicate higher severity.",
		"title":          "A short human-readable title for the finding that describes the type of activity detected.",
		"type":           "The type of finding, indicating the category of security issue detected.",
		"created_at":     "The timestamp when the finding was first created.",
		"updated_at":     "The timestamp when the finding was last updated.",
		"service":        "Details about the GuardDuty service that generated the finding, including actions, additional information, and detection details.",
		"resource":       "Information about the AWS resource that was targeted.",

		// Tailpipe-specific metadata fields
		"tp_akas":      "The Amazon Resource Names (ARNs) associated with the finding.",
		"tp_index":     "The AWS account ID where the finding was generated.",
		"tp_usernames": "Usernames associated with the finding, including IAM users and access key IDs.",
	}
}
