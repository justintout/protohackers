package main

import (
	"fmt"

	"github.com/justintout/protohackers"
)

func main() {
	fmt.Printf("starting problem 0 server: %s\n", protohackers.Addr0)
	s0 := protohackers.EchoServer{}
	s0.ListenAndServe(protohackers.Addr0)
}
