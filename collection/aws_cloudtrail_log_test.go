package collection

import (
	"testing"

	"github.com/turbot/tailpipe-plugin-sdk/collection"
)

func TestConformance(t *testing.T) {
	collection.RunConformanceTests(t, &AwsCloudTrailLogCollection{})
}
