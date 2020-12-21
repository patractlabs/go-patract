package test

import (
	"flag"
	"testing"
	"time"

	"github.com/patractlabs/go-patract/test/canvas"
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

// ByCanvasEnv test with canvas env
func ByCanvasEnv(t *testing.T, c func(log.Logger, Env)) {
	if !flag.Parsed() {
		flag.Parse()
	}

	argList := flag.Args()
	for _, arg := range argList {
		if arg == "extern" {
			ByExternCanvasEnv(t, c)
			return
		}
	}

	logger := log.NewLogger()
	env := canvas.NewCanvasEnv(logger)
	defer func() {
		env.Stop()
		env.Wait()
	}()

	time.Sleep(waitTimesForChainStarted) // wait chain boot
	c(logger, env)
}

// ByExternCanvasEnv test with canvas env other
func ByExternCanvasEnv(t *testing.T, c func(log.Logger, Env)) {
	c(log.NewLogger(), &envExtern{url: "ws://localhost:9944"})
}
