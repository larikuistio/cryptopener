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


func handleConnection(conn net.Conn) {
	
	buf := make([]byte, 0, 16384)
	
	_, err := conn.Read(buf)
	if err != nil {
	  	fmt.Println("Error reading:", err.Error())
	}
	
	conn2, err := net.Dial("tcp", HTTPSERVER_HOST + ":" + HTTPSERVER_PORT)
	if err != nil {
		fmt.Println("Error connecting to http server:", err.Error())
		os.Exit(1)
	}
	conn2.Write([]byte(buf))
	conn2.Read(buf)
	respLen := len(buf)
	fmt.Printf("Received response of length %d from http server\n", respLen)
	conn2.Close()
	conn.Write([]byte(buf))
	
	conn.Close()
}



func main() {
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
        go handleConnection(conn)
    }
}