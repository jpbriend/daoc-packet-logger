package internal

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

// IN Packet //

type DAOCInPacket struct {
	Size        uint16
	PacketCount uint16
	SessionID   uint16
	Parameter   uint16
	Code        uint
	Message     []byte
	Checksum    uint16
}

func (p *PacketLogger) parseDAOCInPacket(buf []byte) {
	packet := DAOCInPacket{
		Size:        binary.BigEndian.Uint16(buf[0:2]),
		PacketCount: binary.BigEndian.Uint16(buf[2:4]),
		SessionID:   binary.BigEndian.Uint16(buf[4:6]),
		Parameter:   binary.BigEndian.Uint16(buf[6:8]),
		Code:        uint(buf[9]),
		Message:     buf[10 : 10+binary.BigEndian.Uint16(buf[0:2])],
		Checksum:    binary.BigEndian.Uint16(buf[10+binary.BigEndian.Uint16(buf[0:2]) : 12+binary.BigEndian.Uint16(buf[0:2])]),
	}

	fmt.Printf("IN/TCP - Time: %v %v\n",
		time.Since(p.StartTime).Milliseconds(),
		packet.ToString())

	// if there is a following message, parse it
	remainingBuffer := buf[12+packet.Size:]

	if len(remainingBuffer) > 0 {
		p.parseDAOCInPacket(buf[12+binary.BigEndian.Uint16(buf[0:2]):])
	}
	return
}

func (d *DAOCInPacket) ToString() string {
	result := fmt.Sprintf("Size: %v #%v Code: 0x%X SessionID: 0x%X Param: 0x%X\n",
		d.Size,
		d.PacketCount,
		d.Code,
		d.SessionID,
		d.Parameter)

	result = result + hex.Dump(d.Message)
	return result
}
