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

type FlowLogCollectionConfig struct {
	// the path to the flow log files
	Paths []string `hcl:"paths"`

	// the fields to extract from the flow log
	Fields []string `hcl "fields"`
}
