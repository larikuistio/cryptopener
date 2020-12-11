package main

import (
	"fmt"
	"net"
	"os"
	"io/ioutil"
)

const (
	LISTEN_HOST = "localhost"
	LISTEN_PORT = "50000"
	HTTPSERVER_HOST = "localhost"
	HTTPSERVER_PORT = "8080"
)

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

func main() {

	for {
		var buf [16384]byte

		tcpAddrc, err := net.ResolveTCPAddr("tcp", LISTEN_HOST + ":" + LISTEN_PORT)
		checkError(err)

		listener, err := net.ListenTCP("tcp", tcpAddrc)
		checkError(err)
		fmt.Println("Listening on " + LISTEN_HOST + ":" + LISTEN_PORT)

		tcpAddrs, err := net.ResolveTCPAddr("tcp", HTTPSERVER_HOST + ":" + HTTPSERVER_PORT)
		checkError(err)

		conns, err := net.DialTCP("tcp", nil, tcpAddrs)
		checkError(err)
		fmt.Println("Connected to http server at " + HTTPSERVER_HOST + ":" + HTTPSERVER_PORT)

		for {
			connc, err := listener.Accept()
			checkError(err)
			
			_, err = connc.Read(buf[0:])
			checkError(err)

			_, err = conns.Write([]byte(buf[0:]))
			checkError(err)

			response, err := ioutil.ReadAll(conns)
			checkError(err)
			respLen := len(response)
			fmt.Printf("Received response of length %d from http server\n", respLen)

			_, err = connc.Write([]byte(response))
			checkError(err)
		}
	}
}
