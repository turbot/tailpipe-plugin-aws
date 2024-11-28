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
	EbsVolumeScanDetails *types.EbsVolumeScanDetails `json:"ebs_volume_scan_details,omitempty"`
	EventFirstSeen       *string                     `json:"event_first_seen,omitempty"`
	EventLastSeen        *string                     `json:"event_last_seen,omitempty"`
	Evidence             *types.Evidence             `json:"evidence,omitempty"`
	FeatureName          *string                     `json:"feature_name,omitempty"`
	MalwareScanDetails   *types.MalwareScanDetails   `json:"malware_scan_details,omitempty"`
	ResourceRole         *string                     `json:"resource_role,omitempty"`
	RuntimeDetails       *types.RuntimeDetails       `json:"runtime_details,omitempty"`
	ServiceName          *string                     `json:"service_name,omitempty"`
	UserFeedback         *string                     `json:"user_feedback,omitempty"`
}

type GuardDutyFindingResource struct {
	ResourceType         *string                     `json:"resource_type,omitempty"`
	AccessKeyDetails     *types.AccessKeyDetails     `json:"access_key_details,omitempty"`
	ContainerDetails     *types.Container            `json:"container_details,omitempty"`
	EbsVolumeDetails     *types.EbsVolumeDetails     `json:"ebs_volume_details,omitempty"`
	EcsClusterDetails    *types.EcsClusterDetails    `json:"ecs_cluster_details,omitempty"`
	EksClusterDetails    *types.EksClusterDetails    `json:"eks_cluster_details,omitempty"`
	InstanceDetails      *types.InstanceDetails      `json:"instance_details,omitempty"`
	KubernetesDetails    *types.KubernetesDetails    `json:"kubernetes_details,omitempty"`
	LambdaDetails        *types.LambdaDetails        `json:"lambda_details,omitempty"`
	RdsDbInstanceDetails *types.RdsDbInstanceDetails `json:"rds_db_instance_details,omitempty"`
	RdsDbUserDetails     *types.RdsDbUserDetails     `json:"rds_db_user_details,omitempty"`
	S3BucketDetails      []types.S3BucketDetail      `json:"s3_bucket_details,omitempty"`
}
