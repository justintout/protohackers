package speed

import (
	"fmt"
	"strings"
)

// client -> server message types
const (
	typEOF uint8 = 0x04

	typPlate         uint8 = 0x20
	typWantHeartbeat uint8 = 0x40
	typIAmCamera     uint8 = 0x80
	typIAmDispatcher uint8 = 0x81
)

type eofMessage struct{}

func (m eofMessage) typ() uint8 {
	return typErr
}

func (m eofMessage) String() string {
	return "EOF"
}

type plateMessage struct {
	plate     lstring
	timestamp uint32
}

func (m plateMessage) typ() uint8 {
	return typPlate
}

func (m plateMessage) String() string {
	return fmt.Sprintf("plate<plate=%s; timestamp=%d>", m.plate, m.timestamp)
}

type wantHearbeatMessage struct {
	interval uint32
}

func (m wantHearbeatMessage) String() string {
	return fmt.Sprintf("wantHeartbeat<interval=%d>", m.interval)
}

func (m wantHearbeatMessage) typ() byte {
	return typWantHeartbeat
}

type iAmCameraMessage struct {
	road  uint16
	mile  uint16
	limit uint16
}

func (m iAmCameraMessage) String() string {
	return fmt.Sprintf("iAmCamera<road=%d; mile=%d; limit=%d>", m.road, m.mile, m.limit)
}

func (m iAmCameraMessage) typ() uint8 {
	return typIAmCamera
}

type iAmDispatcherMessage struct {
	numroads uint8
	roads    []uint16
}

func (m iAmDispatcherMessage) String() string {
	r := make([]string, len(m.roads))
	for i, road := range m.roads {
		r[i] = fmt.Sprintf("%d", road)
	}
	return fmt.Sprintf("iAmDispatcher<numroads=%d; roads=%s>", m.numroads, strings.Join(r, ","))
}

func (m iAmDispatcherMessage) typ() uint8 {
	return typIAmDispatcher
}
