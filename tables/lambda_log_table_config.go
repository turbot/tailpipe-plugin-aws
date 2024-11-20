package tables

type LambdaLogTableConfig struct {
}

func (a LambdaLogTableConfig) Validate() error {
	// TODO: #config validate the config
	return nil
}

func (LambdaLogTableConfig) Identifier() string {
	return LambdaLogTableIdentifier
}
