package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/illiafox/wireshark-dissector-example/proto"
)

// Client represents a TCP client
type Client struct {
	Addr   string
	conn   net.Conn
	closed bool

	ackCh chan time.Time

	logger slog.Logger
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := &Client{
		Addr:  addr,
		conn:  conn,
		ackCh: make(chan time.Time, 1),
	}

	go client.readPump()

	return client, nil
}

// Close closes the TCP connection
func (c *Client) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	close(c.ackCh)
	return c.conn.Close()
}

// SendMessage sends a SendMessage and waits for an acknowledgment (AckMessage)
func (c *Client) SendMessage(ctx context.Context, msg proto.SendMessage) error {
	data, err := proto.Marshall(msg)
	if err != nil {
		return fmt.Errorf("marshall message: %w", err)
	}

	if deadline, ok := ctx.Deadline(); ok {
		err = c.conn.SetWriteDeadline(deadline)
		if err != nil {
			return fmt.Errorf("set write deadline: %w", err)
		}
	}

	_, err = c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("write data: %w", err)
	}

	select {
	case ackReceived := <-c.ackCh:
		fmt.Println("Received ACK from server at", ackReceived)
	case <-ctx.Done():
		return fmt.Errorf("operation timed out waiting for ACK: %w", ctx.Err())
	}

	return nil
}

// readPump continuously listens for messages from the server
func (c *Client) readPump() {
	opBuf := make([]byte, 1) // 1 byte for AckMessage OpCode (ACK is 0x1)
	for {
		// Read data from the server
		n, err := c.conn.Read(opBuf)
		if err != nil {
			if err != net.ErrClosed && err != context.Canceled {
				log.Printf("Error reading from server: %v", err)
			}
			break
		}
		if n == 0 {
			break
		}

		opCode := proto.OpCode(opBuf[0])

		switch opCode {
		case proto.ACK:
			c.ackCh <- time.Now()
		default:
			log.Printf("Unexpected opCode from server: %x\n", opCode)
		}
	}
}
