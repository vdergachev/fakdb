package main

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}

func (c *Connection) Init() error {

	timeout := time.Duration(30) * time.Second
	deadline := time.Now().Add(timeout)

	if err := c.conn.SetDeadline(deadline); err != nil {

		// TODO: Send error to the socket before close

		_ = c.conn.Close() // TODO: Close in other routine ???

		return err
	}
	return nil
}

func (c *Connection) Handle() {

	for {
		_, buf, err := c.receive()
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by client: ", c.conn.RemoteAddr())
			} else {
				fmt.Println("Error reading from connection: ", err.Error())
				_ = c.conn.Close()
			}
			return
		}
		var op Operation
		if err := op.Parse(buf); err != nil {
			_, _ = c.conn.Write([]byte("ERR " + err.Error() + "\n"))
			continue
		}

		res, err := op.Execute()
		if err != nil {
			fmt.Println("Op exec error: ", err.Error()) // TODO: Debug log + metrics
		}
		switch res.Code {
		case ORC_OK:
			_, _ = c.conn.Write([]byte("OK\n"))
			continue
		case ORC_ERR:
			_, _ = c.conn.Write([]byte("ERR " + res.Body + "\n"))
			continue
		case ORC_INT:
			_, _ = c.conn.Write([]byte("BYE\n"))
			_ = c.conn.Close()
			return
		}
	}

}

func (c *Connection) receive() (readBytes int, buf []byte, err error) {
	var received int
	// The buffer grows as we write into it.
	// Ref: https://pkg.go.dev/bytes#Buffer
	buffer := bytes.NewBuffer(nil)
	// Read the data in chunks.
	for {
		// c.ReceiveChunkSize = 8192
		chunk := make([]byte, 1024)
		read, err := c.conn.Read(chunk)
		if err != nil {
			return received, buffer.Bytes(), err
		}
		received += read
		buffer.Write(chunk[:read])

		if read == 0 || read < 1024 {
			break
		}
	}
	return received, buffer.Bytes(), nil

}
