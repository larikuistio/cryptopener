all:
	go build -o cryptopener-cli cmd/main.go
	go build -o testserver-cli testserver/testserver.go

gen-certs:
	openssl req  -nodes -new -x509 -keyout testserver/key.pem -out testserver/cert.pem -subj "/C=FI/ST=PP/L=Oulu/O=My Inc/OU=DevOps/CN=www.breach-test-server.com/emailAddress=test@www.breach-test-server.com"
