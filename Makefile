all:
	go build -o cryptopener-cli cmd/main.go

gen-certs:
	openssl req  -nodes -new -x509 -keyout key.pem -out cert.pem -subj "/C=FI/ST=PP/L=Oulu/O=My Inc/OU=DevOps/CN=www.breach-test-server.com/emailAddress=test@www.breach-test-server.com"
