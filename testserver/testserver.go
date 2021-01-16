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
	"os"
	"runtime"
	"path/filepath"
	"flag"
)

var token string

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

type TestServer struct {
	Token string
}

func NewTestServer(length int) *TestServer {
	rand.Seed(time.Now().UTC().UnixNano())
	return &TestServer{
		Token: RandomString(length),
	}
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

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
	reader := r.Body
	body, err := ioutil.ReadAll(reader)
	checkError(err)
	w.Header().Set("Content-Type", "text/plain")
	response := "token=" + token + string(body)
    switch r.Method {
	case "GET":
		fmt.Println("testserver: received GET")
		fmt.Println("testserver: request body: " + string(body))
		fmt.Println("testserver: sending response with body: " + response)
        w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	case "POST":
		fmt.Println("testserver: received POST")
		fmt.Println("testserver: request body: " + string(body))
		fmt.Println("testserver: sending response with body: " + response)
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

	var length int
	flag.IntVar(&length, "length", 64, "length of secret token")
	flag.Parse()

	rand.Seed(time.Now().Unix())
	s := NewTestServer(length)
	token = s.Token
	fmt.Println("testserver: token=" + string(token))
	fmt.Println("testserver: starting server on 127.0.0.1:8080")
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
    log.Fatal(http.ListenAndServeTLS("127.0.0.1:8080", dir + "/cert.pem", dir + "/key.pem", makeGzipHandler(handler)))
}
