package rpc

import (
	"github.com/patractlabs/go-patract/api"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/rpc/native"
	"github.com/patractlabs/go-patract/types"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
)

var (
	NewCtx = api.NewCtx
)

type Contract struct {
	logger    log.Logger
	native    *native.ContractAPI
	metaData  metadata.Data
	ss58Codec *types.SS58Codec
}

func NewContractAPI(url string) (*Contract, error) {
	cli, err := api.NewClient(log.NewNopLogger(), url)
	if err != nil {
		return nil, errors.Wrap(err, "create client")
	}

	return &Contract{
		logger:    log.NewNopLogger(),
		native:    native.NewContractAPI(cli),
		ss58Codec: types.GetDefaultSS58Codec(),
	}, nil
}

func (c *Contract) WithLogger(logger log.Logger) {
	c.logger = logger
	c.native.WithLogger(logger)
}

func (c *Contract) WithSS58Codec(ss58Codec *types.SS58Codec) {
	c.ss58Codec = ss58Codec
}

func (c *Contract) WithMetaData(bz []byte) error {
	metaData, err := metadata.New(bz)
	if err != nil {
		return err
	}

	c.metaData = *metaData
	return nil
}

// Native get native api for contract
func (c Contract) Native() *native.ContractAPI {
	return c.native
}

// Logger get logger for contract api
func (c Contract) Logger() log.Logger {
	return c.logger
}
