package internal

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var (
	rHost string
	rPort int
)

func Start(listenPort int, remoteHost string, remotePort int) error {
	rHost = remoteHost
	rPort = remotePort

	fmt.Printf("Starting proxy listening on port %v\n", listenPort)

	service, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", listenPort))
	listener, err := net.ListenTCP("tcp", service)

	if err != nil {
		fmt.Print("Error while starting the listener: ", err)
		return err
	}
	defer listener.Close()
	fmt.Println("Accepting connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return err
		}
		go handleConnection(conn)
	}
}

type PacketLogger struct {
	Way string
}

func (p *PacketLogger) Write(content []byte) (n int, err error) {
	p.processDAOCMessage(content)
	return len(content), nil
}

type DAOCPacket struct {
	Size        uint16
	PacketCount uint16
	SessionID   uint16
	Parameter   uint16
	Code        uint
	Message     []byte
	Checksum    uint16
}

func (p *PacketLogger) processDAOCMessage(content []byte) {
	// Remove first packet as it's always 0
	p.parsePacket(content[0:])
}

func (p *PacketLogger) parsePacket(buf []byte) DAOCPacket {
	len := len(buf)
	packet := DAOCPacket{
		Size:        binary.BigEndian.Uint16(buf[0:2]),
		PacketCount: binary.BigEndian.Uint16(buf[2:4]),
		SessionID:   binary.BigEndian.Uint16(buf[4:6]),
		Parameter:   binary.BigEndian.Uint16(buf[6:8]),
		Code:        uint(buf[9]),
		Message:     buf[10 : len-2],
		Checksum:    binary.BigEndian.Uint16(buf[len-2 : len]),
	}

	fmt.Printf("%s/TCP - %v\n", p.Way, packet.ToString())
	return packet
}

func (d *DAOCPacket) ToString() string {
	result := fmt.Sprintf("Size: %v #%v Code: 0x%X SessionID: 0x%X Param: 0x%X\n",
		d.Size,
		d.PacketCount,
		d.Code,
		d.SessionID,
		d.Parameter)

	result = result + hex.Dump(d.Message)
	return result
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("- Client detected -")
	fmt.Printf("Connecting to remote host %s on port %v\n\n", rHost, rPort)
	backendConn, err := net.Dial("tcp", net.JoinHostPort(rHost, fmt.Sprint(rPort)))
	if err != nil {
		log.Printf("Error while connecting to remote %s on port %v: %v\n", rHost, rPort, err)
		return
	}
	defer backendConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		//serverToClient := PacketLogger{Way: "OUT"}
		//writer := io.MultiWriter(&serverToClient, conn)
		io.Copy(conn, backendConn)
		conn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()
	go func() {
		clientToServer := PacketLogger{Way: "IN"}
		writer := io.MultiWriter(&clientToServer, backendConn)
		io.Copy(writer, conn)
		backendConn.(*net.TCPConn).CloseWrite()
		wg.Done()
	}()

	wg.Wait()

}
