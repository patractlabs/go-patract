package api

import (
	"context"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/patractlabs/go-patract/utils/log"
)

// Context for api
type Context struct {
	context.Context

	logger log.Logger
	from   signature.KeyringPair
}

// NewCtx creates a new Context
func NewCtx(ctx context.Context) Context {
	return Context{
		Context: ctx,
		logger:  log.NewLogger(),
	}
}

// WithLogger set logger for ctx
func (c Context) WithLogger(logger log.Logger) Context {
	c.logger = logger
	return c
}

// WithFrom set from keyring pair for ctx
func (c Context) WithFrom(from signature.KeyringPair) Context {
	c.from = from
	return c
}

// From returns keyring pair for ctx
func (c Context) From() signature.KeyringPair {
	return c.from
}
