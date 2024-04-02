package main

import (
	"log"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

const (
	httpProtocol = "HTTP"
	httpVersion  = "1.1"
	httpStatusOK = "200 OK"
	doubleCrlf   = "\r\n\r\n"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer func() {
		err := l.Close()
		if err != nil {
			log.Println("Error closing listener: ", err.Error())
			return
		}
	}()

	conn, err := l.Accept()
	if err != nil {
		log.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection: ", err.Error())
			return
		}
	}()

	okResponse := httpProtocol + "/" + httpVersion + " " + httpStatusOK + " " + httpVersion + doubleCrlf
	_, err = conn.Write([]byte(okResponse))
	if err != nil {
		log.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
}
