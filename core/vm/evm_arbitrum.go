// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"math/big"

	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	// Defines prefix bytes for Polyglot WASM program bytecode
	// when deployed on-chain via a user-initiated transaction.
	// These byte prefixes are meant to conflict with the L1 contract EOF
	// validation rules so they can be sufficiently differentiated from EVM bytecode.
	// This allows us to store WASM programs as code in the stateDB side-by-side
	// with EVM contracts, but match against these prefix bytes when loading code
	// to execute the WASMs through Polyglot rather than the EVM.
	polyglotEOFMagic         = byte(0xEF)
	polyglotEOFVersion       = byte(0x00)
	polyglotEOFSectionHeader = byte(0x00)
)

// Depth returns the current depth
func (evm *EVM) Depth() int {
	return evm.depth
}

func (evm *EVM) IncrementDepth() {
	evm.depth += 1
}

func (evm *EVM) DecrementDepth() {
	evm.depth -= 1
}

func (evm *EVM) AddArbDB(arbDB ethdb.Database) {
	if !evm.chainRules.IsArbitrum {
		return
	}
	evm.arbDB = arbDB
}

func (evm *EVM) ArbDB() (ethdb.Database, error) {
	if !evm.chainRules.IsArbitrum {
		return nil, errors.New("must be using Arbitrum EVM chain rules to fetch ArbDB from the EVM")
	}
	return evm.arbDB, nil
}

type TxProcessingHook interface {
	StartTxHook() (bool, uint64, error, []byte) // return 4-tuple rather than *struct to avoid an import cycle
	GasChargingHook(gasRemaining *uint64) error
	PushCaller(addr common.Address)
	PopCaller()
	ForceRefundGas() uint64
	NonrefundableGas() uint64
	EndTxHook(totalGasUsed uint64, evmSuccess bool)
	ScheduledTxes() types.Transactions
	L1BlockNumber(blockCtx BlockContext) (uint64, error)
	L1BlockHash(blockCtx BlockContext, l1BlocKNumber uint64) (common.Hash, error)
	GasPriceOp(evm *EVM) *big.Int
	FillReceiptInfo(receipt *types.Receipt)
}

type DefaultTxProcessor struct{}

func (p DefaultTxProcessor) StartTxHook() (bool, uint64, error, []byte) {
	return false, 0, nil, nil
}

func (p DefaultTxProcessor) GasChargingHook(gasRemaining *uint64) error {
	return nil
}

func (p DefaultTxProcessor) PushCaller(addr common.Address) {}

func (p DefaultTxProcessor) PopCaller() {
}

func (p DefaultTxProcessor) ForceRefundGas() uint64 {
	return 0
}

func (p DefaultTxProcessor) NonrefundableGas() uint64 {
	return 0
}

func (p DefaultTxProcessor) EndTxHook(totalGasUsed uint64, evmSuccess bool) {}

func (p DefaultTxProcessor) ScheduledTxes() types.Transactions {
	return types.Transactions{}
}

func (p DefaultTxProcessor) L1BlockNumber(blockCtx BlockContext) (uint64, error) {
	return blockCtx.BlockNumber.Uint64(), nil
}

func (p DefaultTxProcessor) L1BlockHash(blockCtx BlockContext, l1BlocKNumber uint64) (common.Hash, error) {
	return blockCtx.GetHash(l1BlocKNumber), nil
}

func (p DefaultTxProcessor) GasPriceOp(evm *EVM) *big.Int {
	return evm.GasPrice
}

func (p DefaultTxProcessor) FillReceiptInfo(*types.Receipt) {}

// Is PolyglotProgram checks if a specified bytecode is a user-submitted WASM program.
// Polyglot differentiates WASMs from EVM bytecode via the prefix 0xEF0000 which will safely fail
// to pass through EVM-bytecode EOF validation rules.
func IsPolyglotProgram(b []byte) bool {
	if len(b) < 3 {
		return false
	}
	return b[0] == polyglotEOFMagic && b[1] == polyglotEOFVersion && b[2] == polyglotEOFSectionHeader
}

// StripPolyglotPrefix if the specified input is a polyglot program.
func StripPolyglotPrefix(b []byte) ([]byte, error) {
	if !IsPolyglotProgram(b) {
		return nil, errors.New("specified bytecode is not a Polyglot program")
	}
	return b[3:], nil
}
