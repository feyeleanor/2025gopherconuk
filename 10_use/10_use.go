package main

import (
	"log"
	"os"
	"unsafe"
)

/*
#cgo CFLAGS: -g -Wall
#include "10_semaphore.h"
*/
import "C"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("name of semaphore required")
	} else {
		log.Print(os.Args[1])
	}
	name := C.CString(os.Args[1])
	defer C.free(unsafe.Pointer(name))

	s, e := C.go_sem_open(name)
	if e != nil {
		log.Println(e)
		e = nil
	}
	if C.sem_wait(s) == 0 {
		log.Print("acquired lock")
	} else {
		log.Print("cannot acquire lock")
	}
	C.sem_post(s)
	log.Print("lock released")
	C.sem_close(s)
	C.sem_unlink(name)
}
