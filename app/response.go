package main

import "strconv"

func contentLengthHeader(body string) map[string]string {
	return map[string]string{contentLengthHeaderKey: strconv.Itoa(len(body))}
}

func contentEncodingHeader(encoding string) map[string]string {
	return map[string]string{contentEncodingHeaderKey: encoding}
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
