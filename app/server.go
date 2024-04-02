package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const okResponseHead = "HTTP/1.1 200 OK"
const crlf = "\r\n"
const notFoundResponseHead = "HTTP/1.1 404 Not Found"

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

func splitRequest(req string) (string, map[string]string, string) {
	splitted := strings.Split(req, crlf+crlf)
	startLine := strings.Split(splitted[0], crlf)[0]
	headers := strings.Split(splitted[0], crlf)[1:]

	headersMap := make(map[string]string)
	for _, header := range headers {
		splitted := strings.Split(header, ": ")
		headersMap[splitted[0]] = splitted[1]
	}

	body := splitted[1]
	return startLine, headersMap, body
}

func splitStartLine(startLine string) (string, string, string) {
	splitted := strings.Split(startLine, " ")
	return splitted[0], splitted[1], splitted[2]
}

func readConn(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func handler(dir, method, path, version string, headers map[string]string, body string, conn net.Conn) {
	log.Println("Method: ", method)
	log.Println("Path: ", path)
	log.Println("Version: ", version)
	log.Println("Headers: ", headers)
	log.Println("Body: ", body)

	if path == "/user-agent" {
		responseBody := headers["User-Agent"]
		conn.Write(buildResponse(
			okResponseHead,
			mergeMaps(contentLengthHeader(responseBody), map[string]string{"Content-Type": "text/plain"}),
			responseBody,
		))
		return
	}

	if strings.HasPrefix(path, "/echo") {
		responseBody := path[6:]
		conn.Write(buildResponse(
			okResponseHead,
			mergeMaps(contentLengthHeader(responseBody), map[string]string{"Content-Type": "text/plain"}),
			responseBody,
		))
		return
	}

	if path == "/" {
		conn.Write(buildResponse(okResponseHead, nil, ""))
		return
	}

	if strings.HasPrefix(path, "/files") {
		filename := path[7:]
		file, err := os.Open(dir + "/" + filename)
		if err != nil {
			conn.Write(buildResponse(notFoundResponseHead, nil, ""))
			return
		}
		responseBody := make([]byte, 1024)
		// read the file
		n, err := file.Read(responseBody)
		if err != nil {
			conn.Write(buildResponse(notFoundResponseHead, nil, ""))
			return
		}
		conn.Write(buildResponse(
			okResponseHead,
			mergeMaps(contentLengthHeader(string(responseBody[:n])), map[string]string{"Content-Type": "text/plain"}),
			string(responseBody[:n]),
		))
		return
	}

	conn.Write(buildResponse(notFoundResponseHead, nil, ""))
}

func contentLengthHeader(body string) map[string]string {
	return map[string]string{"Content-Length": strconv.Itoa(len(body))}
}

func buildResponse(
	statusText string,
	headers map[string]string,
	body string,
) []byte {
	response := statusText + crlf
	for key, value := range headers {
		response += key + ": " + value + crlf
	}
	response += crlf
	if body != "" {
		response += body
	}
	return []byte(response)
}

func mergeMaps(map1, map2 map[string]string) map[string]string {
	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}
