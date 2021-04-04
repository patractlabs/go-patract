package chainnode_test

import (
	"testing"
	"time"

	"github.com/patractlabs/go-patract/test/chainnode"
	"github.com/patractlabs/go-patract/utils/log"
)

func TestCanvas(t *testing.T) {
	env := chainnode.NewCanvasEnv(log.NewLogger())

	time.Sleep(1 * time.Second)

	env.Stop()
	env.Wait()
}
