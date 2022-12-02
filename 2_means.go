package protohackers

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const Addr2 = "0.0.0.0:10002"

type MeanServer struct {
	listener net.Listener
}

func (s *MeanServer) ListenAndServe(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s.listener = l
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handle2(conn)
	}
}

func (s *MeanServer) Close() {
	s.listener.Close()
}

type meanMessage struct {
	Typ [1]byte
	F   [4]byte
	S   [4]byte
}

func (m meanMessage) IsInsert() bool {
	return m.Typ[0] == 'I'
}

func (m meanMessage) IsQuery() bool {
	return m.Typ[0] == 'Q'
}

func (m meanMessage) First() int32 {
	return int32(binary.BigEndian.Uint32(m.F[:]))
}

func (m meanMessage) Second() int32 {
	return int32(binary.BigEndian.Uint32(m.S[:]))
}

type storage map[int32]int32

func handle2(conn net.Conn) {
	defer conn.Close()
	s := make(storage)

	for {
		var m meanMessage
		if err := binary.Read(conn, binary.BigEndian, &m); err != nil {
			if err == io.EOF {
				conn.Close()
				return
			}
			fmt.Printf("2: read err: %v\n", err)
			return
		}
		// fmt.Printf("Decoded:\t%c\t%d\t%d\n", m.Typ, m.First(), m.Second())
		switch {
		case m.IsInsert():
			handleInsert(conn, s, m)
		case m.IsQuery():
			handleQuery(conn, s, m)
		default:
			conn.Close()
		}
	}
}

func handleInsert(conn net.Conn, s storage, m meanMessage) {
	s[m.First()] = m.Second()
}

func handleQuery(conn net.Conn, store storage, m meanMessage) {
	var (
		s int
		n int
	)
	for ts, p := range store {
		if m.First() <= ts && ts <= m.Second() {
			s += int(p)
			n++
		}
	}
	if m.First() > m.Second() || n == 0 {
		binary.Write(conn, binary.BigEndian, int32(0))
		return
	}
	// fmt.Printf("Mean: %d / %d = %d\n", s, n, s/n)
	binary.Write(conn, binary.BigEndian, int32(s/n))
}
