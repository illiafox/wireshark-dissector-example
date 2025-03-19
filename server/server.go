package main

import (
	"context"
	"fmt"
	"github.com/illiafox/wireshark-dissector-example/proto"
	"log"
	"net"
)

type Server struct {
	Addr   string
	ln     net.Listener
	closed bool
}

func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on address %s: %w", addr, err)
	}
	server := &Server{
		Addr: addr,
		ln:   ln,
	}

	go server.listenForClients()

	return server, nil
}

// Close closes the server listener
func (s *Server) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	return s.ln.Close()
}

func (s *Server) listenForClients() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			return
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	opBuf := make([]byte, 1) // 1 byte for AckMessage OpCode (ACK is 0x1)
	for {
		// Read data from the server
		n, err := conn.Read(opBuf)
		if err != nil {
			if err != net.ErrClosed && err != context.Canceled {
				log.Printf("Error reading from client: %v", err)
			}
			break
		}
		if n == 0 {
			break
		}

		opCode := proto.OpCode(opBuf[0])
		switch opCode {
		case proto.SEND:
			var msg proto.SendMessage
			err = msg.ReadBinary(conn)
			if err != nil {
				log.Printf("Error reading from client: %v", err)
				continue
			}

			fmt.Printf("Received message: %+v\n", msg)

			out, err := proto.Marshall(proto.AckMessage{})
			if err != nil {
				log.Printf("Error marshalling ack message: %v", err)
				continue
			}

			_, err = conn.Write(out)
			if err != nil {
				log.Printf("Error writing ACK to client: %v", err)
			}

		default:
			log.Printf("Unexpected opCode from server: %x\n", opCode)
		}
	}
}

func main() {
	// Create a new server and start listening on the desired address
	server, err := NewServer("localhost:12345") // Change this to your desired server address
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	log.Println("Server is listening on localhost:12345")
	// Block the main goroutine to keep the server running
	select {}
}
