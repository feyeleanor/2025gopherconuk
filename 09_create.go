package main

import (
	"log"
	"os"
	"time"
	"unsafe"
)

/*
#include <sys/errno.h>
#include <stdlib.h>
#include <semaphore.h>

sem_t *go_sem_open(const char *name) {
	return sem_open(name, O_CREAT, 0644, 1);
}
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
	log.Print("open")
	if C.sem_wait(s) == 0 {
		log.Print("locked")
		time.Sleep(10 * time.Second)
		log.Print("unlocking")
		C.sem_post(s)
		time.Sleep(10 * time.Second)
	} else {
		log.Print("can't lock")
	}
	C.sem_close(s)
	C.sem_unlink(name)
}
