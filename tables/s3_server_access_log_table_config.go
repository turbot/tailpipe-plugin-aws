package tables

type S3ServerAccessLogTableConfig struct {
}

func (a S3ServerAccessLogTableConfig) Validate() error {
	// TODO: #config validate the config
	return nil
}

func (S3ServerAccessLogTableConfig) Identifier() string {
	return S3ServerAccessLogTableIdentifier
}
