package vpc_flow_log

import (
	"fmt"
	"strings"

	"github.com/turbot/tailpipe-plugin-sdk/formats"
	"github.com/turbot/tailpipe-plugin-sdk/mappers"
	"github.com/turbot/tailpipe-plugin-sdk/types"
)

type VPCFlowLogTableFormat struct {
	// the name of this format instance
	Name string `hcl:"name,label"`
	// Description of the format
	Description string `hcl:"description,optional"`
	// the layout of the log line
	Layout string `hcl:"layout"`
}

func NewVPCFlowLogTableFormat() formats.Format {
	return &VPCFlowLogTableFormat{}
}

func (a *VPCFlowLogTableFormat) Validate() error {
	var invalid []string
	layoutParts := strings.Fields(a.Layout)
	for _, part := range layoutParts {
		if _, exists := getValidTokensAndColumnNames()[part]; !exists {
			invalid = append(invalid, part)
		}
	}

	if len(invalid) > 0 {
		return fmt.Errorf("the following tokens are not valid: %s", strings.Join(invalid, ", "))
	}

	return nil
}

// Identifier returns the format TYPE
func (a *VPCFlowLogTableFormat) Identifier() string {
	// format name is same as table name
	return VpcFlowLogTableIdentifier
}

// GetName returns the format instance name
func (a *VPCFlowLogTableFormat) GetName() string {
	return a.Name
}

// SetName sets the name of this format instance
func (a *VPCFlowLogTableFormat) SetName(name string) {
	a.Name = name
}

func (a *VPCFlowLogTableFormat) GetDescription() string {
	return a.Description
}

func (a *VPCFlowLogTableFormat) GetMapper() (mappers.Mapper[*types.DynamicRow], error) {
	// convert the layout to a regex
	regex, err := a.GetRegex()

	if err != nil {
		return nil, err
	}
	return mappers.NewRegexMapper[*types.DynamicRow](regex)
}

func (a *VPCFlowLogTableFormat) GetRegex() (string, error) {
	// validate checks all tokens in layout are valid else returns error
	err := a.Validate()
	if err != nil {
		return "", err
	}

	tokens := strings.Fields(a.Layout)
	var segments []string

	// get the regex segment for each token
	for _, token := range tokens {
		segments = append(segments, getRegexForSegment(token))
	}

	// return the regex pattern (space separated)
	return strings.Join(segments, " "), nil
}

func (a *VPCFlowLogTableFormat) GetProperties() map[string]string {
	return map[string]string{
		"layout": a.Layout,
	}
}

func getRegexForSegment(segment string) string {
	const defaultRegexFormat = "(?P<%s>[^ ]*)"

	if columnName, exists := getValidTokensAndColumnNames()[segment]; exists {
		return fmt.Sprintf(defaultRegexFormat, columnName)
	}

	return segment
}

func getValidTokensAndColumnNames() map[string]string {
	return map[string]string{
		"version":                    "version",
		"account-id":                 "account_id",
		"interface-id":               "interface_id",
		"srcaddr":                    "src_addr",
		"dstaddr":                    "dst_addr",
		"srcport":                    "src_port",
		"dstport":                    "dst_port",
		"protocol":                   "protocol",
		"packets":                    "packets",
		"bytes":                      "bytes",
		"start":                      "start_time",
		"end":                        "end_time",
		"action":                     "action",
		"log-status":                 "log_status",
		"vpc-id":                     "vpc_id",
		"subnet-id":                  "subnet_id",
		"instance-id":                "instance_id",
		"tcp-flags":                  "tcp_flags",
		"type":                       "type",
		"pkt-srcaddr":                "pkt_src_addr",
		"pkt-dstaddr":                "pkt_dst_addr",
		"region":                     "region",
		"az-id":                      "az_id",
		"sublocation-type":           "sublocation_type",
		"sublocation-id":             "sublocation_id",
		"pkt-src-aws-service":        "pkt_src_aws_service",
		"pkt-dst-aws-service":        "pkt_dst_aws_service",
		"flow-direction":             "flow_direction",
		"traffic-path":               "traffic_path",
		"ecs-cluster-arn":            "ecs_cluster_arn",
		"ecs-cluster-name":           "ecs_cluster_name",
		"ecs-container-instance-arn": "ecs_container_instance_arn",
		"ecs-container-instance-id":  "ecs_container_instance_id",
		"ecs-container-id":           "ecs_container_id",
		"ecs-second-container-id":    "ecs_second_container_id",
		"ecs-service-name":           "ecs_service_name",
		"ecs-task-definition-arn":    "ecs_task_definition_arn",
		"ecs-task-arn":               "ecs_task_arn",
		"ecs-task-id":                "ecs_task_id",
		"reject-reason":              "reject_reason",

		// This is not present in the log but is added by AWS while exporting the log from CloudWatch to S3, the timestamp is the time when the log was exported.
		// Here it the sample line from the log:
		// 2025-02-25T12:25:04.000Z i-085c7a43a498c2f5d eni-0416a1c81c87ab9c9 - - - use1-az2 - 1740486335 - - - - - - 1740486304 subnet-027e9a6d4add894eb
		"exported-at": "exported_at",
	}
}
