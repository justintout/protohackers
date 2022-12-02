package speed

import (
	"encoding/binary"
	"fmt"
	"io"
)

// server -> client message types.
const (
	typErr       uint8 = 0x10
	typTicket    uint8 = 0x21
	typHeartbeat uint8 = 0x41
)

type ticketMessage struct {
	plate      lstring
	road       uint16
	mile1      uint16
	timestamp1 uint32
	mile2      uint16
	timestamp2 uint32
	speed      uint16
}

func (m ticketMessage) typ() uint8 {
	return typTicket
}

func (m ticketMessage) String() string {
	return fmt.Sprintf("ticket<plate=%s; road=%d; mile1=%d; timestamp1=%d; mile2=%d; timestamp2=%d; speed=%d>", m.plate, m.road, m.mile1, m.timestamp1, m.mile2, m.timestamp2, m.speed)
}

func (m ticketMessage) Write(w io.Writer) (int, error) {
	var n int
	nn, err := w.Write([]byte{m.typ()})
	n += nn
	if err != nil {
		return n, err
	}
	nn, err = m.plate.write(w)
	n += nn
	if err != nil {
		return n, err
	}
	b := make([]byte, 10)
	binary.BigEndian.PutUint16(b, m.road)
	binary.BigEndian.PutUint16(b, m.mile1)
	binary.BigEndian.PutUint32(b, m.timestamp1)
	binary.BigEndian.PutUint16(b, m.mile2)
	binary.BigEndian.PutUint32(b, m.timestamp2)
	binary.BigEndian.PutUint16(b, m.speed)
	nn, err = w.Write(b)
	n += nn
	if err != nil {
		return n, err
	}
	return n, nil
}

type errorMessage struct {
	msg []byte
}

func (m errorMessage) typ() uint8 {
	return typErr
}

func (m errorMessage) String() string {
	return fmt.Sprintf("error<msg=%s>", m.msg)
}

func (m errorMessage) Write(w io.Writer) (int, error) {
	n, err := w.Write([]byte{m.typ()})
	if err != nil {
		return n, err
	}
	nn, err := w.Write(m.msg)
	n += nn
	if err != nil {
		return n, err
	}
	return n, nil
}

type heartbeatMessage struct{}

func (m heartbeatMessage) typ() uint8 {
	return typHeartbeat
}
func (m heartbeatMessage) Write(w io.Writer) (int, error) {
	return w.Write([]byte{m.typ()})
}
