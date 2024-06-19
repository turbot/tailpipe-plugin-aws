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

	// started hook
	t.OnStarted(req)

	// completion hook
	t.OnComplete(req, nil)
}
