package main

import (
	"fmt"
	"net"
	"os"
)

const (
	LISTEN_HOST = "localhost"
	LISTEN_PORT = "50000"
	HTTPSERVER_HOST = "localhost"
	HTTPSERVER_PORT = "8080"
)


func handleConnFromClient(conn net.Conn, reqc chan []byte, respc chan []byte) {
	
	buf := make([]byte, 0, 16384)
	
	_, err := conn.Read(buf)
	if err != nil {
	  	fmt.Println("Error reading:", err.Error())
	}
	reqc <- buf
	buf = <- respc
	conn.Write([]byte(buf))
	conn.Close()
}

func handleConnToServer(reqc chan []byte, respc chan []byte) {
	
	buf := make([]byte, 0, 16384)
	
	for {
		buf = <- reqc
		conn, err := net.Dial("tcp", HTTPSERVER_HOST + ":" + HTTPSERVER_PORT)
		if err != nil {
			fmt.Println("Error connecting to http server:", err.Error())
			os.Exit(1)
		}
		conn.Write([]byte(buf))
		conn.Read(buf)
		respLen := len(buf)
		fmt.Printf("Received response of length %d from http server\n", respLen)
		respc <- buf
		
		conn.Close()
	}
}

func main() {

	reqc := make(chan []byte)
	respc := make(chan []byte)

    l, err := net.Listen("tcp", LISTEN_HOST + ":" + LISTEN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    defer l.Close()
    fmt.Println("Listening on " + LISTEN_HOST + ":" + LISTEN_PORT)

	for {        
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
		go handleConnFromClient(conn, reqc, respc)
		go handleConnToServer(reqc, respc)
    }
}