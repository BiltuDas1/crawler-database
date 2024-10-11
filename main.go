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
		if _, docker := os.LookupEnv("RunInDocker"); docker {
			log.Fatalln("\033[31mError in net.Listen: " + err.Error() + "\033[0m")
		}
		log.Fatalln(failed("Error in net.Listen: ", err))
	}

	// Start the server with default settings.
	// Create Server instance for adjusting server settings.
	//
	// Serve returns on ln.Close() or error, so usually it blocks forever.
	if _, docker := os.LookupEnv("RunInDocker"); docker {
		log.Println("\033[32mServer is running on: " + listener.Addr().String() + "\033[0m")
	} else {
		log.Println(success("Server is running on: ", listener.Addr()))
	}

	if err := http.Serve(listener, db.Handler); err != nil {
		if _, docker := os.LookupEnv("RunInDocker"); docker {
			log.Fatalln("\033[31mError in Serve: " + err.Error() + "\033[0m")
		}
		log.Fatalln(failed("Error in Serve: ", err))
	}
}
