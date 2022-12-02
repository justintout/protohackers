package protohackers

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

type SpeedServer struct {
	listener net.Listener
	messages chan speedMessage
	log      []speedMessage
}

func (s *SpeedServer) ListenAndServe(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s.listener = l
	defer l.Close()
	messages := make(chan speedMessage)
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

func (s *SpeedServer) Close() {
	s.listener.Close()
}

func (s *SpeedServer) listen(messages chan speedMessage) {
	for message := range messages {
		s.log = append(s.log, message)
		switch message.typ() {
		}
		s.handleMessage(message)
	}
}

func (s *SpeedServer) handleMessage(message speedMessage) {

}

const eof = 0x00

// client types
const (
	clientTypCamera     uint8 = 0x80
	clientTypDispatcher uint8 = 0x81
)

// message types
const (
	typEOF uint8 = 0x04

	// server -> client message types.
	typErr       uint8 = 0x10
	typTicket    uint8 = 0x21
	typHeartbeat uint8 = 0x41

	// client -> server message types
	typPlate         uint8 = 0x20
	typWantHeartbeat uint8 = 0x40
	typIAmCamera     uint8 = 0x80
	typIAmDispatcher uint8 = 0x81
)

type lstring struct {
	len uint8
	msg []byte
}

func (l *lstring) bytes() []byte {
	return append([]byte{l.len}, l.msg...)
}

type speedMessage interface {
	typ() uint8
}

type plateMessage struct {
	plate     []byte
	timestamp uint32
}

func (m *plateMessage) typ() uint8 {
	return typPlate
}

func (m *plateMessage) bytes() []byte {
	return append([]byte{m.typ()}, append(m.plate, []byte(m.timestamp))...)
}

type errorMessage struct {
	msg []byte
}

func typ() uint8 {
	return typErr
}

func newParser(conn net.Conn) (*parser, chan speedMessage) {
	p := parser{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
	go p.run()
	return &p, p.messages
}

type parser struct {
	conn       net.Conn
	reader     *bufio.Reader
	clientType uint8
	messages   chan speedMessage
	close      chan struct{}
}

// TODO: treating eof as an immediate hungup
func (p *parser) run() {
	for state := parseBegin; state != nil; {
		state = state(p)
	}
	close(p.messages)
}

func (p *parser) nextU8() (byte, bool) {
	var b byte
	err := binary.Read(p.conn, binary.BigEndian, &b)
	if err == io.EOF {
		return eof, true
	}
	if err != nil {
		panic(err)
	}
	return b, false
}

func (p *parser) nextU16() (uint16, bool) {
	var b uint16
	err := binary.Read(p.conn, binary.BigEndian, &b)
	if err == io.EOF {
		return b, true
	}
	if err != nil {
		panic(err)
	}
	return b, false
}

func (p *parser) nextU32() (uint32, bool) {
	var b uint32
	err := binary.Read(p.conn, binary.BigEndian, &b)
	if err == io.EOF {
		return b, true
	}
	if err != nil {
		panic(err)
	}
	return b, false
}

func (p *parser) nextString() (lstring, bool) {
	l, e := p.nextU8()
	if e {
		return lstring{len: 0, msg: []byte{eof}}, true
	}
	b := make([]byte, l)
	n, err := p.conn.Read(b)
	if err != nil && err != io.EOF {
		panic(err)
	}
	if err == io.EOF || len(b) != n {
		return lstring{len: uint8(n), msg: b[:n]}, true
	}
	return lstring{len: l, msg: b}, false
}

func (p *parser) nextN(n int) ([]byte, bool) {
	b := make([]byte, n)
	rn, err := p.reader.Read(b)
	if err != nil && err != io.EOF {
		panic(err)
	}
	if err == io.EOF {
		return b[:rn], true
	}
	return b[:rn], false
}

func (p *parser) emit(msg speedMessage) {
	p.messages <- msg
}

type stateFn func(*parser) stateFn

func parseBegin(p *parser) stateFn {
	for {
		next, eof := p.nextU8()
		if eof {
			break
		}
		if next == typErr || next == typTicket || next == typHeartbeat {
			p.emit(speedMessage{Typ: typErr, Data: []byte("illegal message for client")})
			break
		}
		if next == typPlate {
			if p.clientType != clientTypCamera {
				p.emit(speedMessage{Typ: typErr, Data: []byte("illegal message for client type")})
			}
			return parsePlate
		}
	}
	p.emit(speedMessage{Typ: typEOF})
	return nil
}

func parsePlate(p *parser) stateFn {
	var (
		plate     []byte
		timestamp uint32
		eof       bool
	)
	plate, eof = p.nextString()
	timestamp, eof = p.nextU32()
	p.emit(speedMessage{Typ: typPlate, D})
	if eof {
		return nil
	}
	return parseBegin
}
