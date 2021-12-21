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
		Size:        binary.BigEndian.Uint16(buf[0:4]),
		PacketCount: binary.BigEndian.Uint16(buf[4:6]),
		SessionID:   binary.BigEndian.Uint16(buf[6:8]),
		Parameter:   binary.BigEndian.Uint16(buf[8:10]),
		Code:        uint(buf[11]),
		Message:     buf[12 : len-2],
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
