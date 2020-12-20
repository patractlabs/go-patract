package erc20

import "github.com/patractlabs/go-patract/types"

type EventTransfer struct {
	From  types.OptionAccountID `scale:"from"`
	To    types.OptionAccountID `scale:"to"`
	Value types.U128            `scale:"value"`
}

type EventApproval struct {
	From  types.AccountID `scale:"from"`
	To    types.AccountID `scale:"to"`
	Value types.U128      `scale:"value"`
}
