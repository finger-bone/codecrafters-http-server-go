package main

import (
	"flag"
	"log"
	"net"
	"os"
)

func readConn(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func main() {
	directory := flag.String("directory", "", "Directory path")
	flag.Parse()

	log.Println("Directory:", *directory)

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

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			continue
		}

		go serve(*directory, conn)
	}
}

func serve(dir string, conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection: ", err.Error())
			return
		}
	}()

	req, err := readConn(conn)
	if err != nil {
		log.Println("Error reading connection: ", err.Error())
		return
	}

	startLine, headers, body := splitRequest(req)
	method, path, version := splitStartLine(startLine)

	handler(dir, method, path, version, headers, body, conn)
}
