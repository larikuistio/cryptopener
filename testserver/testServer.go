package main

import (
	"net/http"
	"log"
	"fmt"
	"math/rand"
	"time"
	"io/ioutil"
	"compress/gzip"
	"strings"
	"io"
)

var token string

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "text/plain")
	response := "token=" + token + string(body)
    switch r.Method {
	case "GET":
		fmt.Println("GET")
        w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
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
	token = RandomString(64)
	fmt.Println("Token is: " + string(token))

    log.Fatal(http.ListenAndServeTLS("127.0.0.1:8080", "cert.pem", "key.pem", makeGzipHandler(handler)))
}