package main

import (
	"log"
	"net/http"
	"io"
	"os"
	"crypto/tls"
)


type Client struct {
	config tls.Config
	addr string
	bufferSize int
}



func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
		bufferSize: 4096,
		config: tls.Config{
			// with this we can use our own cert
			InsecureSkipVerify: true,
		},
	}
}

func (client *Client) SendMessage(message []byte) ([]byte, error) {
	connection, err := tls.Dial("tcp", client.addr, &client.config)
	log.Printf("Sendig message %b to server", message)
	if err != nil {
		log.Fatalf("Failed to create connection, error %e", err)
		return nil, err
	}
	defer connection.Close()

	_, err = connection.Write(message)
	if err != nil {
		log.Fatalln("Failed to write message")
		return nil, err
	}

	// buffer that contains whole message
	buffer := make([]byte, 0, client.bufferSize)
	for {
		log.Printf("Reading bytes to buffer %b", buffer)
		b := make([]byte, 1056)
		size, eof := connection.Read(b)
		if eof != nil {
            break
		}
		buffer = append(buffer, b[:size]...)
	}
	log.Printf("Return buffer %b", buffer)
	_, e := connection.Read(buffer)
	if e != nil {
		log.Fatalf("error %e while reading buffer", e)
	}
	return buffer, nil
}

func main() {
	client := NewClient("127.0.0.1:8080")
	log.Println(client.SendMessage([]byte("test")))
	return
}
