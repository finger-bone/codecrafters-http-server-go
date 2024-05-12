package main

import (
	"bytes"
	"compress/gzip"
	"log"
	"net"
	"os"
	"strings"
)

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
		responseHeaders := make(map[string]string)
		responseHeaders = mergeMaps(responseHeaders, contentLengthHeader(responseBody))
		responseHeaders = mergeMaps(responseHeaders, map[string]string{"Content-Type": "text/plain"})

		if enc, ok := headers["Accept-Encoding"]; ok {
			allEncodings := splitArray(enc)
			for _, encoding := range allEncodings {
				if encoding == "gzip" {
					responseHeaders = mergeMaps(responseHeaders, contentEncodingHeader("gzip"))
					// use gzip to compress the responseBody
					var buf bytes.Buffer
					zw := gzip.NewWriter(&buf)
					_, _ = zw.Write([]byte(responseBody))
					zw.Close()
					responseBody = buf.String()
					break
				}
			}
		}
		conn.Write(buildResponse(
			okResponseHead,
			responseHeaders,
			responseBody,
		))
		return
	}

	if path == "/" {
		conn.Write(buildResponse(okResponseHead, nil, ""))
		return
	}

	if strings.HasPrefix(path, "/files") && method == getMethod {
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
			mergeMaps(contentLengthHeader(string(responseBody[:n])), map[string]string{"Content-Type": "application/octet-stream"}),
			string(responseBody[:n]),
		))
		return
	}

	if strings.HasPrefix(path, "/files") && method == postMethod {
		filename := path[7:]
		file, err := os.Create(dir + "/" + filename)
		if err != nil {
			conn.Write(buildResponse(notFoundResponseHead, nil, ""))
			return
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Println("Error closing file: ", err.Error())
				return
			}
		}()

		_, err = file.Write([]byte(body))
		if err != nil {
			conn.Write(buildResponse(notFoundResponseHead, nil, ""))
			return
		}

		conn.Write(buildResponse(createdResponseHead, nil, ""))
		return
	}

	conn.Write(buildResponse(notFoundResponseHead, nil, ""))
}
