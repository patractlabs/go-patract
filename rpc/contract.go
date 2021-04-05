package rpc

import (
	"bytes"

	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils"
	"github.com/patractlabs/go-patract/utils/log"
)

// Instantiate is a new contract from the existent `codeHash`
// optionally transferring some balance.
func (c *Contract) Instantiate(
	ctx api.Context,
	endowment types.CompactBalance,
	gasLimit types.CompactGas,
	codeHash types.CodeHash,
	salt []byte,
	args ...interface{}) (types.Hash, types.AccountID, error) {
	data, err := c.getConstructorsData([]string{"new"}, args...)
	if err != nil {
		return types.Hash{}, types.AccountID{}, err
	}

	if salt == nil {
		salt = []byte(utils.GenerateSalt())
	}

	hash, err := c.native.Instantiate(ctx, endowment, gasLimit, codeHash, data, salt)
	if err != nil {
		return types.Hash{}, types.AccountID{}, err
	}

	contractAccount := GetContractAccountID(types.NewAccountID(ctx.From().PublicKey), codeHash, salt)

	return hash, contractAccount, nil
}

// InstantiateWithCode is a new contract from the supplied `code`
// optionally transferring some balance.
func (c *Contract) InstantiateWithCode(
	ctx api.Context,
	logger log.Logger,
	endowment types.CompactBalance,
	gasLimit types.CompactGas,
	codeHash types.CodeHash,
	code []byte,
	salt []byte,
	args ...interface{}) (types.Hash, types.AccountID, error) {
	data, err := c.getConstructorsData([]string{"new"}, args...)
	if err != nil {
		return types.Hash{}, types.AccountID{}, err
	}

	if salt == nil {
		salt = []byte(utils.GenerateSalt())
	}

	hash, err := c.native.InstantiateWithCode(
		ctx, endowment, gasLimit, code, data, salt)

	if err != nil {
		return types.Hash{}, types.AccountID{}, err
	}

	contractAccount := GetContractAccountID(types.NewAccountID(ctx.From().PublicKey), codeHash, salt)
	contractAddress, _ := c.ss58Codec.EncodeAccountID(contractAccount)

	logger.Info("InstantiateWithCode", "deployer", ctx.From().Address, "salt", string(salt),
		"hash", types.HexEncodeToString(hash[:]), "contractAccount", contractAddress)
	return hash, contractAccount, nil
}

// GetContractAccountID get constract account ID
func GetContractAccountID(deployer types.AccountID, codeHash types.CodeHash, salt []byte) types.AccountID {
	contractHashs := bytes.NewBuffer(make([]byte, 0, 256))

	contractHashs.Write(deployer[:])
	contractHashs.Write(codeHash[:])
	contractHashs.Write(salt)

	return types.NewAccountID(utils.Hash256(contractHashs.Bytes()))
}

// UpdateSchedule update schedule for contract
func (c *Contract) UpdateSchedule(ctx api.Context, schedule types.Schedule) (types.Hash, error) {
	return c.native.Cli.SubmitAndWaitExtrinsic(ctx, "Contracts.update_schedule", schedule)
}

// ClaimSurcharge claim surcharge
func (c *Contract) ClaimSurcharge(
	ctx api.Context,
	dest types.AccountID,
	auxSender types.OptionAccountID) (types.Hash, error) {
	return c.native.Cli.SubmitAndWaitExtrinsic(ctx, "Contracts.claim_surcharge", dest, auxSender)
}
