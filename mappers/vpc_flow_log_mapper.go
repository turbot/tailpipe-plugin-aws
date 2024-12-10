package mappers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-aws/rows"
)

// VpcFlowLogMapper is a mapper that receives string objects and extracts VpcFlowLog record
type VpcFlowLogMapper struct {
	schema []string
}

func NewVpcFlowLogMapper(schema []string) *VpcFlowLogMapper {
	return &VpcFlowLogMapper{
		schema: schema,
	}

}

func (c *VpcFlowLogMapper) Identifier() string {
	return "vpc_flow_log_mapper"
}

// Map casts the data item as an string and returns the VpcFlowLog records
func (c *VpcFlowLogMapper) Map(_ context.Context, a any) (*rows.VpcFlowLog, error) {
	rowString, ok := a.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %T", a)
	}

	fields := strings.Fields(rowString)

	if len(fields) > len(c.schema) {
		slog.Error("row has more fields than c.schema allows", "fields", fields, "schema", c.schema)
		return nil, fmt.Errorf("row has more fields than c.schema allows")
	}

	flowLog := &rows.VpcFlowLog{}
	for i, field := range fields {
		// skip empty fields
		if field == "-" {
			continue
		}
		switch c.schema[i] {
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
			return nil, fmt.Errorf("unknown field: %s", c.schema[i])
		}
	}

	return flowLog, nil

}
