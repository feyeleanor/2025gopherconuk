package main

import (
	"log"
	"os"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("parameter required: name of fifo to open")
	}

	i := 1
	HandleSignal(syscall.SIGUSR2, func() {
		log.Println("received SIGUSR2 x", i)
		i++
	})

	FifoDial(os.Args[1], func(f *os.File) {
		LogPidln("connected to Fifi", f.Fd())
		EverySecond(10, func() {
			SendPid(f)
		}).Wait()
	})
}

func FifoDial(s string, f func(*os.File)) {
	LogPidln("connecting to Fifo", s)
	if _, e := os.Stat(s); e != nil {
		log.Fatal("cannot access fifo")
	}
	ForFifo(s, os.O_WRONLY, f)
}

func SendPid(f *os.File) {
	SendMessage(f, strconv.Itoa(syscall.Getpid()))
}
