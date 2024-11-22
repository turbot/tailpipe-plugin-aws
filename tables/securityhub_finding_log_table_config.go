package tables

type SecurityHubFindingLogTableConfig struct {
}

func (a SecurityHubFindingLogTableConfig) Validate() error {
	// TODO: #config validate the config
	return nil
}

func (SecurityHubFindingLogTableConfig) Identifier() string {
	return SecurityHubFindingLogTableIdentifier
}
