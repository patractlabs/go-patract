package dns

import (
	"fmt"

	"github.com/patractlabs/go-patract/types"
)

type Register struct {
	Name types.Hash            `scale:"name"`
	From types.OptionAccountID `scale:"from"`
}

func (e Register) String() string {
	return fmt.Sprintf("register: %s - %s", e.Name, e.From)
}

type SetAddress struct {
	Name       types.Hash            `scale:"name"`
	From       types.AccountID       `scale:"from"`
	OldAddress types.OptionAccountID `scale:"old_address"`
	NewAddress types.AccountID       `scale:"new_address"`
}

func (e SetAddress) String() string {
	return fmt.Sprintf("set Address: %s : %s - %s -> %s",
		e.Name, types.NewOptionAccountID(e.From), e.OldAddress, types.NewOptionAccountID(e.NewAddress))
}

type Transfer struct {
	Name     types.Hash            `scale:"name"`
	From     types.AccountID       `scale:"from"`
	OldOwner types.OptionAccountID `scale:"old_owner"`
	NewOwner types.AccountID       `scale:"new_owner"`
}

func (e Transfer) String() string {
	return fmt.Sprintf("transfer: %s : %s - %s -> %s",
		e.Name, types.NewOptionAccountID(e.From), e.OldOwner, types.NewOptionAccountID(e.NewOwner))
}
