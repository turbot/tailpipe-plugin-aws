package vpc_flow_log

import (
	"fmt"
	"regexp"
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
	return nil
}

// Identifier returns the format TYPE
func (a *VPCFlowLogTableFormat) Identifier() string {
	// format name is same as table name
	return VpcFlowLogTableIdentifier
}

// GetName returns the format instance name
func (a *VPCFlowLogTableFormat) GetName() string {
	// format name is same as table name
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
	format := regexp.QuoteMeta(a.Layout)
	var unsupportedTokens []string

	// regex to grab tokens
	re := regexp.MustCompile(`\\\$\w+`)

	// check for concatenated tokens (e.g. $body_bytes$status)
	tokens := re.FindAllStringIndex(format, -1)
	for i := 1; i < len(tokens); i++ {
		// With QuoteMeta, tokens will be 2 characters further apart due to the backslash escape
		if tokens[i][0]-tokens[i-1][1] < 1 {
			return "", fmt.Errorf("concatenated tokens detected in format '%s', this is currently unsupported in this format, if this is a requirement a Regex format can be used", a.Layout)
		}
	}

	// replace tokens with regex patterns
	format = re.ReplaceAllStringFunc(format, func(match string) string {
		if pattern, exists := getRegexForSegment(match); exists {
			return pattern
		} else {
			unsupportedTokens = append(unsupportedTokens, strings.TrimPrefix(match, `\`))
		}

		return match
	})

	if len(unsupportedTokens) > 0 {
		return "", fmt.Errorf("the following tokens are not currently supported in this format: %s", strings.Join(unsupportedTokens, ", "))
	}

	if len(format) > 0 {
		format = fmt.Sprintf("^%s", format)
	}

	return format, nil
}

func (a *VPCFlowLogTableFormat) GetProperties() map[string]string {
	return map[string]string{
		"layout": a.Layout,
	}
}

func getRegexForSegment(segment string) (string, bool) {
	const defaultRegexFormat = `(?P<%s>[^ ]*)`

	if _, exists := getValidTokenMap()[segment]; !exists {
		return segment, false
	}

	return fmt.Sprintf(defaultRegexFormat, strings.TrimPrefix(segment, `\$`)), true
}

func getValidTokenMap() map[string]struct{} {
	return map[string]struct{}{
		`\$account-id`:                 {},
		`\$action`:                     {},
		`\$az-id`:                      {},
		`\$bytes`:                      {},
		`\$dstaddr`:                    {},
		`\$dstport`:                    {},
		`\$end`:                        {},
		`\$flow-direction`:             {},
		`\$instance-id`:                {},
		`\$interface-id`:               {},
		`\$log-status`:                 {},
		`\$packets`:                    {},
		`\$pkt-dst-aws-service`:        {},
		`\$pkt-dstaddr`:                {},
		`\$pkt-src-aws-service`:        {},
		`\$pkt-srcaddr`:                {},
		`\$protocol`:                   {},
		`\$region`:                     {},
		`\$reject-reason`:              {},
		`\$srcaddr`:                    {},
		`\$srcport`:                    {},
		`\$start`:                      {},
		`\$sublocation-id`:             {},
		`\$sublocation-type`:           {},
		`\$subnet-id`:                  {},
		`\$tcp-flags`:                  {},
		`\$traffic-path`:               {},
		`\$type`:                       {},
		`\$version`:                    {},
		`\$vpc-id`:                     {},
		`\$ecs-cluster-name`:           {},
		`\$ecs-cluster-arn`:            {},
		`\$ecs-container-instance-id`:  {},
		`\$ecs-container-instance-arn`: {},
		`\$ecs-service-name`:           {},
		`\$ecs-task-definition-arn`:    {},
		`\$ecs-task-id`:                {},
		`\$ecs-task-arn`:               {},
		`\$ecs-container-id`:           {},
		`\$ecs-second-container-id`:    {},
	}
}
