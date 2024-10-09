package db

import (
	color "github.com/fatih/color"
	http "github.com/valyala/fasthttp"
)

type Response struct {
	Result      bool        `json:"result"`
	Description interface{} `json:"description,omitempty"`
	Content     interface{} `json:"content,omitempty"`
}

type NodeDetails struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

// Function for sending HTML data to the user side
func writeHTMLResp(ctx *http.RequestCtx, statusCode int, body []byte) {
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.SetBody(body)
}

// Function for sending JSON Data to the user side
func writeJSONResp(ctx *http.RequestCtx, statusCode int, body []byte) {
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetBody(body)
}

var methodNotAllowed = []byte("<center><h1>405 Not Allowed</h1><hr /><p>crawler-db</p></center>")
var unauthorized = []byte("<center><h1>401 Authorization Required</h1><hr /><p>crawler-db</p></center>")
var internalError = []byte("<center><h1>500 Internal Server Error</h1><hr /><p>crawler-db</p></center>")
var invalidInput = []byte(`{"result": false, "description": "hash parameter is empty or invalid"}`)
var failedLog = color.New(color.FgRed).SprintFunc()

const (
	StatusMethodNotAllowed    = 405
	StatusUnauthorized        = 401
	StatusInternalError       = 500
	StatusUnprocessableEntity = 422
	StatusSuccessful          = 200
)
