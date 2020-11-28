package native

import "github.com/patractlabs/go-patract/api"

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
