package aws

import (
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/sources/cloudwatch_log_group"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-aws/tables/alb_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/clb_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/cloudtrail_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/cost_and_usage_focus"
	"github.com/turbot/tailpipe-plugin-aws/tables/cost_and_usage_report"
	"github.com/turbot/tailpipe-plugin-aws/tables/cost_optimization_recommendation"
	"github.com/turbot/tailpipe-plugin-aws/tables/guardduty_finding"
	"github.com/turbot/tailpipe-plugin-aws/tables/nlb_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/s3_server_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/vpc_flow_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/waf_traffic_log"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type Plugin struct {
	plugin.PluginImpl
}

func init() {
	// Register tables, with type parameters:
	// 1. row struct
	// 2. table implementation
	table.RegisterTable[*alb_access_log.AlbAccessLog, *alb_access_log.AlbAccessLogTable]()
	table.RegisterTable[*clb_access_log.ClbAccessLog, *clb_access_log.ClbAccessLogTable]()
	table.RegisterTable[*cloudtrail_log.CloudTrailLog, *cloudtrail_log.CloudTrailLogTable]()
	table.RegisterTable[*cost_and_usage_focus.CostUsageFocus, *cost_and_usage_focus.CostUsageFocusTable]()
	table.RegisterTable[*cost_and_usage_report.CostUsageReport, *cost_and_usage_report.CostUsageReportTable]()
	table.RegisterTable[*cost_optimization_recommendation.CostOptimizationRecommendation, *cost_optimization_recommendation.CostOptimizationRecommendationsTable]()
	table.RegisterTable[*guardduty_finding.GuardDutyFinding, *guardduty_finding.GuardDutyFindingTable]()
	table.RegisterTable[*nlb_access_log.NlbAccessLog, *nlb_access_log.NlbAccessLogTable]()
	table.RegisterTable[*s3_server_access_log.S3ServerAccessLog, *s3_server_access_log.S3ServerAccessLogTable]()
	table.RegisterTable[*waf_traffic_log.WafTrafficLog, *waf_traffic_log.WafTrafficLogTable]()

	// register custom table
	table.RegisterCustomTable[*vpc_flow_log.VpcFlowLogTable]()

	// register sources
	row_source.RegisterRowSource[*s3_bucket.AwsS3BucketSource]()
	row_source.RegisterRowSource[*cloudwatch_log_group.AwsCloudWatchLogGroupSource]()

	// register formats
	table.RegisterFormatPresets(vpc_flow_log.VPCFlowLogTableFormatPresets...)
	table.RegisterFormat[*vpc_flow_log.VPCFlowLogTableFormat]()
	row_source.RegisterRowSource[*cloudwatch_log_group.AwsCloudWatchLogGroupSource]()
}

func NewPlugin() (_ plugin.TailpipePlugin, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = helpers.ToError(r)
		}
	}()

	p := &Plugin{
		PluginImpl: plugin.NewPluginImpl(config.PluginName),
	}

	return p, nil
}
