package main

import (
	"fmt"
	cryptopener "github.com/larikuistio/cryptopener"
	"time"
)


func main() {
	fmt.Println("Not yet implemented")
	
	cryptopener := cryptopener.NewCryptopener("127.0.0.1:8080", "token=")
	go cryptopener.Run()
	for {
		time.Sleep(1 * time.Second)
	}

	return
}
