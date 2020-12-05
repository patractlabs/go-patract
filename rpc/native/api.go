package native

import (
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
)

// ContractAPI the api for contract
type ContractAPI struct {
	Cli *api.Client
}

// NewContractAPI creates a Contract api instance from cli
func NewContractAPI(cli *api.Client) *ContractAPI {
	return &ContractAPI{
		Cli: cli,
	}
}

// WithLogger set logger
func (c *ContractAPI) WithLogger(logger log.Logger) {
	c.Cli.WithLogger(logger)
}

// UpdateSchedule Updates the schedule for metering contracts.
func (c *ContractAPI) UpdateSchedule(ctx api.Context, schedule types.Schedule) (string, error) {
	return c.Cli.SubmitAndWaitExtrinsic(ctx, "Contracts.update_schedule", schedule)
}

// PutCode submit put_code to chain
func (c *ContractAPI) PutCode(ctx api.Context, code []byte) (string, error) {
	return c.Cli.SubmitAndWaitExtrinsic(ctx, "Contracts.put_code", code)
}

// Instantiate is a new contract from the `codehash` generated by `put_code`,
// optionally transferring some balance.
func (c *ContractAPI) Instantiate(
	ctx api.Context,
	endowment types.CompactBalance,
	gasLimit types.CompactGas,
	codeHash types.CodeHash,
	data []byte) (string, error) {
	return c.Cli.SubmitAndWaitExtrinsic(
		ctx, "Contracts.instantiate", endowment, gasLimit, codeHash, data)
}

// Call Makes a call to an account, optionally transferring some balance.
func (c *ContractAPI) Call(
	ctx api.Context,
	dest types.AccountID,
	value types.Balance,
	gasLimit types.Gas,
	data []byte,
) (string, error) {
	return c.Cli.SubmitAndWaitExtrinsic(
		ctx, "Contracts.call", dest, value, gasLimit, data)
}

// ClaimSurcharge Allows block producers to claim a small reward for evicting a contract.
// If a block producer fails to do so, a regular users will be allowed to claim the reward.
func (c *ContractAPI) ClaimSurcharge(
	ctx api.Context,
	dest types.AccountID,
	auxSender types.AccountID,
) (string, error) {
	return c.Cli.SubmitAndWaitExtrinsic(
		ctx, "Contracts.claim_surcharge", dest, auxSender)
}
