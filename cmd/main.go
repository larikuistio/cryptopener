package main

import (
	"fmt"
	cryptopener "github.com/larikuistio/cryptopener"
	testserver "github.com/larikuistio/cryptopener/testserver"
)


func main() {
	fmt.Println("Not yet implemented")
	cryptopener.NewCryptopener("127.0.0.1:8080", "/")
	testserver.TestServer()
	
	return
}
