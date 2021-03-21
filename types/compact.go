package types

import (
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
)

// Some type by Compact

type (
	// UCompact for big int
	UCompact = types.UCompact
	// CompactGas Compact<Gas> type
	CompactGas = types.UCompact
	// CompactBalance Compact<Balance> type
	CompactBalance = types.UCompact
)

var (
	// NewUCompactFromUInt create a compact big int from uint64
	NewUCompactFromUInt = types.NewUCompactFromUInt
)

// NewCompactGas creates a new compact of gas
func NewCompactGas(gas Gas) types.UCompact {
	return types.NewUCompactFromUInt(uint64(gas))
}

// NewCompactBalance create a new CompactBalance
func NewCompactBalance(amount uint64) types.UCompact {
	return types.NewUCompactFromUInt(amount)
}

// NewCompactBalanceByInt create a new CompactBalance
func NewCompactBalanceByInt(amount *big.Int) types.UCompact {
	return types.NewUCompact(amount)
}
