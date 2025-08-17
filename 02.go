package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	Parallelize(os.Args[1:], func(s string) {
		ShellCommand(s, func(c *exec.Cmd) {
			log.Println(
				Launch(c, func() {
					WaitSeconds(5, func() {
						c.Process.Kill()
					})
				}))
		})
	})
}

func ShellCommand(s string, f func(*exec.Cmd)) {
	c := exec.Command(os.Getenv("SHELL"), "-c", s)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	p := syscall.Getpid()
	log.Println("Process", p, "launching", c.String())
	f(c)
}

func Launch(c *exec.Cmd, f func()) (e error) {
	if e = c.Start(); e == nil {
		f()
		c.Wait()
	}
	if e == nil {
		e = fmt.Errorf("%v: OK", c.String())
	}
	return
}
