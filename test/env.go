package test

import (
	"flag"
	"testing"
	"time"

	"github.com/patractlabs/go-patract/test/chainnode"
	"github.com/patractlabs/go-patract/utils/log"
)

const (
	waitTimesForChainStarted = 300 * time.Millisecond
)

// Env env interface for test
type Env interface {
	URL() string
	IsUseExtToTest() bool
}

type envExtern struct {
	url string
}

func (e envExtern) URL() string {
	return e.url
}

func (e envExtern) IsUseExtToTest() bool {
	return true
}

// ByNodeEnv test with europa env
func ByNodeEnv(t *testing.T, c func(log.Logger, Env)) {
	if !flag.Parsed() {
		flag.Parse()
	}

	argList := flag.Args()
	for _, arg := range argList {
		if arg == "extern" {
			ByExternEnv(t, c)
			return
		}
	}

	logger := log.NewLogger()
	env := chainnode.NewCanvasEnv(logger)
	defer func() {
		env.Stop()
		env.Wait()
	}()

	time.Sleep(waitTimesForChainStarted) // wait chain boot
	c(logger, env)
}

// ByExternEnv test with ext env other
func ByExternEnv(t *testing.T, c func(log.Logger, Env)) {
	c(log.NewLogger(), &envExtern{url: "ws://localhost:9944"})
}
