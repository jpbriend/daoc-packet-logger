package internal

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// OUT Packet //

type DAOCOutPacket struct {
	Size    uint16
	Code    uint
	Message []byte
}

func (p *PacketLogger) parseDAOCOutPacket(buf []byte) DAOCOutPacket {
	len := len(buf)
	packet := DAOCOutPacket{
		Size:    binary.BigEndian.Uint16(buf[0:2]),
		Code:    uint(buf[2]),
		Message: buf[3:len],
	}

	fmt.Printf("OUT/TCP - %v\n", packet.ToString())
	return packet
}

func (d *DAOCOutPacket) ToString() string {
	result := fmt.Sprintf("Size: %v Code: 0x%X\n",
		d.Size,
		d.Code)

	result = result + hex.Dump(d.Message)
	return result
}
