package models

import (
	"fmt"
	"github.com/turbot/tailpipe-plugin-sdk/enrichment"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

// TODO is trhere an existing amazon sdk type we can use
type AwsVpcFlowLog struct {
	// embed required enrichment fields (be sure to skip in parquet)
	enrichment.CommonFields

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
	Start                   *int64     `json:"start,omitempty"`
	End                     *int64     `json:"end,omitempty"`
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

// fromString

func FlowLogFromString(rowString string, schema []string) (*AwsVpcFlowLog, error) {

	fields := strings.Fields(rowString)

	if len(fields) > len(schema) {
		slog.Error("row has more fields than schema allows", "fields", fields, "schema", schema)
		return nil, fmt.Errorf("row has more fields than schema allows")
	}

	flowLog := &AwsVpcFlowLog{}
	for i, field := range fields {
		// skip empty fields
		if field == "-" {
			continue
		}
		switch schema[i] {
		case "timestamp":
			timestamp, err := time.Parse(time.RFC3339, field)
			if err != nil {
				return nil, fmt.Errorf("invalid timestamp: %s", field)
			}
			flowLog.Timestamp = &timestamp
		case "version":
			version, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid version: %s", field)
			}
			v := int32(version)
			flowLog.Version = &v
		case "account-id":
			flowLog.AccountID = &field
		case "interface-id":
			flowLog.InterfaceID = &field
		case "srcaddr":
			flowLog.SrcAddr = &field
		case "dstaddr":
			flowLog.DstAddr = &field
		case "srcport":
			srcPort, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid srcport: %s", field)
			}
			v := int32(srcPort)
			flowLog.SrcPort = &v
		case "dstport":
			dstPort, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid dstport: %s", field)
			}
			v := int32(dstPort)
			flowLog.DstPort = &v
		case "protocol":
			protocol, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid protocol: %s", field)
			}
			v := int32(protocol)
			flowLog.Protocol = &v
		case "packets":
			packets, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid packets: %s", field)
			}
			flowLog.Packets = &packets
		case "bytes":
			bytes, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid bytes: %s", field)
			}
			flowLog.Bytes = &bytes
		case "start":
			start, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid start: %s", field)
			}
			flowLog.Start = &start
		case "end":
			end, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid end: %s", field)
			}
			flowLog.End = &end
		case "action":
			flowLog.Action = &field
		case "log-status":
			flowLog.LogStatus = &field
		case "vpc-id":
			flowLog.VPCID = &field
		case "subnet-id":
			flowLog.SubnetID = &field
		case "instance-id":
			flowLog.InstanceID = &field
		case "tcp-flags":
			tcpFlags, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid tcp-flags: %s", field)
			}
			v := int32(tcpFlags)
			flowLog.TCPFlags = &v
		case "type":
			flowLog.Type = &field
		case "pkt-srcaddr":
			flowLog.PktSrcAddr = &field
		case "pkt-dstaddr":
			flowLog.PktDstAddr = &field
		case "region":
			flowLog.Region = &field
		case "az-id":
			flowLog.AzID = &field
		case "sublocation-type":
			flowLog.SublocationType = &field
		case "sublocation-id":
			flowLog.SublocationID = &field
		case "pkt-src-aws-service":
			flowLog.PktSrcAWSService = &field
		case "pkt-dst-aws-service":
			flowLog.PktDstAWSService = &field
		case "flow-direction":
			flowLog.FlowDirection = &field
		case "traffic-path":
			trafficPath, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid traffic-path: %s", field)
			}
			v := int32(trafficPath)
			flowLog.TrafficPath = &v
		case "ecs-cluster-arn":
			flowLog.ECSClusterARN = &field
		case "ecs-cluster-name":
			flowLog.ECSClusterName = &field
		case "ecs-container-instance-arn":
			flowLog.ECSContainerInstanceARN = &field
		case "ecs-container-instance-id":
			flowLog.ECSContainerInstanceID = &field
		case "ecs-container-id":
			flowLog.ECSContainerID = &field
		case "ecs-second-container-id":
			flowLog.ECSSecondContainerID = &field
		case "ecs-service-name":
			flowLog.ECSServiceName = &field
		case "ecs-task-definition-arn":
			flowLog.ECSTaskDefinitionARN = &field
		case "ecs-task-arn":
			flowLog.ECSTaskARN = &field
		case "ecs-task-id":
			flowLog.ECSTaskID = &field
		default:
			return nil, fmt.Errorf("unknown field: %s", schema[i])
		}
	}

	return flowLog, nil
}
