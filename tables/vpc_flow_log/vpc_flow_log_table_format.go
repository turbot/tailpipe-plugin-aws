package vpc_flow_log

var DefaultFlowLogFields = []string{
	"version",
	"account-id",
	"interface-id",
	"srcaddr",
	"dstaddr",
	"srcport",
	"dstport",
	"protocol",
	"packets",
	"bytes",
	"start",
	"end",
	"action",
	"log-status",
	"ecs-cluster-name",
	"ecs-cluster-arn",
	"ecs-container-instance-id",
	"ecs-container-instance-arn",
	"ecs-service-name",
	"ecs-task-definition-arn",
	"ecs-task-id",
	"ecs-task-arn",
	"ecs-container-id",
	"ecs-second-container-id",
}

type VpcFlowLogTableFormat struct {
	// the fields to extract from the flow log
	Fields []string `hcl:"fields,optional"`
}

func (c *VpcFlowLogTableFormat) Validate() error {
	// set default fields if none are specified
	if len(c.Fields) == 0 {
		c.Fields = DefaultFlowLogFields
	}

	return nil
}

func (*VpcFlowLogTableFormat) Identifier() string {
	return VpcFlowLogTableIdentifier
}