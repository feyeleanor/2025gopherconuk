package main

import (
	"log"
	"os"
	"strings"
	"syscall"
)

func main() {
	Parallelize(os.Args[1:], func(s string) {
		ShellCommand(s, func(p *os.Process) {
			WaitSeconds(10, func() {
				p.Kill()
			})
		})
	})
}

func ShellCommand(s string, f func(*os.Process)) {
	c := []string{os.Getenv("SHELL"), "-c", s}
	log.Println("Process", syscall.Getpid(), "preparing", strings.Join(c, " "))
	if p, e := os.StartProcess(c[0], c, &os.ProcAttr{Files: Stdio()}); e == nil {
		f(p)
		p.Wait()
	}
}
