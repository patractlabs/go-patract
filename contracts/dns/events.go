package dns

import (
	"fmt"

	"github.com/patractlabs/go-patract/types"
)

type EventRegister struct {
	Name types.Hash            `scale:"name"`
	From types.OptionAccountID `scale:"from"`
}

func (e EventRegister) String() string {
	return fmt.Sprintf("event Register: %s - %s", e.Name, e.From)
}

type EventSetAddress struct {
	Name       types.Hash            `scale:"name"`
	From       types.AccountID       `scale:"from"`
	OldAddress types.OptionAccountID `scale:"old_address"`
	NewAddress types.AccountID       `scale:"new_address"`
}

func (e EventSetAddress) String() string {
	return fmt.Sprintf("event SetAddress: %s : %s - %s -> %s",
		e.Name, types.NewOptionAccountID(e.From), e.OldAddress, types.NewOptionAccountID(e.NewAddress))
}

type EventTransfer struct {
	Name     types.Hash            `scale:"name"`
	From     types.AccountID       `scale:"from"`
	OldOwner types.OptionAccountID `scale:"old_owner"`
	NewOwner types.AccountID       `scale:"new_owner"`
}

func (e EventTransfer) String() string {
	return fmt.Sprintf("event Transfer: %s : %s - %s -> %s",
		e.Name, types.NewOptionAccountID(e.From), e.OldOwner, types.NewOptionAccountID(e.NewOwner))
}
