package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func HoldToClose(waitFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	waitFunc()
}
