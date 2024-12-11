package tables

import "fmt"

type CloudTrailLogTableConfig struct{}

func (c *CloudTrailLogTableConfig) Identifier() string {
	return CloudTrailLogTableIdentifier
}

func (c *CloudTrailLogTableConfig) Validate() error {
	return fmt.Errorf("not implemented")
}
