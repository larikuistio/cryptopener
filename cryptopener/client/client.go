package client

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"
)

const ConnectionTimeout = 60

type Client struct {
	config tls.Config
	Addr string
	bufferSize int
	// base attack vector
	Entrypoint string
}

func NewClient(addr string, entrypoint string) *Client {
	return &Client{
		Addr: addr,
		bufferSize: 4096,
		Entrypoint: entrypoint,
		config: tls.Config{
			// with this we can use our own cert
			InsecureSkipVerify: true,
		},
	}
}

func (client *Client) getRequestBody(message []byte, padding []byte) []byte {
	token := fmt.Sprintf("%s%s%s", string(padding), client.Entrypoint, string(message))
	request := fmt.Sprintf("POST / HTTP/1.1\r\nHost: %s\r\naccept-encoding: gzip, deflate, br\r\ncontent-length: %d\r\n\r\n%s\r\n", client.Addr, len(token), token)
	return []byte(request)
}

// Send a message into socket
func (client *Client) SendMessage(message []byte, padding []byte) []byte {
	connection, err := tls.Dial("tcp", client.Addr, &client.config)

	if err != nil {
		log.Printf("Failed to create connection, error %s", []byte(fmt.Sprintf("%s", err)))
		return nil
	}
	defer connection.Close()
	connection.SetReadDeadline(time.Now().Add(ConnectionTimeout * time.Second))
	m := client.getRequestBody(message, padding)
	_, err = connection.Write(m)
	if err != nil {
		log.Printf("Failed to write message, error %e", err)
		return nil
	}

	// buffer that contains whole message
	buffer := make([]byte, 0, client.bufferSize)
	for {
		b := make([]byte, 1056)
		size, eof := connection.Read(b)
		if size <= 1056 || eof != nil {
			buffer = append(buffer, b[:size]...)
			break
		}
		buffer = append(buffer, b[:size]...)
	}
	defer connection.Close()
	return buffer
}

// Communicates a return value into a channel thus it can used with runtime routines
func (client *Client) SendMessageConcurrent(channel chan []byte, message []byte, padding []byte) {
	channel <-client.SendMessage(message, padding)
	return
}
