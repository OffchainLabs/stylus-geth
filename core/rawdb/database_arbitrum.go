package rawdb

import (
	"errors"
	"github.com/ethereum/go-ethereum/ethdb"
	"reflect"
)

// Arbitrum: ArbDB provider setter.
func (t *table) SetArbDB(db ethdb.KeyValueWriter) {
	t.arbDB = db
}

// Arbitrum: ArbDB provider getter.
func (t *table) ArbDB() (ethdb.KeyValueWriter, error) {
	return loadArbDB(frdb.arbDB)
}

// Arbitrum: ArbDB provider setter.
func (db *nofreezedb) SetArbDB(arbDB ethdb.KeyValueWriter) {
	db.arbDB = arbDB
}

// Arbitrum: ArbDB provider getter.
func (db *nofreezedb) ArbDB() (ethdb.KeyValueWriter, error) {
	return loadArbDB(db.arbDB)
}

// Arbitrum: ArbDB provider setter.
func (frdb *freezerdb) SetArbDB(db ethdb.KeyValueWriter) {
	frdb.arbDB = db
}

// Arbitrum: ArbDB provider getter.
func (frdb *freezerdb) ArbDB() (ethdb.KeyValueWriter, error) {
	return loadArbDB(frdb.arbDB)
}

func loadArbDB(db ethdb.KeyValueWriter) (ethdb.KeyValueWriter, error) {
	if db == nil || reflect.ValueOf(db).IsNil() {
		return nil, errors.New("nil arbDB retrieved")
	}
	return db, nil
}
