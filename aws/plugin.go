package aws

import (
	"context"
	"github.com/turbot/tailpipe-plugin-aws/aws_collection"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"log"
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
	p.RegisterCollections(aws_collection.NewCloudTrailLogCollection, aws_collection.NewFlowlogLogCollection)

	return p, nil
}

func (t *Plugin) Identifier() string {
	return "aws"
}

func (t *Plugin) Collect(req *proto.CollectRequest) error {
	log.Println("[INFO] Collect")

	// TODO can we do this in base
	go t.doCollect(context.Background(), req)

	return nil
}

func (t *Plugin) doCollect(ctx context.Context, req *proto.CollectRequest) {
	var col plugin.Collection
	var config any

	// TODO most of this can be done from base
	// todo config parsing, identify collection type etc.

	// TODO parse config and use to build collection

	switch req.CollectionName {
	case "aws_cloudtrail_log":
		col = aws_collection.NewCloudTrailLogCollection()
		config = &aws_collection.CloudTrailLogCollectionConfig{
			Paths: req.Paths,
		}
	case "aws_flow_log":
		col = aws_collection.NewFlowlogLogCollection()

		config = &aws_collection.FlowLogCollectionConfig{
			Paths: req.Paths,
			Fields: []string{"timestamp",
				"version",
				"account-id",
				"interface-id",
				"srcaddr",
				"dstaddr",
				"srcport",
				"dstport",
				"protocol",
				"packets",
				"bytes",
				"start",
				"end",
				"action",
				"log-status",
			},
		}
	}

	// TEMP call init
	col.Init(config)

	// add ourselves as an observer
	col.AddObserver(t)

	// signal we have started
	t.OnStarted(req)

	// tell the collection to start collecting - this is a blocking call
	err := col.Collect(ctx, req)

	// signal we have completed - pass error if there was one
	t.OnComplete(req, err)
}
