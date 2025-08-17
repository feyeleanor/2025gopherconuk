package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for s := range c {
		LogPidln("received", s)
		if s == syscall.SIGTERM {
			LogPidln("exiting")
			close(c)
		}
	}
}
