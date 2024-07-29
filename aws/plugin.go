package aws

import (
	"github.com/turbot/tailpipe-plugin-aws/aws_collection"
	"github.com/turbot/tailpipe-plugin-aws/aws_source"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
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

	// register collections which we support
	err := p.RegisterCollections(aws_collection.NewCloudTrailLogCollection, aws_collection.NewVPCFlowLogLogCollection)
	if err != nil {
		return nil, err
	}

	// register artifact mappers
	err = p.RegisterArtifactMappers(aws_source.NewCloudtrailMapper)

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "aws"
}
