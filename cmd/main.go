package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/justintout/protohackers"
)

// TODO: shutdown is sketchy and leaks. needs channels.
func main() {
	fmt.Printf("starting problem 0 server: %s\n", protohackers.Addr0)
	s0 := protohackers.EchoServer{}
	go s0.ListenAndServe(protohackers.Addr0)

	fmt.Printf("starting problem 1 server: %s\n", protohackers.Addr1)
	s1 := protohackers.PrimeServer{}
	go s1.ListenAndServe(protohackers.Addr1)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	s0.Close()
	s1.Close()
}
