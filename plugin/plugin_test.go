package plugin

import (
	"testing"

	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

func TestConformance(t *testing.T) {
	plugin.RunConformanceTests(t, &AwsPlugin{})
}
