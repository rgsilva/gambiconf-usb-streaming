package stream

import (
	"bufio"
	"fmt"
	"net"
)

type (
	Receiver interface {
		Start() error
		Stop()
	}

	receiver struct {
		port     uint16
		active   bool
		bq       *BufferQueue
		listener net.Listener
	}
)

func NewReceiver(bq *BufferQueue, port uint16) Receiver {
	return &receiver{bq: bq, port: port}
}

func (r *receiver) Start() error {
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}
	r.listener = ln
	r.active = true

	go func() {
		for {
			if !r.active {
				return
			}

			println("[RECV] Waiting for connections...")

			// Accept a new connection
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("[RECV] Error accepting connection:", err)
				continue
			}
			println("[RECV] Connection accepted.")

			reader := bufio.NewReader(conn)

			// Handle the connection
			for {
				if !r.active {
					return
				}

				buffer := make([]byte, 8192)
				n, err := reader.Read(buffer)
				if err != nil {
					fmt.Println("[RECV] Error reading from connection:", err)
					break
				}

				// Queue the data
				tmp := make([]byte, n)
				copy(tmp, buffer[:n])
				r.bq.Put(tmp)
			}

			conn.Close()
		}
	}()
	return nil
}

func (r *receiver) Stop() {
	r.active = false
	r.listener.Close()
}
