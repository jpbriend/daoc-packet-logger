package internal

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

// OUT Packet //

type DAOCOutPacket struct {
	Size    uint16
	Code    uint
	Message []byte
}

func (p *PacketLogger) parseDAOCOutPacket(buf []byte) {
	packet := DAOCOutPacket{
		Size:    binary.BigEndian.Uint16(buf[0:2]),
		Code:    uint(buf[2]),
		Message: buf[3 : 3+binary.BigEndian.Uint16(buf[0:2])],
	}

	fmt.Printf("OUT/TCP - Time: %v %v\n",
		time.Since(p.StartTime).Milliseconds(),
		packet.ToString())

	// if there is a following message, parse it
	remainingBuffer := buf[3+packet.Size:]

	if len(remainingBuffer) > 0 {
		p.parseDAOCOutPacket(buf[3+binary.BigEndian.Uint16(buf[0:2]):])
	}

	return
}

func (d *DAOCOutPacket) ToString() string {
	result := fmt.Sprintf("Size: %v Code: 0x%X\n",
		d.Size,
		d.Code)

	result = result + hex.Dump(d.Message)
	return result
}
