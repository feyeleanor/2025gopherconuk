package main

import (
	"os"
	"syscall"
)

func main() {
	RunProgram("./05_child", func(p *os.Process) {
		HandleSignal(syscall.SIGUSR1, func() {
			LogPidln("received SIGUSR1 from", p.Pid)
		})

		WaitSeconds(30, func() {
			p.Kill()
		})

		WaitSeconds(20, func() {
			p.Signal(syscall.SIGTERM)
		})

		EverySecond(20, func() {
			p.Signal(syscall.SIGINT)
		})
	})
}

func RunProgram(s string, f func(*os.Process)) {
	c := []string{s}
	a := os.ProcAttr{Files: Stdio()}
	if p, e := os.StartProcess(c[0], c, &a); e == nil {
		LogPidf("process %v launched", p.Pid)
		f(p)
		p.Wait()
		LogPidf("process %v finished", p.Pid)
	}
}
