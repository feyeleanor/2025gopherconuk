package main

import (
	"io"
	"os"
	"syscall"
)

func main() {
	defer func() {
		LogPidln("exiting")
	}()

	ForParent(func(p *os.Process) {
		i := 1
		HandleSignal(syscall.SIGUSR2, func() {
			LogPidf("received SIGUSR2 x %v", i)
			i++
		})

		EverySecond(5, func() {
			SendMessage(os.Stdout, "HELLO")
			ReceivePids(func(i int) {
				ForProcess(i, func(p *os.Process) {
					if i != syscall.Getpid() {
						LogPidln("sending signal to", p.Pid)
						p.Signal(syscall.SIGUSR2)
					}
				})
			})
		}).Wait()
	})
}

func ReceivePids(f func(int)) {
	switch b, e := ReceiveMessage(os.Stdin); e {
	case nil, io.EOF:
		ForEachInt(b, func(i int) {
			f(i)
		})
	default:
		LogPidln(e)
	}
}
