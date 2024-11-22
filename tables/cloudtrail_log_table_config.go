package tables

type CloudTrailLogTableConfig struct{}

func (c *CloudTrailLogTableConfig) Identifier() string {
	return CloudTrailLogTableIdentifier
}

func (c *CloudTrailLogTableConfig) Validate() error {
	return nil
}
