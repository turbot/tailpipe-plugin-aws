package main

import (
	"log/slog"

	"github.com/turbot/tailpipe-plugin-aws/aws"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

func main() {
	err := plugin.Serve(&plugin.ServeOpts{
		PluginFunc: aws.NewPlugin,
	})

	if err != nil {
		slog.Error("Error starting aws", "error", err)
	}
}
