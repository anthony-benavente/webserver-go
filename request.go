package webserver

import (
	"bufio"
	"strings"
)

// Request The Request struct holds information about what the connected client
// sent to the web server.
type Request struct {
	action  string
	uri     string
	params  map[string]string
	form    map[string]string
	body    []byte
	headers map[string]string
}

func newRequest(action, uri string, body []byte) (results *Request) {
	results = &Request{
		action:  action,
		uri:     uri,
		body:    body,
		params:  make(map[string]string),
		form:    make(map[string]string),
		headers: make(map[string]string),
	}
	return
}

// NewRequest Creates a new request object which contains information the
// client that connected to the server sent.
func NewRequest(reader *bufio.Reader) (results *Request) {
	scanner := bufio.NewScanner(reader)
	lineIndex := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		if lineIndex > 0 {
			lineSplit := strings.Split(line, ":")
			if len(lineSplit) >= 2 {
				val := strings.Join(lineSplit[1:], ":")
				results.headers[lineSplit[0]] = val
			}
		} else {
			lineSplit := strings.Split(line, " ")
			if len(lineSplit) > 0 {
				results = newRequest(lineSplit[0], lineSplit[1], make([]byte, 512))
			} else {
				lineIndex--
			}
		}
		lineIndex++
	}
	return
}
