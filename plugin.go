package main

import (
	"github.com/turbot/tailpipe-plugin-aws/tables"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/table"
)

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := plugin.NewPlugin("aws")

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
