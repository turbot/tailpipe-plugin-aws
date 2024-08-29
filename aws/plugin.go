package aws

import (
	"github.com/turbot/tailpipe-plugin-aws/aws_table"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type Plugin struct {
	plugin.PluginBase
}

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := &Plugin{}

	//slog.Info("AWS Plugin starting")
	//time.Sleep(10 * time.Second)
	//slog.Info("YAWN")

	// register the tables that we provide
	resources := &plugin.ResourceFunctions{
		Tables: []func() table.Table{
			aws_table.NewCloudTrailLogTable,
			aws_table.NewVPCFlowLogLogTable,
			aws_table.NewElbAccessLogTable,
			aws_table.NewS3ServerAccessLogTable,
			aws_table.NewLambdaLogTable,
		},
	}

	if err := p.RegisterResources(resources); err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "aws"
}
