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

	fmt.Printf("starting problem 2 server: %s\n", protohackers.Addr2)
	s2 := protohackers.MeanServer{}
	go s2.ListenAndServe(protohackers.Addr2)

	fmt.Printf("starting problem 3 server: %s\n", protohackers.Addr3)
	s3 := protohackers.ChatServer{}
	go s3.ListenAndServe(protohackers.Addr3)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	s0.Close()
	s1.Close()
	s2.Close()
	s3.Close()
}
