package tables

type SecurityHubFindingTableConfig struct {
}

func (a SecurityHubFindingTableConfig) Validate() error {
	// TODO: #config validate the config
	return nil
}

func (SecurityHubFindingTableConfig) Identifier() string {
	return SecurityHubFindingTableIdentifier
}
