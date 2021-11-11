package rand_extension

import (
	"fmt"

	"github.com/patractlabs/go-patract/types"
)

type TransferSingle struct {
	Operator types.OptionAccountID `scale:"operator"`
	From     types.OptionAccountID `scale:"from"`
	To       types.OptionAccountID `scale:"to"`
	Id       types.U32             `scale:"id"`
}

func (e TransferSingle) String() string {
	return fmt.Sprintf("transfer Single: %s - %s -> %s by %s", e.Operator, e.From, e.To, e.Id)
}

type ApprovalForAll struct {
	Owner    types.AccountID `scale:"owner"`
	Operator types.AccountID `scale:"operator"`
	approved types.Bool      `scale:"approved"` // TODO: did not have #[ink(topic)]
}

func (e ApprovalForAll) String() string {
	return fmt.Sprintf("approval For All: %s -> %s by %s",
		types.NewOptionAccountID(e.Owner), types.NewOptionAccountID(e.Operator), e.approved)
}

type Uri struct {
	value   string     `scale:"value"` // TODO: did not have #[ink(topic)]
	TokenId types.U128 `scale:"token_id"`
}

func (e Uri) String() string {
	return fmt.Sprintf("Uri: %s : %s", e.TokenId, e.value)
}
