package aws_collection

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
}

type VpcFlowLogCollectionConfig struct {
	// the fields to extract from the flow log
	Fields []string `hcl:"fields,optional"`
}

func (c *VpcFlowLogCollectionConfig) Validate() error {
	// set default fields if none are specified
	if len(c.Fields) == 0 {
		c.Fields = DefaultFlowLogFields
	}

	return nil
}
