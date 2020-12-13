package main

import (
	"fmt"
	"os"
	"net"
	"time"
	"runtime"
)

const (
	CLIENT_HOST = "127.0.0.1"
	CLIENT_PORT = "50000"
	HTTPSERVER_HOST = "127.0.0.1"
	HTTPSERVER_PORT = "8080"
)

func checkError(err error) {
    if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		_, file, no, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("\nin: %s#%d\n", file, no)
		}
		fmt.Println("\n")
        os.Exit(1)
    }
}

func clientSocket(ctoschan chan []byte, stocchan chan []byte) {
	
	netaddr, err := net.ResolveIPAddr("ip4:ip", CLIENT_HOST)
	checkError(err)

	conn, err := net.ListenIP("ip4:ipv4", netaddr)
	checkError(err)
	
	go clientSocketWriter(conn, stocchan)
	go clientSocketReader(conn, ctoschan)

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

func clientSocketWriter(conn *net.IPConn, stocchan chan []byte) {
	buf := make([]byte, 16384)
	for {
		buf = <- stocchan
		_, err := conn.Write(buf)
		checkError(err)
		fmt.Println("clientSocketWriter")
	}
}

func clientSocketReader(conn *net.IPConn, ctoschan chan []byte) {
	buf := make([]byte, 16384)
	for {
		_, err := conn.Read(buf)
		checkError(err)
		fmt.Println("clientSocketReader")
		ctoschan <- buf
	}
}

func serverSocket(ctoschan chan []byte, stocchan chan []byte) {

	netaddr, err := net.ResolveIPAddr("ip4:ip", HTTPSERVER_HOST)
	checkError(err)

	conn, err := net.DialIP("ip4:ipv4", nil, netaddr)
	checkError(err)

	go serverSocketWriter(conn, ctoschan)
	go serverSocketReader(conn, stocchan)

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

func serverSocketWriter(conn *net.IPConn, ctoschan chan []byte) {
	buf := make([]byte, 16384)
	for {
		buf = <- ctoschan
		_, err := conn.Write(buf)
		checkError(err)
		fmt.Println("serverSocketWriter")
	}
}

func serverSocketReader(conn *net.IPConn, stocchan chan []byte) {
	buf := make([]byte, 16384)
	for {
		_, err := conn.Read(buf)
		checkError(err)
		fmt.Println("serverSocketReader")
		stocchan <- buf
	}
}

func main() {

	ctoschan := make(chan []byte)
	stocchan := make(chan []byte)

	go clientSocket(ctoschan, stocchan)
	go serverSocket(ctoschan, stocchan)

	for {
		time.Sleep(1000 * time.Millisecond)
	}

}
