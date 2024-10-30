package main

import (
	"log"
	"net"
	"os"

	db "crawler-db/src"

	"github.com/fatih/color"
	http "github.com/valyala/fasthttp"
)

func main() {
	runningInDocker := os.Getenv("RunInDocker") != ""
	failed := color.New(color.FgRed).SprintFunc()
	warning := color.New(color.FgYellow).SprintFunc()
	success := color.New(color.FgGreen).SprintFunc()

	// Getting the local IP Address
	ipAddress, err := db.GetLocalIP(runningInDocker)
	if err != nil {
		log.Println(warning("Warning: " + err.Error()))
	}

	// Initialize DB
	if runningInDocker {
		db.InitializeDBList("/run/secrets/db.txt", true)
	} else {
		db.InitializeDBList("db.txt", false)
	}

	listener, err := net.Listen("tcp4", ipAddress+":"+os.Getenv("PORT"))
	if err != nil {
		if runningInDocker {
			log.Fatalln("\033[31mError in net.Listen: " + err.Error() + "\033[0m")
		}
		log.Fatalln(failed("Error in net.Listen: ", err))
	}

	// Start the server with default settings.
	// Create Server instance for adjusting server settings.
	//
	// Serve returns on ln.Close() or error, so usually it blocks forever.
	if runningInDocker {
		log.Println("\033[32mServer is running on: " + listener.Addr().String() + "\033[0m")
	} else {
		log.Println(success("Server is running on: ", listener.Addr()))
	}

	if err := http.Serve(listener, db.Handler); err != nil {
		if runningInDocker {
			log.Fatalln("\033[31mError in Serve: " + err.Error() + "\033[0m")
		}
		log.Fatalln(failed("Error in Serve: ", err))
	}
}
