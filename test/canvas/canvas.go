package canvas

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/jesselucas/executil"
	"github.com/patractlabs/go-patract/utils/log"
	"github.com/pkg/errors"
)

const (
	portWs         = "39944"
	portPrometheus = "39615"
)

// Env a canvas environment for testing
type Env struct {
	wg sync.WaitGroup

	log   log.Logger
	mutex sync.RWMutex
	pID   int
}

// NewCanvasEnv create a canvas chain to test
func NewCanvasEnv(log log.Logger) *Env {
	res := &Env{
		log: log,
	}
	if err := res.Start(); err != nil {
		panic(err)
	}

	return res
}

// URL get the url to the canvas
func (c *Env) URL() string {
	return fmt.Sprintf("ws://localhost:%s", portWs)
}

// PID get canvas process id
func (c *Env) PID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.pID
}

func (c *Env) setPID(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pID = id
}

// Start start the canvas process
func (c *Env) Start() error {
	c.log.Debug("start canvas env")

	outputChan := make(chan string)

	cmd := executil.Command("canvas", "--tmp", "--dev",
		"--ws-port", portWs,
		"--prometheus-port", portPrometheus,
	)
	cmd.OutputChan = outputChan

	c.wg.Add(1)
	go func() {
		defer func() {
			c.wg.Done()
			c.log.Debug("stop canvas cmd goroutine", "PID", c.PID())
			close(outputChan)
		}()

		err := cmd.Start()

		if err != nil {
			c.log.Error("start canvas cmd error", "err", err)
			panic(err)
		}

		c.setPID(cmd.Cmd.Process.Pid)

		// wait stop
		if err := cmd.Wait(); err != nil {
			c.log.Error("wait cmd error", "err", err)
		}
	}()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			out, stop := <-outputChan
			if !stop {
				c.log.Debug("stop output goroutine", "PID", c.PID())
				return
			}
			c.processOutput(out)
		}
	}()

	return nil
}

func (c *Env) processOutput(str string) {
	c.log.Debug("canvas log", "str", str)
}

// Stop stop the environment
func (c *Env) Stop() {
	pID := c.PID()
	if pID == 0 {
		return
	}

	c.log.Debug("kill", "PID", c.PID())

	var (
		runCmd = "kill"
		pIDStr = fmt.Sprintf("%d", pID)
	)

	cmd := exec.Command(runCmd, pIDStr)

	if err := cmd.Start(); err != nil {
		panic(errors.Wrap(err, "kill error"))
	}

	if err := cmd.Wait(); err != nil {
		panic(errors.Wrap(err, "kill wait error"))
	}
}

// Wait wait for env stopped
func (c *Env) Wait() {
	c.wg.Wait()
}
