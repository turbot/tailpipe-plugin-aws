package aws

import (
	"github.com/turbot/tailpipe-plugin-aws/aws_collection"
	"github.com/turbot/tailpipe-plugin-aws/aws_source"
	"github.com/turbot/tailpipe-plugin-sdk/artifact"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
	"log/slog"
	"time"
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
	collections := []func() plugin.Collection{aws_collection.NewCloudTrailLogCollection, aws_collection.NewVPCFlowLogLogCollection}
	sources := []func() row_source.RowSource{aws_source.NewCloudwatchSource}
	mappers := []func() artifact.Mapper{aws_source.NewCloudtrailMapper}

	// register collections which we support
	if err := p.RegisterCollections(collections...); err != nil {
		return nil, err
	}
	// register sources
	if err := p.RegisterSources(sources...); err != nil {
		return nil, err
	}

	// register artifact mappers
	if err := p.RegisterArtifactMappers(mappers...); err != nil {
		return nil, err
	}

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "aws"
}
