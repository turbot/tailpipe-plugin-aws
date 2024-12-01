package rows

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"

	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
)

type GuardDutyFinding struct {
	enrichment.CommonFields

	AccountId     *string                   `json:"account_id"`
	Arn           *string                   `json:"arn"`
	Description   *string                   `json:"description"`
	Id            *string                   `json:"id"`
	Partition     *string                   `json:"partition"`
	Region        *string                   `json:"region"`
	SchemaVersion *string                   `json:"schema_version"`
	Severity      *float64                  `json:"severity"`
	Title         *string                   `json:"title"`
	Type          *string                   `json:"type"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     *time.Time                `json:"updated_at"`
	Service       *GuardDutyFindingService  `json:"service,omitempty"`
	Resource      *GuardDutyFindingResource `json:"resource,omitempty"`
}

type GuardDutyFindingService struct {
	Action               *types.Action               `json:"action,omitempty" parquet:"type=JSON"` // contains maps
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
	RuntimeDetails       *types.RuntimeDetails       `json:"runtime_details,omitempty" parquet:"type=JSON"` // contains []struct
	ServiceName          *string                     `json:"service_name,omitempty"`
	UserFeedback         *string                     `json:"user_feedback,omitempty"`
}

type GuardDutyFindingResource struct {
	ResourceType     *string                 `json:"resource_type,omitempty"`
	AccessKeyDetails *AccessKeyDetails       `json:"access_key_details,omitempty"`
	ResourceDetails  *map[string]interface{} `json:"resource_details,omitempty" parquet:"type=JSON"`
}

type AccessKeyDetails struct {
	AccessKeyId *string `json:"access_key_id,omitempty"`
	PrincipalId *string `json:"principal_id,omitempty"`
	UserName    *string `json:"user_name,omitempty"`
	UserType    *string `json:"user_type,omitempty"`
}
