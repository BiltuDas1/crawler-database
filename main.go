package main

import (
	"log"
	"net"
	"os"

	db "crawler-db/src"

	color "github.com/fatih/color"
	http "github.com/valyala/fasthttp"
)

func main() {
	failed := color.New(color.FgRed).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()

	// Initialize DB
	db.InitializeDBList("db.txt")

	listener, err := net.Listen("tcp4", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalf(failed("Error in net.Listen: ", err))
	}

	// Start the server with default settings.
	// Create Server instance for adjusting server settings.
	//
	// Serve returns on ln.Close() or error, so usually it blocks forever.
	log.Println(success("Server is running on: ", listener.Addr()))

	if err := http.Serve(listener, db.Handler); err != nil {
		log.Fatalf(failed("Error in Serve: ", err))
	}
}
