package rpc

import (
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/rpc/native"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
)

var (
	NewCtx = api.NewCtx
)

type Contract struct {
	logger   log.Logger
	native   *native.ContractAPI
	metaData metadata.Data
}

func NewContractAPI(url string) (*Contract, error) {
	cli, err := api.NewClient(log.NewNopLogger(), url)
	if err != nil {
		return nil, errors.Wrap(err, "create client")
	}

	return &Contract{
		logger: log.NewNopLogger(),
		native: native.NewContractAPI(cli),
	}, nil
}

func (c *Contract) WithLogger(logger log.Logger) {
	c.logger = logger
	c.native.WithLogger(logger)
}

func (c *Contract) WithMetaData(bz []byte) error {
	metaData, err := metadata.New(bz)
	if err != nil {
		return err
	}

	c.metaData = *metaData
	return nil
}

func (c Contract) Native() *native.ContractAPI {
	return c.native
}
