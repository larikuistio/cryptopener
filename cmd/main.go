package main

import (
	"fmt"
	cryptopener "github.com/larikuistio/cryptopener"
)


func main() {
	fmt.Println("Not yet implemented")
	cryptopener := cryptopener.NewCryptopener("127.0.0.1:8080", "token=")
	cryptopener.Run()
	defer func ()  {
		fmt.Printf("Found token, tokens is %s", cryptopener.ResultToken)
	}()
}
