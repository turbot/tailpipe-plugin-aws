package aws

import (
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/sources/s3_bucket"
	"github.com/turbot/tailpipe-plugin-aws/tables/alb_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/clb_access_log"
	"github.com/turbot/tailpipe-plugin-aws/tables/cloudtrail_log"
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
	table.RegisterTable[*cloudtrail_log.CloudTrailLog, *cloudtrail_log.CloudTrailLogTable]()
	table.RegisterTable[*nlb_access_log.NlbAccessLog, *nlb_access_log.NlbAccessLogTable]()
	table.RegisterTable[*s3_server_access_log.S3ServerAccessLog, *s3_server_access_log.S3ServerAccessLogTable]()
	table.RegisterTable[*vpc_flow_log.VpcFlowLog, *vpc_flow_log.VpcFlowLogTable]()
	table.RegisterTable[*waf_traffic_log.WafTrafficLog, *waf_traffic_log.WafTrafficLogTable]()
	table.RegisterTable[*clb_access_log.ClbAccessLog, *clb_access_log.ClbAccessLogTable]()

	// register sources
	row_source.RegisterRowSource[*s3_bucket.AwsS3BucketSource]()
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
