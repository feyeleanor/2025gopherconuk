package main

import (
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("name of semaphore required")
	}

	SemUnlink(os.Args[1])
	if s, e := SemOpen(os.Args[1]); e == nil {
		log.Print("open")
		if e := SemWait(s); e == nil {
			log.Println("locked")
			time.Sleep(10 * time.Second)
			log.Println("unlocking")
			SemPost(s)
			time.Sleep(10 * time.Second)
		} else {
			log.Print(e)
		}
		SemClose(s)
		SemUnlink(os.Args[1])
	} else {
		log.Fatal(e)
	}
}
