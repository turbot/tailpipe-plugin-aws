package aws

import (
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/tailpipe-plugin-aws/config"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/table"
	// reference the table package to ensure that the tables are registered by the init functions
	_ "github.com/turbot/tailpipe-plugin-aws/tables"
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

	//fmt.Println(tables.AlbAccessLogTable{})
	p := &Plugin{
		PluginImpl: plugin.NewPluginImpl("aws", config.NewAwsConnection),
	}

	// initialise table factory
	if err := table.Factory.Init(); err != nil {
		return nil, err
	}

	return p, nil
}
