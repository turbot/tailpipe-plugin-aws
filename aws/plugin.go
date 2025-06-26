package aws

import (
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/sources/cloudwatch_log_group"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-aws/tables/alb_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/alb_connection_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/clb_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/cloudtrail_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/cost_and_usage_focus"
	"github.com/turbot/tailpipe-plugin-aws/tables/cost_and_usage_report"
	"github.com/turbot/tailpipe-plugin-aws/tables/cost_optimization_recommendation"
	"github.com/turbot/tailpipe-plugin-aws/tables/guardduty_finding"
	"github.com/turbot/tailpipe-plugin-aws/tables/nlb_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/s3_server_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/securityhub_finding"
	"github.com/turbot/tailpipe-plugin-aws/tables/vpc_flow_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/waf_traffic_log"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type Plugin struct {
	plugin.PluginImpl
}

func (p *Plugin) GetTables() []any {
	return []any{
		&alb_access_log.AlbAccessLogTable{},
		&alb_connection_log.AlbConnectionLogTable{},
		&clb_access_log.ClbAccessLogTable{},
		&cloudtrail_log.CloudTrailLogTable{},
		&cost_and_usage_focus.CostUsageFocusTable{},
		&cost_and_usage_report.CostUsageReportTable{},
		&cost_optimization_recommendation.CostOptimizationRecommendationsTable{},
		&guardduty_finding.GuardDutyFindingTable{},
		&nlb_access_log.NlbAccessLogTable{},
		&s3_server_access_log.S3ServerAccessLogTable{},
		&securityhub_finding.SecurityHubFindingTable{},
		&vpc_flow_log.VpcFlowLogTable{},
		&waf_traffic_log.WafTrafficLogTable{},
	}
}

func (p *Plugin) GetSources() []row_source.RowSource {
	return []row_source.RowSource{
		&cloudwatch_log_group.AwsCloudWatchLogGroupSource{},
		&s3_bucket.AwsS3BucketSource{},
	}
}

func (p *Plugin) GetConnectionConfig() *config.AwsConnection {
	return &config.AwsConnection{}
}

func (p *Plugin) GetName() string {
	return "aws"
}

func (p *Plugin) GetVersion() string {
	return "0.1.0"
}

func init() {
	// Register tables with their row types
	table.RegisterTable[cloudtrail_log.CloudTrailLog, cloudtrail_log.CloudTrailLogTable]()
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
