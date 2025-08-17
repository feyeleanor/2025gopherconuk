package main

import (
	"log"
	"os"
	"time"
	"unsafe"
)

/*
#cgo CFLAGS: -g -Wall
#include "10_semaphore.h"
*/
import "C"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no semaphore name")
	} else {
		log.Print(os.Args[1])
	}
	name := C.CString(os.Args[1])
	defer C.free(unsafe.Pointer(name))

	s := C.go_sem_open(name)
	log.Print("semaphore opening:")
	if C.sem_wait(s) == 0 {
		log.Print("locked semaphore")
		time.Sleep(10 * time.Second)
		log.Print("unlocking semaphore")
		C.sem_post(s)
		time.Sleep(10 * time.Second)
	} else {
		log.Print("can't lock semaphore")
	}
	C.sem_close(s)
	C.sem_unlink(name)
}
