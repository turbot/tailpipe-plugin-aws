package aws

import (
	"log/slog"
	"time"

	"github.com/turbot/tailpipe-plugin-aws/aws_collection"
	"github.com/turbot/tailpipe-plugin-aws/aws_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_mapper"
	"github.com/turbot/tailpipe-plugin-sdk/collection"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
)

type Plugin struct {
	plugin.Base
}

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := &Plugin{}

	slog.Info("AWS Plugin starting")
	time.Sleep(10 * time.Second)
	slog.Info("YAWN")

	// register the collections, sources and mappers that we provide
	resources := &plugin.ResourceFunctions{
		Collections:     []func() collection.Collection{aws_collection.NewCloudTrailLogCollection, aws_collection.NewVPCFlowLogLogCollection},
		Sources:         []func() row_source.RowSource{aws_source.NewCloudwatchSource},
		ArtifactMappers: []func() artifact_mapper.Mapper{aws_source.NewCloudtrailMapper},
	}

	if err := p.RegisterResources(resources); err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "aws"
}
