package aws

import (
	"context"
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

func NewPlugin(_ context.Context) *Plugin {
	return &Plugin{}
}

func (t *Plugin) Collect(req *proto.CollectRequest) error {
	log.Println("[INFO] Collect")

	go t.doCollect(req)

	return nil
}

func (t *Plugin) doCollect(req *proto.CollectRequest) {
	//
	//"sources": {
	//	"my_file": {
	//		"name": "my_file",
	//			"plugin": "file",
	//			"config": {
	//			"path": "/Users/nathan/src/play-duckdb/2023/02",
	//				"extensions": []
	//}
	//}
	//},
	//"collections": {
	//"my_aws_log": {
	//"plugin": "aws_cloudtrail_log",
	//"name": "my_aws_log",
	//"source": "my_file"
	//}
	//}
	// tactical
	//create source and collection
	source := &AwsS3BucketSource{}
	collection := &AwsCloudTrailLogCollection{}

	// started hook
	t.OnStarted(req)

	// completion hook
	t.OnComplete(req, nil)
}
