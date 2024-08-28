package aws

import (
	"github.com/turbot/tailpipe-plugin-aws/aws_partition"
	"github.com/turbot/tailpipe-plugin-sdk/partition"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

type Plugin struct {
	plugin.PluginBase
}

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := &Plugin{}

	//slog.Info("AWS Plugin starting")
	//time.Sleep(10 * time.Second)
	//slog.Info("YAWN")

	// register the partitions that we provide
	resources := &plugin.ResourceFunctions{
		Partitions: []func() partition.Partition{
			aws_partition.NewCloudTrailLogPartition,
			aws_partition.NewVPCFlowLogLogPartition,
			aws_partition.NewElbAccessLogPartition,
			aws_partition.NewS3ServerAccessLogPartition,
			aws_partition.NewLambdaLogPartition,
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
