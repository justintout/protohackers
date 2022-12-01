package protohackers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
)

const Addr1 = "0.0.0.0:10001"

type PrimeServer struct {
	listener net.Listener
}

func (s *PrimeServer) ListenAndServe(addr string) {
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
		go handle1(conn)
	}
}

func (s *PrimeServer) Close() {
	s.listener.Close()
}

const reqMethod = "isPrime"

type request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func handle1(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("err: %v\n", err)
			continue
		}
		var req request
		if err := json.Unmarshal(b, &req); err != nil {
			writeMalformed(conn)
			continue
		}
		if req.Method != reqMethod || req.Number == nil {
			writeMalformed(conn)
			continue
		}
		x := big.NewFloat(*req.Number)
		y, _ := x.Int(nil)
		res := response{
			Method: reqMethod,
			Prime:  x.IsInt() && y.ProbablyPrime(20),
		}
		if err := json.NewEncoder(conn).Encode(res); err != nil {
			panic(err)
		}
	}
	conn.Close()
}

func writeMalformed(conn net.Conn) {
	conn.Write([]byte("\n"))
}
