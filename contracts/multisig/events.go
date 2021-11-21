package multisig

import (
	"fmt"

	"github.com/patractlabs/go-patract/types"
)

type EventConfirmation struct {
	Transaction TransactionId         `scale:"transaction"`
	From        types.OptionAccountID `scale:"from"`
	// TODO: status should
}

func (e EventConfirmation) String() string {
	return fmt.Sprintf("event Confirmation: %s - %s", e.Transaction, e.From)
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
