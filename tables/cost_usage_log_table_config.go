package tables

type CostAndUsageLogTableConfig struct{}

func (c *CostAndUsageLogTableConfig) Validate() error {
	return nil
}

func (*CostAndUsageLogTableConfig) Identifier() string {
	return CostUsageLogTableIdentifier
}
