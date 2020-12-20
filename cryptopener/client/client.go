package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"
)

const ConnectionTimeout = 60

type Client struct {
	config tls.Config
	addr string
	bufferSize int
	// base attack vector
	entrypoint string
}



func NewClient(addr string, entrypoint string) *Client {
	return &Client{
		addr: addr,
		bufferSize: 4096,
		entrypoint: entrypoint,
		config: tls.Config{
			// with this we can use our own cert
			InsecureSkipVerify: true,
		},
	}
}

func (client *Client) getRequestBody(message string) []byte {
	token := fmt.Sprintf("%s%s", client.entrypoint, message)
	request := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nAccept-Encoding: gzip;deflate\r\n\r\n", token, client.addr)
	return []byte(request)
}

func (client *Client) SendMessage(message string) ([]byte, error) {
	connection, err := tls.Dial("tcp", client.addr, &client.config)
	// set timeout for connection
	connection.SetReadDeadline(time.Now().Add(ConnectionTimeout * time.Second))

	log.Printf("Sendig message %s to server", message)
	if err != nil {
		log.Fatalf("Failed to create connection, error %e", err)
		return nil, err
	}
	defer connection.Close()

	m := client.getRequestBody(message)
	_, err = connection.Write(m)
	if err != nil {
		log.Fatalln("Failed to write message")
		return nil, err
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
	log.Printf("Return buffer %s", string(buffer))
	return buffer, nil
}
