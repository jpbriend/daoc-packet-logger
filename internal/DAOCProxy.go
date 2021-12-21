package internal

import (
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

func (p *PacketLogger) processDAOCMessage(content []byte) {
	if p.Way == "IN" {
		p.parseDAOCInPacket(content)
	} else {
		p.parseDAOCOutPacket(content)
	}
}

// Connection //

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
		serverToClient := PacketLogger{Way: "OUT"}
		writer := io.MultiWriter(&serverToClient, conn)
		io.Copy(writer, backendConn)
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
