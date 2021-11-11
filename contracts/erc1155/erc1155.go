package erc1155

import (
	"github.com/patractlabs/go-patract/test"
	"github.com/patractlabs/go-patract/types"
)

func (a *API) Create(ctx Context, value Balance) (Hash, error) {
	valueParam := struct {
		Value Balance
	}{
		Value: value,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"create"},
		valueParam,
	)
}

func (a *API) Mint(ctx Context, tokenId TokenId, value Balance) (Hash, error) {
	tokenIdParam := struct {
		Id TokenId
	}{
		Id: tokenId,
	}

	valueParam := struct {
		value Balance
	}{
		value: value,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"mint"},
		tokenIdParam, valueParam,
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

	tokenIdParam := struct {
		Id TokenId
	}{
		Id: tokenId,
	}

	valueParam := struct {
		value Balance
	}{
		value: value,
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
		[]string{"safe_transfer_from"},
		fromParam, toParam, tokenIdParam, valueParam, dataParam,
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
		[]string{"safe_batch_transfer_from"},
		fromParam, toParam, tokenIdsParam, valuesParam, dataParam,
	)
}

func (a *API) BalanceOf(ctx Context, owner AccountID, tokenId TokenId) (Balance, error) {
	ownerParam := struct {
		Address AccountID
	}{
		Address: owner,
	}

	tokenIdParam := struct {
		Id TokenId
	}{
		Id: tokenId,
	}

	var res Balance

	err := a.CallToRead(ctx,
		&res,
		a.ContractAccountID,
		[]string{"balance_of"},
		ownerParam, tokenIdParam,
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
		[]string{"balance_of_batch"},
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

	approvedParam := struct {
		Address Bool
	}{
		Address: approved,
	}

	return a.CallToExec(ctx,
		a.ContractAccountID,
		types.NewCompactBalance(0),
		types.NewCompactGas(test.DefaultGas),
		[]string{"set_approval_for_all"},
		operatorParam, approvedParam,
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
		[]string{"is_approved_for_all"},
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

	tokenIdParam := struct {
		Id TokenId
	}{
		Id: _token_id,
	}

	valueParam := struct {
		Value Balance
	}{
		Value: _value,
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
		[]string{"on_received"},
		operatorParam, fromParam, tokenIdParam, valueParam, dataParam,
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
		[]string{"on_batch_received"},
		operatorParam, fromParam, tokenIdsParam, valuesParam, dataParam,
	)
}
