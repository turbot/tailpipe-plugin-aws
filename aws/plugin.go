package aws

import (
	"context"
	"github.com/turbot/tailpipe-plugin-aws/aws_collection"
	"github.com/turbot/tailpipe-plugin-aws/aws_source"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"log"
)

type Plugin struct {
	plugin.Base
}

func (t *Plugin) Identifier() string {
	return "aws"
}

func (t *Plugin) Collect(req *proto.CollectRequest) error {
	log.Println("[INFO] Collect")

	go t.doCollect(context.Background(), req)

	return nil
}

//// GetSchema returns the schema (i.e. an instance of the row struct) for all collections
//// it is used primarily to validate the row structs provide the required fields
//func (t *Plugin) GetSchema(collection string) map[string]any {
//	return map[string]any{
//		aws_collection.CloudTrailLogCollection{}.Identifier(): aws_types.AWSCloudTrail{},
//	}
//}

func (t *Plugin) doCollect(ctx context.Context, req *proto.CollectRequest) {
	// todo config parsing, identify collection type etc.

	// TODO parse config and use to build collection
	//  tactical - create collection
	sourceConfig := aws_source.CompressedFileSourceConfig{Paths: []string{"/Users/kai/Downloads/flaws_cloudtrail_logs"}}
	var source = aws_source.NewCompressedFileSourceConfig(sourceConfig)

	collectionConfig := aws_collection.CloudTrailLogCollectionConfig{}
	var col = aws_collection.NewCloudTrailLogCollection(collectionConfig, source)

	// add ourselves as an observer
	col.AddObserver(t)

	// signal we have started
	t.OnStarted(req)

	// tell the collection to start collecting - this is a blocking call
	err := col.Collect(ctx, req)

	// signal we have completed
	t.OnComplete(req, err)
}
