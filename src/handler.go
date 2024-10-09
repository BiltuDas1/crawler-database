package db

import (
	"encoding/json"
	"log"
	"strconv"

	http "github.com/valyala/fasthttp"
)

// Function to handle every http/https Request
func Handler(ctx *http.RequestCtx) {
	// If the HTTP method is non GET then throws an Exception
	if !ctx.IsGet() {
		writeHTMLResp(ctx, StatusMethodNotAllowed, methodNotAllowed)
		return
	}

	// If the request is asking for to choose the database according to the query
	if string(ctx.URI().Path()) == "/v1/db/node/ip" {
		// If no hash parameter has been passed then show error to the client side
		if !ctx.URI().QueryArgs().Has("hash") || len(ctx.URI().QueryArgs().Peek("hash")) == 0 {
			writeJSONResp(ctx, StatusUnprocessableEntity, invalidInput)
			return
		}

		hashString := string(ctx.URI().QueryArgs().Peek("hash"))
		hash, err := strconv.ParseUint(hashString, 10, 64)
		if err != nil {
			writeJSONResp(ctx, StatusUnprocessableEntity, invalidInput)
			return
		}

		ip, port := GetDatabase(hash)
		resp, err := json.Marshal(Response{Result: true, Content: NodeDetails{IP: ip, Port: port}})

		// If JSON parse failed
		if err != nil {
			writeHTMLResp(ctx, StatusInternalError, internalError)
			log.Println(failedLog(err))
			return
		}

		// Send response when process successful
		writeJSONResp(ctx, StatusSuccessful, resp)
		return
	}

	// Generic handler
	writeHTMLResp(ctx, StatusUnauthorized, unauthorized)
}
