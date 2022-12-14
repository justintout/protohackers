package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/justintout/protohackers"
	"github.com/justintout/protohackers/speed"
)

const (
	addr6 = "0.0.0.0:10006"
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

	fmt.Printf("starting problem 4 server: %s\n", protohackers.Addr4)
	s4 := protohackers.DatabaseServer{}
	go s4.ListenAndServe(protohackers.Addr4)

	fmt.Printf("starting problem 5 server: %s\n", protohackers.Addr5)
	s5 := protohackers.ChatProxyServer{}
	go s5.ListenAndServe(protohackers.Addr5, protohackers.UpstreamChat)

	fmt.Printf("starting problem 6 server: %s\n", addr6)
	s6 := speed.Server{}
	go s6.ListenAndServe(addr6)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	s0.Close()
	s1.Close()
	s2.Close()
	s3.Close()
	s4.Close()
	s5.Close()
}
