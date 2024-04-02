package main

import (
	"log"
	"strconv"
	"strings"

	"net"
	"os"
)

const okResponseHead = "HTTP/1.1 200 OK"
const crlf = "\r\n"
const notFoundResponseHead = "HTTP/1.1 404 Not Found"

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

	req, err := readConn(conn)
	if err != nil {
		log.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}

	splitted := strings.Split(req, crlf)
	startLine := strings.Split(splitted[0], " ")
	_, path, _ := startLine[0], startLine[1], startLine[2]

	splittedPath := splitPath(path)
	log.Println(splittedPath)

	if splittedPath[0] == "echo" {
		contentLength := strconv.Itoa(len(splittedPath[1]))
		writeConn(conn,
			okResponseHead+crlf+
				"Content-Type: text/plain"+crlf+
				"Content-Length: "+contentLength+crlf+crlf+
				splittedPath[1])
		return
	}

	if path == "/" {
		writeConn(conn, okResponseHead+crlf+crlf)
		return
	} else {
		writeConn(conn, notFoundResponseHead+crlf+crlf)
		return
	}

}

func readConn(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func writeConn(conn net.Conn, resp string) error {
	_, err := conn.Write([]byte(resp))
	if err != nil {
		return err
	}

	return nil
}

func splitPath(path string) []string {
	// splitPath only into 2 parts
	// /echo/foo/bar to [echo, foo/bar]
	ret := make([]string, 2)
	splitted := strings.Split(path, "/")
	ret[0] = splitted[1]
	ret[1] = strings.Join(splitted[2:], "/")
	return ret
}
