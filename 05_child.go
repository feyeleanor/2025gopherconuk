package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if pp, e := os.FindProcess(syscall.Getppid()); e == nil {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		for s := range c {
			LogPidln("received", s)
			switch s {
			case syscall.SIGTERM:
				LogPidln("exiting")
				close(c)
			case syscall.SIGINT:
				pp.Signal(syscall.SIGUSR1)
			}
		}
	}

}
