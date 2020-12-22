all:
	go build -o cryptopener-cli cmd/main.go

gen-certs:
	openssl req -x509 -newkey rsa:4096 -keyout server-key.pem -out server-cert.pem -days 365
