package speed

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
	messages chan message
	log      []message
}

func (s *Server) ListenAndServe(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s.listener = l
	defer l.Close()
	messages := make(chan message)
	go s.listen(messages)
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		_, m := newParser(conn)
		go func() {
			for message := range m {
				messages <- message
			}
		}()
	}
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) listen(messages chan message) {
	for message := range messages {
		s.log = append(s.log, message)
		switch message.typ() {
		}
		s.handleMessage(message)
	}
}

func (s *Server) handleMessage(message message) {
	fmt.Printf("%s\n", message)
}

const eof = 0x00

// client types
const (
	clientTypCamera     uint8 = 0x80
	clientTypDispatcher uint8 = 0x81
)

type message interface {
	typ() uint8
	String() string
}
