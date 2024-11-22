package tables

type WafTrafficLogTableConfig struct{}

func (c WafTrafficLogTableConfig) Validate() error {
	return nil
}

func (WafTrafficLogTableConfig) Identifier() string {
	return WaftTrafficLogTableIdentifier
}
