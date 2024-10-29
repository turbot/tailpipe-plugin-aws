package main

import (
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"log/slog"
)

func main() {
	err := plugin.Serve(&plugin.ServeOpts{
		PluginFunc: NewPlugin,
	})

	if err != nil {
		slog.Error("Error starting plugin", "error", err)
	}
}
