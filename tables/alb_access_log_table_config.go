package tables

type AlbAccessLogTableConfig struct {
}

func (c *AlbAccessLogTableConfig) Identifier() string {
	return AlbAccessLogTableIdentifier
}

func (c *AlbAccessLogTableConfig) Validate() error {
	// By default, use built-in ALB log format
	// Add any additional validation as needed
	return nil
}
