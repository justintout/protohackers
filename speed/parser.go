package speed

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type lstring struct {
	len uint8
	msg []byte
}

func (l lstring) bytes() []byte {
	return append([]byte{l.len}, l.msg...)
}

func (l lstring) String() string {
	return fmt.Sprintf("str<len=%d; msg=%s>", l.len, l.msg)
}

func (l lstring) write(w io.Writer) (int, error) {
	return w.Write(append([]byte{l.len}, l.msg...))
}

func newParser(conn net.Conn) (*parser, chan message) {
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
	heartbeat  *time.Ticker
	messages   chan message
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

func (p *parser) emit(msg message) {
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
			p.emit(errorMessage{msg: []byte("illegal message for client")})
			break
		}
		if next == typPlate {
			if p.clientType != clientTypCamera {
				p.emit(errorMessage{msg: []byte("illegal message for client type")})
			}
			return parsePlate
		}
	}
	p.emit(eofMessage{})
	return nil
}

func parsePlate(p *parser) stateFn {
	var (
		plate     lstring
		timestamp uint32
		eof       bool
	)
	plate, eof = p.nextString()
	timestamp, eof = p.nextU32()
	p.emit(&plateMessage{
		plate:     plate,
		timestamp: timestamp,
	})
	if eof {
		return nil
	}
	return parseBegin
}

func parseWantHeartbeat(p *parser) stateFn {
	var (
		interval uint32
		eof      bool
	)
	interval, eof = p.nextU32()
	p.emit(wantHearbeatMessage{
		interval: interval,
	})
	if eof {
		return nil
	}
	return parseBegin
}

func parseIAmCamera(p *parser) stateFn {
	p.emit(&errorMessage{
		msg: []byte("client already declared type"),
	})
	var (
		road  uint16
		mile  uint16
		limit uint16
		eof   bool
	)
	road, eof = p.nextU16()
	mile, eof = p.nextU16()
	limit, eof = p.nextU16()
	p.emit(iAmCameraMessage{
		road:  road,
		mile:  mile,
		limit: limit,
	})
	if eof {
		return nil
	}
	return parseBegin
}
