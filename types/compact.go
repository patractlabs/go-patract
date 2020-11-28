package types

import "github.com/centrifuge/go-substrate-rpc-client/types"

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
