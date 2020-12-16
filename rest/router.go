package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/patractlabs/go-patract/metadata"
	"github.com/patractlabs/go-patract/rpc"
	"github.com/patractlabs/go-patract/types"
	"github.com/pkg/errors"
)

type Router struct {
	*gin.Engine

	cli             *rpc.Contract
	runtimeMetadata *types.Metadata
	metadatas       map[string]*metadata.Data
}

func NewRouter(router *gin.Engine) Router {
	return Router{
		Engine:    router,
		metadatas: make(map[string]*metadata.Data, 16),
	}
}

func (r *Router) WithMetaData(data *metadata.Data) {
	if _, ok := r.metadatas[data.Contract.Name]; ok {
		panic(errors.Errorf("has with metadata by %s", data.Contract.Name))
	}
	r.metadatas[data.Contract.Name] = data
}

func (r *Router) WithCli(cli *rpc.Contract) {
	r.cli = cli
}

func (r *Router) WithRuntimeMetadata(m *types.Metadata) {
	r.runtimeMetadata = m
}

func (r *Router) Init() {
	r.messages("/")
}
