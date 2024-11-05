package aws

import (
	"log/slog"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-aws/tables"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

type Plugin struct {
	plugin.PluginImpl
}

func NewPlugin() (_ plugin.TailpipePlugin, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = helpers.ToError(r)
		}
	}()
	slog.Info("AWS Plugin starting")
	//time.Sleep(10 * time.Second)
	slog.Info("AWS Plugin started")

	p := &Plugin{
		PluginImpl: plugin.NewPluginImpl("aws", config.NewAwsConnection),
	}

	// register the tables that we provide
	resources := &plugin.ResourceFunctions{
		Tables: []func() table.Table{
			tables.NewCloudTrailLogTable,
			tables.NewVPCFlowLogLogTable,
			tables.NewElbAccessLogTable,
			tables.NewS3ServerAccessLogTable,
			tables.NewLambdaLogTable,
		},
	}

	if err := p.RegisterResources(resources); err != nil {
		return nil, err
	}

	return p, nil
}
