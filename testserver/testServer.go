package main

import (
	"net/http"
	"log"
	"fmt"
	"math/rand"
	"time"
	"io/ioutil"
)

type testServer struct {
	
}

var token string

func api(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/text/plain")
	tokenJson := "token=" + token + string(body)
    switch r.Method {
	case "GET":
        w.WriteHeader(http.StatusOK)
		w.Write([]byte(tokenJson))
    default:
        w.WriteHeader(http.StatusNotFound)
    }
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().Unix())
	token = RandomString(32)
	fmt.Println("Token is: " + string(token))

    http.HandleFunc("/", api)
    log.Fatal(http.ListenAndServeTLS("127.0.0.1:8080", "cert.pem", "key.pem", nil))
}