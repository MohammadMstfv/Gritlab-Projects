package main

import (
	"fmt"
	"log"
)

func main() {
	op, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := runServer(op); err != nil {
		log.Fatal(err)
	}
}
