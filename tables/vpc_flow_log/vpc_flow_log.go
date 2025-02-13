package vpc_flow_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type VpcFlowLog struct {
	// embed required enrichment fields
	schema.CommonFields

	Timestamp               *time.Time `json:"timestamp,omitempty"`
	Version                 *int32     `json:"version,omitempty"`
	AccountID               *string    `json:"account_id,omitempty"`
	InterfaceID             *string    `json:"interface_id,omitempty"`
	SrcAddr                 *string    `json:"src_addr,omitempty"`
	DstAddr                 *string    `json:"dst_addr,omitempty"`
	SrcPort                 *int32     `json:"src_port,omitempty"`
	DstPort                 *int32     `json:"dst_port,omitempty"`
	Protocol                *int32     `json:"protocol,omitempty"`
	Packets                 *int64     `json:"packets,omitempty"`
	Bytes                   *int64     `json:"bytes,omitempty"`
	Start                   *int64     `json:"start_time,omitempty"`
	End                     *int64     `json:"end_time,omitempty"`
	Action                  *string    `json:"action,omitempty"`
	LogStatus               *string    `json:"log_status,omitempty"`
	VPCID                   *string    `json:"vpc_id,omitempty"`
	SubnetID                *string    `json:"subnet_id,omitempty"`
	InstanceID              *string    `json:"instance_id,omitempty"`
	TCPFlags                *int32     `json:"tcp_flags,omitempty"`
	Type                    *string    `json:"type,omitempty"`
	PktSrcAddr              *string    `json:"pkt_src_addr,omitempty"`
	PktDstAddr              *string    `json:"pkt_dst_addr,omitempty"`
	Region                  *string    `json:"region,omitempty"`
	AzID                    *string    `json:"az_id,omitempty"`
	SublocationType         *string    `json:"sublocation_type,omitempty"`
	SublocationID           *string    `json:"sublocation_id,omitempty"`
	PktSrcAWSService        *string    `json:"pkt_src_aws_service,omitempty"`
	PktDstAWSService        *string    `json:"pkt_dst_aws_service,omitempty"`
	FlowDirection           *string    `json:"flow_direction,omitempty"`
	TrafficPath             *int32     `json:"traffic_path,omitempty"`
	ECSClusterARN           *string    `json:"ecs_cluster_arn,omitempty"`
	ECSClusterName          *string    `json:"ecs_cluster_name,omitempty"`
	ECSContainerInstanceARN *string    `json:"ecs_container_instance_arn,omitempty"`
	ECSContainerInstanceID  *string    `json:"ecs_container_instance_id,omitempty"`
	ECSContainerID          *string    `json:"ecs_container_id,omitempty"`
	ECSSecondContainerID    *string    `json:"ecs_second_container_id,omitempty"`
	ECSServiceName          *string    `json:"ecs_service_name,omitempty"`
	ECSTaskDefinitionARN    *string    `json:"ecs_task_definition_arn,omitempty"`
	ECSTaskARN              *string    `json:"ecs_task_arn,omitempty"`
	ECSTaskID               *string    `json:"ecs_task_id,omitempty"`
}
