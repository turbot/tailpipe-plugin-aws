package main

import (
	"github.com/turbot/tailpipe-plugin-aws/aws"
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Plugin: &aws.Plugin{},
	})
}
