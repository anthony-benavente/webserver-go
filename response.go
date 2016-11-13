package webserver

import (
	"bufio"
	"time"
)

// Response The response object contains information that will be sent back
// to the connected client.
type Response struct {
	writer  *bufio.Writer
	status  int
	message string
	body    []byte
	Headers map[string]string
}

// NewResponse Creates a new response that will be send back to the user
func NewResponse(status int, message string) (results *Response) {
	results = &Response{
		status:  status,
		message: message,
		body:    make([]byte, 0),
		Headers: make(map[string]string),
	}
	results.Headers["Date"] = time.Now().String()
	return
}

// Write Adds to what will be sent back to the connected client.
func (resp *Response) Write(message string) {
	resp.body = append(resp.body, []byte(message)...)
}
