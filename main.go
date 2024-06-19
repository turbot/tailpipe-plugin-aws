package main

import (
	"context"
	"github.com/turbot/tailpipe-plugin-aws/aws"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		// TODO should we pass func, not object? For dynamic plugins? Will we have those?
		Plugin: aws.NewPlugin(context.Background()),
	})
}
