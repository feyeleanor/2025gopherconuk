package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	c := exec.Command("echo", "Hello, Golang!")
	if o, e := c.Output(); e == nil {
		fmt.Println("Output:", string(o))
	} else {
		log.Println("Error:", e)
	}
}