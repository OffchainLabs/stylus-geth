// Arbitrum: Storing the L2 non-consensus DB as a field within the chain database.
type ArbDBProvider interface {
	SetArbDB(db KeyValueWriter) error
	ArbDB() (KeyValueWriter, error)
}
