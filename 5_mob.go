package protohackers

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
)

const Addr5 = "0.0.0.0:10005"
const UpstreamChat = "chat.protohackers.com:16963"

type ChatProxyServer struct {
	listener     net.Listener
	upstreamAddr string
}

func (s *ChatProxyServer) ListenAndServe(addr, up string) {
	s.upstreamAddr = up
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		go handle5(s, conn)
	}
}

func (s *ChatProxyServer) Close() {
	s.listener.Close()
}

const tony = "${1}7YWHMfk9JZe0LM0g1ZauHuiSxhI${2}"

var replace = func(msg string) string {
	return regexp.MustCompile(`(\b)7[a-zA-Z0-9]{25,34}([ \n$])`).ReplaceAllString(msg, tony)
}

// TODO: these connections aren't managed super-duper well.
//       should synchronize readers with a channel or something.
func handle5(svr *ChatProxyServer, conn net.Conn) {
	upstream, err := net.Dial("tcp", svr.upstreamAddr)
	if err != nil {
		panic(err)
	}
	defer func() {
		conn.Close()
		upstream.Close()
	}()
	out := bufio.NewReader(conn)
	in := bufio.NewReader(upstream)
	open := true

	go func() {
		for open {
			t, err := out.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("5: outgoing reader err: %v\n", err)
				break
			}
			fmt.Printf("-> %s   %s", t, replace(t))
			_, err = fmt.Fprint(upstream, replace(t))
			if err != nil {
				fmt.Printf("5: outgoing write error: %v\n", err)
			}
		}
		fmt.Println("-/-> closed in outgoing, closing connections")
		open = false
		conn.Close()
		upstream.Close()
	}()

	for open {
		t, err := in.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("5: incoming scanner error: %v\n", err)
			break
		}
		fmt.Printf("<- %s   %s", t, replace(t))
		_, err = fmt.Fprint(conn, replace(t))
		if err != nil {
			fmt.Printf("5: incoming write error: %v\n", err)
		}
	}
	fmt.Println("<-/- closed in incoming, closing connections")
	open = false
}
