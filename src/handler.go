package modules

import (
	"bytes"
	"log"

	color "github.com/fatih/color"
	http "github.com/valyala/fasthttp"
)

var methodNotAllowed = []byte("<center><h1>405 Not Allowed</h1><hr /><p>crawler-db</p></center>")
var unauthorized = []byte("<center><h1>401 Authorization Required</h1><hr /><p>crawler-db</p></center>")
var failed = color.New(color.FgRed).SprintFunc()

// Function to handle every http Request
func Handler(ctx *http.RequestCtx) {
	// If the HTTP method is non GET then throws an Exception
	if !bytes.Equal(ctx.Method(), []byte("GET")) {
		// Shows Error to the user side
		ctx.SetStatusCode(405) // Method not allowed
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.SetBody(methodNotAllowed)
		return
	}

	// If the request is asking for to choose the database according to the query
	if bytes.Equal(ctx.URI().Path(), []byte("/v1/db/node/ip")) {
		get_database(ctx)
		return
	}

	// Generic handler
	ctx.SetStatusCode(401) // Unauthorized
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.SetBody(unauthorized)
}

// Function to find the database
func get_database(ctx *http.RequestCtx) {
	_, err := ctx.WriteString("Access Granted")
	if err != nil {
		log.Fatalln(failed(err))
	}
}
