package vpc_flow_log

import sdkformats "github.com/turbot/tailpipe-plugin-sdk/formats"

var defaultVPCFlowLogTableFormat = &VPCFlowLogTableFormat{
	Name:        "default",
	Description: "The default format for an VPC Flow Log.",
	Layout:      `version account-id interface-id srcaddr dstaddr srcport dstport protocol packets bytes start end action log-status`,
}

var VPCFlowLogTableFormatPresets = []sdkformats.Format{
	defaultVPCFlowLogTableFormat,
}
