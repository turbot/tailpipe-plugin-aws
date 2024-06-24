package aws

import (
	"context"
	"github.com/turbot/tailpipe-plugin-sdk/grpc/proto"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"log"
	"log/slog"
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

func (t *Plugin) doCollect(ctx context.Context, req *proto.CollectRequest) {
	onRow := func(row any) {
		t.OnRow(row, req)
	}

	//create collection
	collection := &AwsCloudTrailLogCollection{
		paths: []string{"/Users/kai/Downloads/flaws_cloudtrail_logs"},
	}

	// started hook
	t.OnStarted(req)

	err := collection.Collect(ctx, onRow)

	// completion hook
	if err := t.OnComplete(req, err); err != nil {
		// TODO #error
		slog.Error("error notifying observers of completion", "error", err)
	}
}
