package vpc_flow_log

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/mappers"
)

type VPCFlowLogMapper struct {
}

func (m *VPCFlowLogMapper) Identifier() string {
	return "vpc_flow_log_mapper"
}

func (m *VPCFlowLogMapper) Map(_ context.Context, a any, _ ...mappers.MapOption[*VpcFlowLog]) (*VpcFlowLog, error) {
	var log *VpcFlowLog
	var logLines string
	var err error

	switch v := a.(type) {
	case *VpcFlowLog:
		return v, nil
	case VpcFlowLog:
		return &v, nil
	case []byte:
		logLines = string(v)
	case *string:
		logLines = *v
	case string:
		logLines = v
	default:
		return nil, fmt.Errorf("expected byte[], string or rows.CloudTailLog got %T", a)
	}

	mapSlice, err := ConvertToMapSlice(logLines)
	if err != nil {
		return nil, fmt.Errorf("error converting to map slice: %w", err)
	}

	log = &VpcFlowLog{}
	for _, mapSliceValue := range mapSlice {
		err = log.MapValues(mapSliceValue)
		if err != nil {
			return nil, fmt.Errorf("error mapping values: %w", err)
		}
	}

	return log, nil
}

// InitialiseFromMap initializes a VpcFlowLog struct from a map of string values
func (l *VpcFlowLog) InitialiseFromMap(m map[string]string) error {
	for key, value := range m {
		if value == "-" {
			continue
		}
		switch key {
		case "account-id":
			l.AccountID = &value
		case "action":
			l.Action = &value
		case "az-id":
			l.AzID = &value
		case "bytes":
			if bytes, err := strconv.ParseInt(value, 10, 64); err == nil {
				l.Bytes = &bytes
			}
		case "dstaddr":
			l.DstAddr = &value
		case "dstport":
			if port, err := strconv.ParseInt(value, 10, 32); err == nil {
				port32 := int32(port)
				l.DstPort = &port32
			}
		case "end":
			if ts, err := strconv.ParseInt(value, 10, 64); err == nil {
				t := time.Unix(ts, 0)
				l.End = &t
			}
		case "flow-direction":
			l.FlowDirection = &value
		case "instance-id":
			l.InstanceID = &value
		case "interface-id":
			l.InterfaceID = &value
		case "log-status":
			l.LogStatus = &value
		case "packets":
			if packets, err := strconv.ParseInt(value, 10, 64); err == nil {
				l.Packets = &packets
			}
		case "pkt-dst-aws-service":
			l.PktDstAWSService = &value
		case "pkt-dstaddr":
			l.PktDstAddr = &value
		case "pkt-src-aws-service":
			l.PktSrcAWSService = &value
		case "pkt-srcaddr":
			l.PktSrcAddr = &value
		case "protocol":
			if proto, err := strconv.ParseInt(value, 10, 32); err == nil {
				proto32 := int32(proto)
				l.Protocol = &proto32
			}
		case "region":
			l.Region = &value
		case "reject-reason":
			l.RejectReason = &value
		case "srcaddr":
			l.SrcAddr = &value
		case "srcport":
			if port, err := strconv.ParseInt(value, 10, 32); err == nil {
				port32 := int32(port)
				l.SrcPort = &port32
			}
		case "start":
			if ts, err := strconv.ParseInt(value, 10, 64); err == nil {
				t := time.Unix(ts, 0)
				l.Start = &t
			}
		case "sublocation-id":
			l.SublocationID = &value
		case "sublocation-type":
			l.SublocationType = &value
		case "subnet-id":
			l.SubnetID = &value
		case "tcp-flags":
			if flags, err := strconv.ParseInt(value, 10, 32); err == nil {
				flags32 := int32(flags)
				l.TCPFlags = &flags32
			}
		case "traffic-path":
			if path, err := strconv.ParseInt(value, 10, 32); err == nil {
				path32 := int32(path)
				l.TrafficPath = &path32
			}
		case "type":
			l.Type = &value
		case "version":
			if ver, err := strconv.ParseInt(value, 10, 32); err == nil {
				ver32 := int32(ver)
				l.Version = &ver32
			}
		case "vpc-id":
			l.VPCID = &value
		case "ecs-cluster-name":
			l.ECSClusterName = &value
		case "ecs-cluster-arn":
			l.ECSClusterARN = &value
		case "ecs-container-instance-id":
			l.ECSContainerInstanceID = &value
		case "ecs-container-instance-arn":
			l.ECSContainerInstanceARN = &value
		case "ecs-service-name":
			l.ECSServiceName = &value
		case "ecs-task-definition-arn":
			l.ECSTaskDefinitionARN = &value
		case "ecs-task-id":
			l.ECSTaskID = &value
		case "ecs-task-arn":
			l.ECSTaskARN = &value
		case "ecs-container-id":
			l.ECSContainerID = &value
		case "ecs-second-container-id":
			l.ECSSecondContainerID = &value
		}
	}

	return nil
}
