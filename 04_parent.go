package main

import (
	"os"
	"syscall"
)

func main() {
	RunProgram("./04_child", func(p *os.Process) {
		WaitSeconds(30, func() {
			p.Kill()
		})

		EverySecond(20, func() {
			p.Signal(syscall.SIGINT)
		})

		WaitSeconds(20, func() {
			p.Signal(syscall.SIGTERM)
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
