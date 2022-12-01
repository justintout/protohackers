package protohackers

import (
	"io"
	"net"
)

const Addr0 = "0.0.0.0:10000"

type EchoServer struct {
}

func (s *EchoServer) ListenAndServe(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	if _, err := io.Copy(conn, conn); err != nil {
		panic(err)
	}
	conn.Close()
}
