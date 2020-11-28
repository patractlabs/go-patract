package types

import "math/big"

type (
	// Balance balance type to chain, which is BalanceIf<T>
	Balance = U128

	// CodeHash the code hash if code
	CodeHash = Hash
)

// NewBalanceByU64 new balance from uint64
func NewBalanceByU64(amount uint64) U128 {
	var i big.Int
	i.SetUint64(amount)
	return NewU128(i)
}
