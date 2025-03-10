package vpc_flow_log

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/artifact_source"
)

// VPCFlowLogExtractor is an extractor that receives JSON serialised VPCFlowLogBatch objects
// and extracts VPCFlowLog records from them
type VPCFlowLogExtractor struct {
}

// NewVPCFlowLogExtractor creates a new VPCFlowLogExtractor
func NewVPCFlowLogExtractor() artifact_source.Extractor {
	return &VPCFlowLogExtractor{}
}

func (c *VPCFlowLogExtractor) Identifier() string {
	return "vpc_flow_log_extractor"
}

// Extract unmarshalls the artifact data as an VPCFlowLogBatch and returns the VPCFlowLog records
func (c *VPCFlowLogExtractor) Extract(_ context.Context, a any) ([]any, error) {
	// the expected input type is a JSON byte[] deserializable to VPCFlowLogBatch
	jsonBytes, ok := a.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected byte[], got %T", a)
	}

	mappedData, err := ConvertToMapSlice(string(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("error in mapping the results in the form of []map[string]string")
	}

	var res = make([]any, len(mappedData))
	for i, record := range mappedData {

		flowLog := &VpcFlowLog{}
		err := flowLog.MapValues(record)
		if err != nil {
			return nil, fmt.Errorf("error in extracting the log: %v", err)
		}
		res[i] = flowLog
	}
	return res, nil
}

// ConvertToMapSlice converts space-separated string data into a slice of map[string]string
func ConvertToMapSlice(data string) ([]map[string]string, error) {
	// Split the input into lines
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("insufficient data: needs at least a header and one row")
	}

	// Extract header row as keys
	keys := strings.Fields(lines[0]) // Splitting header by spaces

	var result []map[string]string

	// Iterate over the data rows
	for _, line := range lines[1:] {
		values := strings.Fields(line) // Splitting row by spaces
		entry := make(map[string]string)

		// Map values to keys
		for i, key := range keys {
			if i < len(values) {
				entry[key] = values[i]
			} else {
				entry[key] = "" // Assign empty string if missing
			}
		}
		result = append(result, entry)
	}

	return result, nil
}

// MapValues splits the input string and assigns each word to a corresponding key in the slice
func (flowLog *VpcFlowLog) MapValues(input map[string]string) error {

	for i, field := range input {
		// skip empty fields
		// https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs-records-examples.html#flow-log-example-no-data
		if field == "-" || field == "SKIPDATA" || field == "NODATA" {
			continue
		}
		switch i {
		case "version":
			version, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid version: %s", field)
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
				return fmt.Errorf("invalid srcport: %s", field)
			}
			v := int32(srcPort)
			flowLog.SrcPort = &v
		case "dstport":
			dstPort, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid dstport: %s", field)
			}
			v := int32(dstPort)
			flowLog.DstPort = &v
		case "protocol":
			protocol, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid protocol: %s", field)
			}
			v := int32(protocol)
			flowLog.Protocol = &v
		case "packets":
			packets, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid packets: %s", field)
			}
			flowLog.Packets = &packets
		case "bytes":
			bytes, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid bytes: %s", field)
			}
			flowLog.Bytes = &bytes
		case "start":
			start, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid start value: %s", field)
			}
			t := time.Unix(start, 0)
			flowLog.Start = &t
		case "end":
			end, err := strconv.ParseInt(field, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid end value: %s", field)
			}
			t := time.Unix(end, 0)
			flowLog.End = &t
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
				return fmt.Errorf("invalid tcp-flags: %s", field)
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
		case "reject-reason":
			flowLog.RejectReason = &field
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
				return fmt.Errorf("invalid traffic-path: %s", field)
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
			return fmt.Errorf("unknown field: %s", i)
		}
	}

	return nil
}
