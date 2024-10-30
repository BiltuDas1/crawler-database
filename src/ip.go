package db

import (
	"errors"
	"log"
	"net"
	"strings"
)

// GetLocalIP Get the public local IP of the container
// If the code is not running on into any docker container then it will return err
func GetLocalIP(runningInDocker bool) (ipAddr string, errorOf error) {
	if !runningInDocker {
		errorOf = errors.New("not running in a container")
		return
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln("\033[31mError in net.Listen: " + err.Error() + "\033[0m")
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		// Handle err
		if err != nil {
			log.Fatalln("\033[31mError in net.Listen: " + err.Error() + "\033[0m")
		}

		for _, addr := range addrs {
			// Choosing the IPV4 Address which is not loopback address
			if !strings.Contains(addr.String(), "127.0.0.1") || !strings.Contains(addr.String(), "::1") {
				ip := strings.Split(addr.String(), "/")
				ipAddr = ip[0]
			}
		}
	}

	return
}
