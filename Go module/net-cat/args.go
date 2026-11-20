package main

import (
	"fmt"
	"os"
)

type opts struct {
	addr string
}

func parseArgs() (opts, error) {
	if len(os.Args) > 2 {
		return opts{}, fmt.Errorf("[USAGE]: ./TCPChat $port")
	}

	op := opts{addr: ":8989"}
	if len(os.Args) == 2 {
		op.addr = ":" + os.Args[1]
	}

	return op, nil
}
