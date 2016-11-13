package main

import (
	webserver "github.com/anthony-benavente/webserver-go"
)

func main() {
	server := webserver.NewWebServer(":8080")

	server.Get("/", func(ctx *webserver.Context) {
		ctx.Res.Write("Hello, world!")
	})

	server.Get("/anthonys", func(ctx *webserver.Context) {
		ctx.Res.Write("My name is anthony!")
	})

	server.Listen()
}
