package state

import (
	"errors"
	"github.com/ethereum/go-ethereum/ethdb"
	"reflect"
)

// Arbitrum: ArbDB provider setter.
func (db *cachingDB) SetArbDB(kv ethdb.KeyValueWriter) {
	db.arbDB = kv
}

// Arbitrum: ArbDB provider getter.
func (db *cachingDB) ArbDB() (ethdb.KeyValueWriter, error) {
	return loadArbDB(db.arbDB)
}

func loadArbDB(db ethdb.KeyValueWriter) (ethdb.KeyValueWriter, error) {
	if db == nil || reflect.ValueOf(db).IsNil() {
		return nil, errors.New("nil arbDB retrieved")
	}
	return db, nil
}
