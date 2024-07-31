package aws_source

import (
	"context"
	"github.com/turbot/tailpipe-plugin-sdk/artifact_row_source"
	"github.com/turbot/tailpipe-plugin-sdk/hcl"
	"github.com/turbot/tailpipe-plugin-sdk/row_source"
)

const CloudwatchSourceIdentifier = "aws_cloudwatch_source"

type CloudwatchSource struct {
	artifact_row_source.Base
}

func NewCloudwatchSource() row_source.RowSource {
	return &CloudwatchSource{}
}

// Identifier returns the name of the row source
func (c *CloudwatchSource) Identifier() string {
	return CloudwatchSourceIdentifier
}

func (c *CloudwatchSource) Init(ctx context.Context, configData *hcl.Data, opts ...row_source.RowSourceOption) error {
	return c.Base.Init(ctx, configData, opts...)
}
