package aws

import (
	"log/slog"
	"time"

	"github.com/turbot/tailpipe-plugin-aws/aws_collection"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

type Plugin struct {
	plugin.PluginBase
}

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := &Plugin{}

	slog.Info("AWS Plugin starting")
	time.Sleep(10 * time.Second)
	slog.Info("YAWN")

	// register the collections that we provide
	resources := &plugin.ResourceFunctions{
		Collections: []func() collection.Collection{aws_collection.NewCloudTrailLogCollection, aws_collection.NewVPCFlowLogLogCollection},
	}

	if err := p.RegisterResources(resources); err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "aws"
}
