package webserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// WebServer This struct is the web server that users can use listen for
// connections
type WebServer struct {
	host     string
	port     uint
	listener net.Listener
	headers  map[string]string
	handlers map[string]map[string][]func(*Context)
}

// Context This is a wrapper struct around a Request and a Response object.
// This is used for handling what happens when URIs are requested.
type Context struct {
	Req *Request
	Res *Response
}

// NewWebServer Creates a new WebServer with the specified address. The format
// for addresses is in the following <host>:<port>
func NewWebServer(addr string) (results *WebServer) {
	addrSplit := strings.Split(addr, ":")
	port, err := strconv.Atoi(addrSplit[1])
	if err != nil {
		port = 80
	}
	results = &WebServer{
		host:     addrSplit[0],
		port:     uint(port),
		listener: nil,
		headers:  make(map[string]string),
		handlers: make(map[string]map[string][]func(*Context)),
	}

	return
}

// NewContext Creates an object that wraps around a request and response object
func NewContext(req *Request, res *Response) (result *Context) {
	result = &Context{
		Req: req,
		Res: res,
	}
	return
}

// Get Adds a handler to the specified route with the given action
func (server *WebServer) Get(route string, action func(*Context)) {
	server.addHandler(route, "GET", action)
}

// Post Adds a handler for POST methods for the specified route
func (server *WebServer) Post(route string, action func(*Context)) {
	server.addHandler(route, "POST", action)
}

// Listen Starts the server listening at the port specified when the server
// was created.
func (server *WebServer) Listen() {
	ln, err := net.Listen("tcp", server.Addr())
	if err != nil {
		log.Fatal("Failed to initalize server at " + server.Addr())
	}
	server.listener = ln

	for {
		conn, _ := ln.Accept()
		go server.handleConnection(conn)
	}
}

// Addr Returns the address of the server as <host>:<port>
func (server *WebServer) Addr() (result string) {
	result = fmt.Sprintf("%v:%v", server.host, server.port)
	return
}

func (server *WebServer) handleRequest(res *Response, req *Request) {
	// Determine which route to get from the handlers
	if handler, valid := server.handlers[req.uri]; !valid {
		res.status = 404
		res.message = "NOT FOUND"
	} else {
		if methods, valid2 := handler[req.action]; !valid2 {
			res.status = 405
			res.message = "METHOD NOT ALLOWED"
		} else {
			for _, method := range methods {
				method(NewContext(req, res))
			}
		}

	}
}

func (server *WebServer) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	request := NewRequest(reader)
	response := NewResponse(200, "OK")

	server.handleRequest(response, request)
	server.writeResponse(writer, response)

	conn.Close()
}

func (server *WebServer) writeResponse(writer *bufio.Writer, response *Response) {
	writer.WriteString(fmt.Sprintf("\nHTTP/1.1 %v %v\n", response.status, response.message))
	for k := range response.Headers {
		writer.WriteString(fmt.Sprintf("%v: %v\n", k, response.Headers[k]))
	}
	writer.WriteString("\n")
	writer.WriteString(string(response.body))
	writer.WriteString("\n")
	writer.Flush()

}

func (server *WebServer) addHandler(route, method string, action func(*Context)) {
	routeMethods := server.handlers[route]
	if routeMethods == nil {
		server.handlers[route] = make(map[string][]func(*Context))
	}
	server.handlers[route][method] = append(server.handlers[route][method], action)
}

func readAll(reader *bufio.Reader) []byte {
	result := make([]byte, 1024)
	data := make([]byte, 1024)
	for n, err := reader.Read(result); n == 1024 && err == nil; {
		for v := range data {
			result = append(result, data[v])
		}
	}
	return result
}
