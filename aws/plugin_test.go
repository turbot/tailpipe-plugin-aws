package aws

import (
	"testing"

	"github.com/turbot/tailpipe-plugin-sdk/plugin"
)

func TestConformance(t *testing.T) {
	plugin.Validate(t, NewPlugin)
}
