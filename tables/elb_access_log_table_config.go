package tables

type ElbAccessLogTableConfig struct {
}

func (a ElbAccessLogTableConfig) Validate() error {
	// TODO: #config validate the config
	return nil
}

func (ElbAccessLogTableConfig) Identifier() string {
	return ElbAccessLogTableIdentifier
}
