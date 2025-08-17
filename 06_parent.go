package main

import (
	"os"
	"strings"
)

func main() {
	t := TaskList("./06_child", 3)
	c := []int{}
	Parallelize(t, func(s string) {
		RunProgram(s, func(p *os.Process, i, o *os.File) {
			c = append(c, p.Pid)
			WaitSeconds(10, func() {
				p.Kill()
				i.Close()
			})

			MessageLoop(i, func(s string) {
				LogPidf("message [%v] from %v", s, p.Pid)
				LogPidln("children:", c)
				x := Peers(p, c...)
				LogPidln("peers:", x)
				SendMessage(o, strings.Join(x, " "))
			})
		})
	})
}

func RunProgram(s string, f func(*os.Process, *os.File, *os.File)) {
	c := []string{s}
	ChildStdio(func(io ...*os.File) {
		a := os.ProcAttr{
			Files: []*os.File{io[2], io[1], os.Stderr}}

		if p, e := os.StartProcess(c[0], c, &a); e == nil {
			LogPidln("process launched", p.Pid)
			f(p, io[0], io[3])
			TryWait(p)
		}
	})
}
