package protohackers

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

const Addr4 = "0.0.0.0:10004"

type DatabaseServer struct {
	listener *net.UDPConn

	mu      *sync.Mutex
	storage map[string]string
}

func (s *DatabaseServer) ListenAndServe(addr string) {
	s.mu = new(sync.Mutex)
	s.storage = make(map[string]string)
	ua, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", ua)
	if err != nil {
		panic(err)
	}
	s.listener = conn
	defer conn.Close()
	for {
		b := make([]byte, 1000)
		n, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			panic(err)
		}
		p := string(b[:n])
		fmt.Println("-> " + p)
		go handle4(s, addr, p)
	}
}

func (s *DatabaseServer) Close() {
	s.listener.Close()
}

func (s *DatabaseServer) set(key, value string) {
	s.mu.Lock()
	s.storage[key] = value
	s.mu.Unlock()
}

func (s *DatabaseServer) get(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.storage[key]
}

const version = "protohackers problem 4 1.0"

func handle4(svr *DatabaseServer, addr *net.UDPAddr, p string) {
	if strings.Contains(p, "=") {
		handleDBInsert(svr, p)
		return
	}
	handleDBRetrieve(svr, addr, p)
}

func handleDBInsert(svr *DatabaseServer, req string) {
	s := strings.SplitN(req, "=", 2)
	if len(s) == 1 {
		s = append(s, "")
	}
	svr.set(s[0], s[1])
}

func handleDBRetrieve(svr *DatabaseServer, addr *net.UDPAddr, req string) {
	val := svr.get(req)
	if req == "version" {
		val = version
	}
	svr.listener.WriteToUDP([]byte(req+"="+val), addr)
}
