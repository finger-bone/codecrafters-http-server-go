package main

import (
	"strings"
)

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
