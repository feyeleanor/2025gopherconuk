package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

func SyscallError(r uintptr, e syscall.Errno, s string) error {
	if r != 0 {
		return ErrorCode(s, e)
	}
	return nil
}

func CString(s string) unsafe.Pointer {
	p, e := syscall.BytePtrFromString(s)
	if e != nil {
		log.Fatalf("unable to convert %v to a C-style string", s)
	}
	return unsafe.Pointer(p)
}

func ErrorCode(s string, e syscall.Errno) error {
	return errors.New(s + ": " + strconv.FormatUint(uint64(e), 10))
}

func ForFifo(s string, flag int, f func(*os.File)) {
	if p, e := os.OpenFile(s, flag, os.ModeNamedPipe); e == nil {
		defer p.Close()
		f(p)
	} else {
		log.Fatal(e)
	}
}

func WaitSeconds(n time.Duration, f func()) {
	time.AfterFunc(n*time.Second, f)
}

func LogPidFatal(v ...any) {
	log.Fatal(v...)
}

func LogPidln(v ...any) {
	log.Println(
		append(
			[]any{fmt.Sprintf("%v:", syscall.Getpid())},
			v...)...)
}

func LogPidf(s string, v ...any) {
	log.Printf("%v: %v", syscall.Getpid(), fmt.Sprintf(s, v...))
}

func ForEachInt(s []byte, f func(int)) {
	for _, m := range Tokens(s) {
		if i, e := strconv.Atoi(m); e == nil {
			f(i)
		} else {
			LogPidln(e, m)
		}
	}
}

func ForEachEntry(m map[int]bool, f func(p *os.Process)) {
	for p, _ := range m {
		ForProcess(p, f)
	}
}

func ToggleEntry(c map[int]bool, s string) {
	if i, e := strconv.Atoi(s); e == nil {
		ForProcess(i, func(p *os.Process) {
			if v, _ := c[p.Pid]; !v {
				log.Println(p.Pid, "toggled on")
				c[p.Pid] = true
			} else {
				log.Println(p.Pid, "toggled off")
				delete(c, p.Pid)
			}
		})
	} else {
		log.Println(e)
	}
}

func ForParent(f func(*os.Process)) {
	ForProcess(syscall.Getppid(), f)
}

func ForProcess(p int, f func(*os.Process)) {
	LogPidln("ForProcess", p)
	if p, e := os.FindProcess(p); e == nil {
		if e := p.Signal(syscall.Signal(0)); e == nil {
			f(p)
		} else {
			LogPidf("no such process as %v", p.Pid)
		}
	} else {
		LogPidFatal(e)
	}
}

func TaskList(s string, i int) (r []string) {
	r = make([]string, 0, i)
	for range i {
		r = append(r, s)
	}
	return
}

func TryWait(p *os.Process) {
	if s, e := p.Wait(); e != nil {
		LogPidln("unable to Wait on process", p.Pid)
	} else {
		LogPidln("process finished", s.Pid)
		LogPidln(s.Pid, "exit code", s.ExitCode())
	}
}

func Peers(p *os.Process, c ...int) (r []string) {
	for _, v := range c {
		if v != p.Pid {
			r = append(r, strconv.Itoa(v))
		}
	}
	return
}

func ChildStdio(f func(...*os.File)) {
	i := Pipeio()
	o := Pipeio()
	defer i.w.Close()
	defer i.r.Close()
	defer o.w.Close()
	defer o.r.Close()
	f(i.r, i.w, o.r, o.w)
}

type Pipe struct {
	r, w *os.File
}

func Pipeio() Pipe {
	in, out, e := os.Pipe()
	if e != nil {
		LogPidFatal(e)
	}
	return Pipe{in, out}
}

func MessageLoop(r io.Reader, f func(string)) {
	for {
		switch s, e := ReceiveMessage(r); e {
		case nil:
			for _, m := range Tokens(s) {
				f(m)
			}
		case io.EOF:
			for _, m := range Tokens(s) {
				f(m)
			}
			return
		default:
			LogPidln(e)
			return
		}
	}
}

func Tokens(b []byte) []string {
	return strings.Split(string(b), string(' '))
}

func SendMessage(w io.Writer, s ...any) {
	for _, v := range s {
		switch v := v.(type) {
		case []byte:
			w.Write(v)
			w.Write([]byte(string('\n')))
		case string:
			SendMessage(w, []byte(v))
		case rune:
			SendMessage(w, string(v))
		case fmt.Stringer:
			SendMessage(w, v.String())
		default:
			LogPidf("unable to send message [%T] %v", v, v)
		}
	}
}

func ReceiveMessage(r io.Reader) (b []byte, e error) {
	switch b, e = bufio.NewReader(r).ReadBytes('\n'); {
	case e == nil, e == io.EOF:
		b = DropTail(b, 1)
	default:
		LogPidln(e)
	}
	return
}

func DropTail(b []byte, n int) []byte {
	l := len(b) - n
	if l < 0 {
		l = 0
	}
	return b[:l]
}

func HandleSignal(s os.Signal, f func()) {
	c := make(chan os.Signal)
	signal.Notify(c, s)

	go func(c chan os.Signal) {
		for _ = range c {
			f()
		}
	}(c)
}

func Stdio() []*os.File {
	return []*os.File{os.Stdin, os.Stdout, os.Stderr}
}

func Parallelize[T any](s []T, f func(T)) {
	var w sync.WaitGroup

	for _, n := range s {
		w.Add(1)
		go func(n T) {
			defer w.Done()

			f(n)
		}(n)
	}
	w.Wait()
}

func EverySecond(n int, f func()) (w *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Second)

	w = new(sync.WaitGroup)
	w.Add(1)
	go func() {
		defer w.Done()

		for {
			select {
			case <-ticker.C:
				if n--; n > 0 {
					f()
				} else {
					return
				}
			}
		}
	}()
	return
}
