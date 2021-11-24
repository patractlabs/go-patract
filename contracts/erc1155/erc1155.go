package erc1155

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Create(ctx Context, value Balance) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"create"},
		value,
	)
}

func (a *API) Mint(ctx Context, tokenId TokenId, value Balance) (Hash, error) {
	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"mint"},
		tokenId, value,
	)
}

func (a *API) SafeTransferFrom(ctx Context, from AccountID, to AccountID, tokenId TokenId, value Balance, data VecU8) (Hash, error) {
	fromParam := struct {
		Address AccountID
	}{
		Address: from,
	}

	toParam := struct {
		Address AccountID
	}{
		Address: to,
	}

	dataParam := struct {
		data VecU8
	}{
		data: data,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"Erc1155", "safe_transfer_from"},
		fromParam, toParam, tokenId, value, dataParam,
	)
}

func (a *API) SafeBatchTransferFrom(ctx Context, from AccountID, to AccountID, tokenIds VecTokenId, values VecBalance, data VecU8) (Hash, error) {
	fromParam := struct {
		Address AccountID
	}{
		Address: from,
	}

	toParam := struct {
		Address AccountID
	}{
		Address: to,
	}

	tokenIdsParam := struct {
		Id VecTokenId
	}{
		Id: tokenIds,
	}

	valuesParam := struct {
		value VecBalance
	}{
		value: values,
	}

	dataParam := struct {
		data VecU8
	}{
		data: data,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"Erc1155", "safe_batch_transfer_from"},
		fromParam, toParam, tokenIdsParam, valuesParam, dataParam,
	)
}

func (a *API) BalanceOf(ctx Context, owner AccountID, tokenId TokenId) (Balance, error) {
	ownerParam := struct {
		Address AccountID
	}{
		Address: owner,
	}

	var res Balance

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"Erc1155", "balance_of"},
		ownerParam, tokenId,
	)

	return res, err
}

func (a *API) BalanceOfBatch(ctx Context, owners VecAccountID, tokenIds VecTokenId) (VecBalance, error) {
	ownersParam := struct {
		Address VecAccountID
	}{
		Address: owners,
	}

	tokenIdsParam := struct {
		Id VecTokenId
	}{
		Id: tokenIds,
	}

	var res VecBalance

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"Erc1155", "balance_of_batch"},
		ownersParam, tokenIdsParam,
	)

	return res, err
}

func (a *API) SetApprovalForAll(ctx Context, operator AccountID, approved Bool) (Hash, error) {
	operatorParam := struct {
		Address AccountID
	}{
		Address: operator,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"Erc1155", "set_approval_for_all"},
		operatorParam, approved,
	)
}

func (a *API) IsApprovedForAll(ctx Context, owner, operator AccountID) (Bool, error) {
	ownerParam := struct {
		Address AccountID
	}{
		Address: owner,
	}

	operatorParam := struct {
		Id AccountID
	}{
		Id: operator,
	}

	var res Bool

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"Erc1155", "is_approved_for_all"},
		ownerParam, operatorParam,
	)

	return res, err
}

func (a *API) OnReceived(ctx Context, _operator, _from AccountID, _token_id TokenId, _value Balance, _data VecU8) (Hash, error) {
	operatorParam := struct {
		Address AccountID
	}{
		Address: _operator,
	}

	fromParam := struct {
		Address AccountID
	}{
		Address: _from,
	}

	dataParam := struct {
		Data VecU8
	}{
		Data: _data,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"Erc1155TokenReceiver", "on_received"},
		operatorParam, fromParam, _token_id, _value, dataParam,
	)
}

func (a *API) OnBatchReceived(ctx Context, _operator, _from AccountID, _token_ids VecTokenId, _values VecBalance, _data VecU8) (Hash, error) {
	operatorParam := struct {
		Address AccountID
	}{
		Address: _operator,
	}

	fromParam := struct {
		Address AccountID
	}{
		Address: _from,
	}

	tokenIdsParam := struct {
		Id VecTokenId
	}{
		Id: _token_ids,
	}

	valuesParam := struct {
		Value VecBalance
	}{
		Value: _values,
	}

	dataParam := struct {
		Data VecU8
	}{
		Data: _data,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"Erc1155TokenReceiver", "on_batch_received"},
		operatorParam, fromParam, tokenIdsParam, valuesParam, dataParam,
	)
}
