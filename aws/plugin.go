package aws

import (
	"github.com/turbot/tailpipe-plugin-aws/tables"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	"log/slog"
	"time"
)

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := plugin.NewPlugin("aws")

	slog.Info("AWS Plugin starting")
	time.Sleep(10 * time.Second)
	slog.Info("AWS Plugin started")
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
