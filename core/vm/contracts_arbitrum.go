package vm

import "github.com/ethereum/go-ethereum/common"

var (
	PrecompiledContractsArbitrum     = make(map[common.Address]PrecompiledContract)
	PrecompiledAddressesArbitrum     []common.Address
	ArbosVersionPrecompiledAddresses = make(map[uint64][]common.Address)
)
