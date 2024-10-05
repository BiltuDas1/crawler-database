package modules

import (
	"fmt"

	http "github.com/valyala/fasthttp"
)

// Function to handle every http Request
func Handler(ctx *http.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!")
}
