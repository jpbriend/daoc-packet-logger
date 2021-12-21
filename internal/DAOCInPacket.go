package internal

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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

func (p *PacketLogger) parseDAOCInPacket(buf []byte) DAOCInPacket {
	len := len(buf)
	packet := DAOCInPacket{
		Size:        binary.BigEndian.Uint16(buf[0:2]),
		PacketCount: binary.BigEndian.Uint16(buf[2:4]),
		SessionID:   binary.BigEndian.Uint16(buf[4:6]),
		Parameter:   binary.BigEndian.Uint16(buf[6:8]),
		Code:        uint(buf[9]),
		Message:     buf[10 : len-2],
		Checksum:    binary.BigEndian.Uint16(buf[len-2 : len]),
	}

	fmt.Printf("IN/TCP - %v\n", packet.ToString())
	return packet
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
