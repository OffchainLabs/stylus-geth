package light

import (
	"errors"
	"github.com/ethereum/go-ethereum/ethdb"
)

// Arbitrum: ArbDB provider setter.
func (db *odrDatabase) SetArbDB(kv ethdb.KeyValueWriter) {
}

// Arbitrum: ArbDB provider getter.
func (db *odrDatabase) ArbDB() (ethdb.KeyValueWriter, error) {
	return nil, errors.New("arbDB unsupported for light trie database")
}
