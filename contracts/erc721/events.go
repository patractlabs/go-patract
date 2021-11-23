package erc721

import (
	"fmt"

	"github.com/patractlabs/go-patract/types"
)

type EventTransfer struct {
	From types.OptionAccountID `scale:"from"`
	To   types.OptionAccountID `scale:"to"`
	Id   TokenId               `scale:"id"`
}

func (e EventTransfer) String() string {
	return fmt.Sprintf("event Transfer: %s -> %s by %v", e.From, e.To, e.Id)
}

type EventApproval struct {
	From types.AccountID `scale:"from"`
	To   types.AccountID `scale:"to"`
	Id   TokenId         `scale:"id"`
}

func (e EventApproval) String() string {
	return fmt.Sprintf("event Approval: %s -> %s by %v",
		types.NewOptionAccountID(e.From), types.NewOptionAccountID(e.To), e.Id)
}

type EventApprovalForAll struct {
	Owner    types.AccountID `scale:"owner"`
	Operator types.AccountID `scale:"operator"`
	approved Bool            `scale:"approved"` // TODO: did not have #[ink(topic)]
}

func (e EventApprovalForAll) String() string {
	return fmt.Sprintf("event Approval For All: %s -> %s by %v",
		types.NewOptionAccountID(e.Owner), types.NewOptionAccountID(e.Operator), e.approved)
}
