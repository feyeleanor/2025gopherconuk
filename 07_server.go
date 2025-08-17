package main

import (
	"log"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("needs name of fifo to create")
	}

	m := map[int]bool{}
	EverySecond(20, func() {
		ForEachEntry(m, func(p *os.Process) {
			log.Println("sending SIGUSR2 to", p.Pid)
			p.Signal(syscall.SIGUSR2)
		})
	})

	for {
		FifoListen(os.Args[1], func(f *os.File) {
			MessageLoop(f, func(s string) {
				log.Printf("received pid [%v]", s)
				ToggleEntry(m, s)
			})
		})
	}
}

func FifoListen(s string, f func(*os.File)) {
	if _, e := os.Stat(s); e == nil {
		if e = os.Remove(s); e != nil {
			log.Fatal("cannot delete existing fifo")
		}
	}

	if e := unix.Mkfifo(s, 0666); e != nil {
		log.Fatal(e)
	}

	defer func() {
		os.Remove(s)
	}()

	ForFifo(s, os.O_RDONLY, f)
}
