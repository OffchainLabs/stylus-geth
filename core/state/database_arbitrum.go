package state

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

var ErrNotFound = errors.New("not found")

// CompiledWasmContractCode retrieves a particular contract's compiled wasm code.
func (db *cachingDB) CompiledWasmContractCode(version uint32, codeHash common.Hash) ([]byte, error) {
	wasmKey := rawdb.CompiledWasmCodeKey(version, codeHash)
	if code := db.compiledWasmCache.Get(nil, wasmKey); len(code) > 0 {
		return code, nil
	}
	code, err := db.db.DiskDB().Get(wasmKey)
	if err != nil {
		if strings.Contains("not found", err.Error()) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if len(code) > 0 {
		db.compiledWasmCache.Set(wasmKey, code)
		return code, nil
	}
	return nil, ErrNotFound
}
