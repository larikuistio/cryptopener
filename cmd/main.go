package main

import (
	"fmt"
	cryptopener "github.com/larikuistio/cryptopener"
	testserver "github.com/larikuistio/cryptopener/testserver"
)


func main() {
	fmt.Println("Not yet implemented")
	go testserver.TestServer()
	cryptopener := cryptopener.NewCryptopener("127.0.0.1:8080", "/")
	go cryptopener.Run()
	return
}
