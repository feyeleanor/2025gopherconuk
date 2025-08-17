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

	if s, e := SemOpen(os.Args[1]); e == nil {
		log.Print("open")
		if e := SemWait(s); e == nil {
			log.Print("locked")
			time.Sleep(10 * time.Second)
			log.Println("unlocking")
			SemPost(s)
		} else {
			log.Println(e)
		}
		SemClose(s)
		SemUnlink(os.Args[1])
	} else {
		log.Fatal(e)
	}
}
