package vpc_flow_log

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/schema"
)

type VpcFlowLog struct {
	// embed required enrichment fields
	schema.CommonFields

	AccountID               *string    `json:"account_id,omitempty"`
	Action                  *string    `json:"action,omitempty"`
	AzID                    *string    `json:"az_id,omitempty"`
	Bytes                   *int64     `json:"bytes,omitempty"`
	DstAddr                 *string    `json:"dst_addr,omitempty"`
	DstPort                 *int32     `json:"dst_port,omitempty"`
	ECSClusterARN           *string    `json:"ecs_cluster_arn,omitempty"`
	ECSClusterName          *string    `json:"ecs_cluster_name,omitempty"`
	ECSContainerID          *string    `json:"ecs_container_id,omitempty"`
	ECSContainerInstanceARN *string    `json:"ecs_container_instance_arn,omitempty"`
	ECSContainerInstanceID  *string    `json:"ecs_container_instance_id,omitempty"`
	ECSSecondContainerID    *string    `json:"ecs_second_container_id,omitempty"`
	ECSServiceName          *string    `json:"ecs_service_name,omitempty"`
	ECSTaskARN              *string    `json:"ecs_task_arn,omitempty"`
	ECSTaskDefinitionARN    *string    `json:"ecs_task_definition_arn,omitempty"`
	ECSTaskID               *string    `json:"ecs_task_id,omitempty"`
	End                     *time.Time `json:"end_time,omitempty"`
	FlowDirection           *string    `json:"flow_direction,omitempty"`
	InstanceID              *string    `json:"instance_id,omitempty"`
	InterfaceID             *string    `json:"interface_id,omitempty"`
	LogStatus               *string    `json:"log_status,omitempty"`
	Packets                 *int64     `json:"packets,omitempty"`
	PktDstAddr              *string    `json:"pkt_dst_addr,omitempty"`
	PktDstAWSService        *string    `json:"pkt_dst_aws_service,omitempty"`
	PktSrcAddr              *string    `json:"pkt_src_addr,omitempty"`
	PktSrcAWSService        *string    `json:"pkt_src_aws_service,omitempty"`
	Protocol                *int32     `json:"protocol,omitempty"`
	Region                  *string    `json:"region,omitempty"`
	SrcAddr                 *string    `json:"src_addr,omitempty"`
	SrcPort                 *int32     `json:"src_port,omitempty"`
	Start                   *time.Time `json:"start_time,omitempty"`
	SublocationID           *string    `json:"sublocation_id,omitempty"`
	SublocationType         *string    `json:"sublocation_type,omitempty"`
	SubnetID                *string    `json:"subnet_id,omitempty"`
	TCPFlags                *int32     `json:"tcp_flags,omitempty"`
	Timestamp               *time.Time `json:"timestamp,omitempty"`
	TrafficPath             *int32     `json:"traffic_path,omitempty"`
	Type                    *string    `json:"type,omitempty"`
	Version                 *int32     `json:"version,omitempty"`
	VPCID                   *string    `json:"vpc_id,omitempty"`
}

func (c *VpcFlowLog) GetColumnDescriptions() map[string]string {
	return map[string]string{
		"account_id":                 "The AWS account ID of the network interface owner.",
		"action":                     "The action associated with the traffic: ACCEPT or REJECT.",
		"az_id":                      "The Availability Zone ID that contains the network interface for which traffic is recorded.",
		"bytes":                      "The number of bytes transferred during the flow.",
		"dst_addr":                   "The destination IP address of the traffic. For incoming traffic, this is the private IP of the network interface.",
		"dst_port":                   "The destination port of the traffic.",
		"ecs_cluster_arn":            "The ARN of the ECS cluster if traffic originates from an ECS task.",
		"ecs_cluster_name":           "The name of the ECS cluster if traffic originates from an ECS task.",
		"ecs_container_id":           "The Docker runtime ID of the first container in an ECS task.",
		"ecs_container_instance_arn": "The ARN of the ECS container instance if traffic is from a running ECS task on an EC2 instance.",
		"ecs_container_instance_id":  "The ID of the ECS container instance if traffic is from a running ECS task on an EC2 instance.",
		"ecs_second_container_id":    "The Docker runtime ID of the second container in an ECS task, if applicable.",
		"ecs_service_name":           "The name of the ECS service if traffic is from an ECS task started by a service.",
		"ecs_task_arn":               "The ARN of the ECS task if traffic is from a running ECS task.",
		"ecs_task_definition_arn":    "The ARN of the ECS task definition if traffic is from a running ECS task.",
		"ecs_task_id":                "The ID of the ECS task if traffic is from a running ECS task.",
		"end_time":                   "The time, in Unix seconds, when the last packet of the flow was received within the aggregation interval.",
		"flow_direction":             "The direction of the flow with respect to the network interface: ingress or egress.",
		"instance_id":                "The ID of the EC2 instance associated with the network interface, if applicable.",
		"interface_id":               "The ID of the network interface for which traffic is recorded.",
		"log_status":                 "The logging status of the flow log: OK, NODATA, or SKIPDATA.",
		"packets":                    "The number of packets transferred during the flow.",
		"pkt_dst_addr":               "The packet-level (original) destination IP address of the traffic.",
		"pkt_dst_aws_service":        "The AWS service associated with the packet destination IP, if applicable.",
		"pkt_src_addr":               "The packet-level (original) source IP address of the traffic.",
		"pkt_src_aws_service":        "The AWS service associated with the packet source IP, if applicable.",
		"protocol":                   "The IANA protocol number of the traffic (e.g., 6 for TCP, 17 for UDP).",
		"region":                     "The AWS region containing the network interface for which traffic is recorded.",
		"reject_reason":              "The reason why traffic was rejected (e.g., BPA).",
		"src_addr":                   "The source IP address of the traffic. For outgoing traffic, this is the private IP of the network interface.",
		"src_port":                   "The source port of the traffic.",
		"start_time":                 "The time, in Unix seconds, when the first packet of the flow was received within the aggregation interval.",
		"sublocation_id":             "The ID of the sublocation that contains the network interface.",
		"sublocation_type":           "The type of sublocation that contains the network interface (e.g., wavelength, outpost, localzone).",
		"subnet_id":                  "The ID of the subnet that contains the network interface for which traffic is recorded.",
		"tcp_flags":                  "The bitmask value for TCP flags recorded during the flow.",
		"timestamp":                  "The timestamp when the request was received.",
		"traffic_path":               "The path that egress traffic takes to the destination.",
		"type":                       "The type of traffic (IPv4, IPv6, or EFA).",
		"version":                    "The VPC Flow Logs version. The version depends on the fields included in the log.",
		"vpc_id":                     "The ID of the VPC that contains the network interface for which traffic is recorded.",

		// Override table specific tp_* column descriptions
		"tp_akas":      "Resource ARNs related to the event. Possible values include ECSClusterARN, ECSContainerInstanceARN, ECSTaskARN, and/or ECSTaskDefinitionARN.",
		"tp_index":     "The AWS account ID that received the request. The default value is default.",
		"tp_ips":       "IP addresses associated with the event, including the source and destination IP addresses, private IPv4 or IPv6 addresses of the network interface for outgoing traffic, packet-level source and destination addresses to distinguish between intermediate and original or final destinations, and addresses affected by NAT gateways or Amazon EKS pod communications within a VPC.",
		"tp_timestamp": "The date and time the event occurred, in ISO 8601 format; if the timestamp is not available, the current timestamp is used.",
	}
}
