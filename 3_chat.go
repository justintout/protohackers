package protohackers

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
)

const Addr3 = "0.0.0.0:10003"

type ChatServer struct {
	listener net.Listener

	mu      *sync.Mutex
	clients []*client
}

func (s *ChatServer) ListenAndServe(addr string) {
	s.mu = new(sync.Mutex)
	s.clients = make([]*client, 0, 10)
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
		go handle3(s, conn)
	}
}

func (s *ChatServer) Close() {
	s.listener.Close()
}

func (s *ChatServer) list() []string {
	l := make([]string, 0, len(s.clients))
	for _, c := range s.clients {
		l = append(l, c.Name)
	}
	return l
}

func (s *ChatServer) join(c *client) {
	s.mu.Lock()
	c.receive(svrLead + strings.Join(s.list(), ", "))
	s.relay("", c.Name+" "+joined)
	s.clients = append(s.clients, c)
	s.mu.Unlock()
}

func (s *ChatServer) leave(c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var (
		idx   int
		found bool
	)
	for i, cl := range s.clients {
		if c == cl {
			found = true
			idx = i
			break
		}
	}
	if !found {
		return
	}
	s.relay("", c.Name+" "+left)
	s.clients[idx] = s.clients[len(s.clients)-1]
	s.clients = s.clients[:len(s.clients)-1]
}

func (s ChatServer) relay(from string, msg string) {
	pre := svrLead
	if from != "" {
		pre = nOpen + from + nClose + " "
	}
	for _, c := range s.clients {
		if c.Name == from {
			continue
		}
		c.receive(pre + msg)
	}
	fmt.Println(pre + msg)
}

const (
	greeting = "Welcome to budgetchat! What shall I call you?"
	joined   = "joined"
	left     = "left"
	svrLead  = "* "
	nOpen    = "["
	nClose   = "]"
)

var (
	validNameRe = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
)

type client struct {
	conn net.Conn
	svr  *ChatServer
	Name string
}

func (c client) relay(msg string) {
	c.svr.relay(c.Name, msg)
}

func (c client) receive(msg string) {
	c.conn.Write([]byte(msg + "\n"))
}

func (c client) validName() bool {
	return validNameRe(c.Name)
}

func handle3(svr *ChatServer, conn net.Conn) {
	fmt.Println("open: " + conn.RemoteAddr().String())
	c := client{
		svr:  svr,
		conn: conn,
	}
	defer func() {
		conn.Close()
		if c.Name != "" {
			c.svr.leave(&c)
		}
	}()
	c.receive(greeting)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if c.Name == "" {
			c.Name = scanner.Text()
			if !c.validName() {
				return
			}
			c.svr.join(&c)
			continue
		}
		c.relay(scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Printf("3: err: %v\n", scanner.Err())
		return
	}
}
