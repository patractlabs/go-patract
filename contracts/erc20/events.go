package erc20

import (
	"fmt"

	"github.com/patractlabs/go-patract/types"
)

type EventTransfer struct {
	From  types.OptionAccountID `scale:"from"`
	To    types.OptionAccountID `scale:"to"`
	Value types.U128            `scale:"value"`
}

func (e EventTransfer) String() string {
	return fmt.Sprintf("event Transfer: %s -> %s by %s", e.From, e.To, e.Value)
}

type EventApproval struct {
	From  types.AccountID `scale:"from"`
	To    types.AccountID `scale:"to"`
	Value types.U128      `scale:"value"`
}

func (e EventApproval) String() string {
	return fmt.Sprintf("event Approval: %s -> %s by %s",
		types.NewOptionAccountID(e.From), types.NewOptionAccountID(e.To), e.Value)
}
