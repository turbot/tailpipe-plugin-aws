package aws

import (
	"github.com/turbot/tailpipe-plugin-sdk/plugin"
	"testing"
)

func TestConformance(t *testing.T) {
	plugin.Validate(t, NewPlugin)
}
