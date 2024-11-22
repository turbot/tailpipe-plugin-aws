package tables

type GuardDutyFindingTableConfig struct{}

func (c GuardDutyFindingTableConfig) Validate() error {
	return nil
}

func (GuardDutyFindingTableConfig) Identifier() string {
	return GuardDutyFindingTableIdentifier
}
