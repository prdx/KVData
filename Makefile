CC=gcc -g -O3
GO=go build
all: 	clean run
default: clean run
clean:
	rm -rf proxy server
proxy:	proxy.go
	$(GO) proxy.go
server:
	$(GO) server.go
run: proxy server

