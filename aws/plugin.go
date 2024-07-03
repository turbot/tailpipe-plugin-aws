package aws

import (
	"context"
	"github.com/turbot/tailpipe-plugin-aws/aws_collection"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"log"
)

type Plugin struct {
	plugin.Base
}

func NewPlugin() (plugin.TailpipePlugin, error) {
	p := &Plugin{}

	//time.Sleep(10 * time.Second)
	// register collections which we support
	p.RegisterCollections(aws_collection.NewCloudTrailLogCollection)

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
	// todo config parsing, identify collection type etc.

	// TODO parse config and use to build collection
	//  tactical - create collection

	collectionConfig := aws_collection.CloudTrailLogCollectionConfig{}
	var col = aws_collection.NewCloudTrailLogCollection()
	// TEMP call init
	col.Init(collectionConfig)

	// add ourselves as an observer
	col.AddObserver(t)

	// signal we have started
	t.OnStarted(req)

	// tell the collection to start collecting - this is a blocking call
	err := col.Collect(ctx, req)

	// signal we have completed
	t.OnComplete(req, err)
}
