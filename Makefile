all:
	go build -o cryptopener-cli cmd/main.go

gen-certs:
	openssl req  -nodes -new -x509 -keyout testserver/key.pem -out testserver/cert.pem -days 365
