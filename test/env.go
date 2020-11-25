package test

import (
	"testing"
	"time"

	"github.com/patractlabs/go-patract/test/canvas"
	"github.com/patractlabs/go-patract/utils/log"
)

// ByCanvasEnv test with canvas env
func ByCanvasEnv(t *testing.T, c func(log.Logger, *canvas.Env)) {
	logger := log.NewLogger()
	env := canvas.NewCanvasEnv(logger)
	defer func() {
		env.Stop()
		env.Wait()
	}()

	time.Sleep(1 * time.Second) // wait chain boot
	c(logger, env)
}
