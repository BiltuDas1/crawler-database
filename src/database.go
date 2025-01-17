package db

import (
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	xxh3 "github.com/zeebo/xxh3"
)

var dbList map[uint64]NodeDetails
var dbPositions []uint64

// Custom Binary Search
func findDB(value uint64) uint64 {
	// If dbPosition empty then return 0
	if len(dbPositions) == 0 {
		return 0
	}

	low, high := 0, len(dbPositions)-1
	var lastdb uint64 = 0

	for low <= high {
		mid := low + (high-low)/2

		if dbPositions[mid] == value {
			return value
		} else if dbPositions[mid] < value {
			low = mid + 1
		} else {
			high = mid - 1
			lastdb = dbPositions[mid]
		}
	}

	// Getting elements which is larger than the target
	if lastdb == 0 {
		// Returning first database hash cause any database can't hold the data
		return dbPositions[0]
	} else {
		return lastdb
	}
}

// Function to load all the ip addresses and ports of all of the databases from the file
func InitializeDBList(filename string) {
	dbList = make(map[uint64]NodeDetails)
	var temp []byte
	var hash uint64
	ipv4PortRegex := `^(([0-9][0-9]?|[0-1][0-9][0-9]|[2][0-4][0-9]|[2][5][0-5])\.){3}([0-9][0-9]?|[0-1][0-9][0-9]|[2][0-4][0-9]|[2][5][0-5]):([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
	ipv4Port := regexp.MustCompile(ipv4PortRegex)

	file, err := os.ReadFile(filename)
	if err != nil {
		if _, docker := os.LookupEnv("RunInDocker"); docker {
			log.Fatalln("\033[31mNo db.txt file found, Please create one and mount it into the /data directory\033[0m")
		}
		log.Fatalln(failedLog("No db.txt file found, Please create one in the same path where the main.go file located. For More information: https://github.com/BiltuDas1/crawler-db/wiki"))
	}

	// Adding each line elements to the db list
	for _, value := range file {
		// When newline occurs
		if value == 10 {
			// If the ip address and port is invalid then clear the temp buffer and continue the next iteration
			if !ipv4Port.MatchString(string(temp)) {
				temp = nil
				continue
			}

			hash = xxh3.Hash(temp)
			dbPositions = append(dbPositions, hash)

			ipList := strings.Split(string(temp), ":")
			port, err := strconv.Atoi(ipList[1])

			if err != nil {
				if _, docker := os.LookupEnv("RunInDocker"); docker {
					log.Fatalln("\033[31m" + err.Error() + "\033[0m")
				}
				log.Fatalln(failedLog(err))
			}

			dbList[hash] = NodeDetails{IP: ipList[0], Port: uint16(port)}
			temp = nil
		} else {
			temp = append(temp, value)
		}
	}

	// If the EOF doesn't have \n character then
	// add the temp variable values to the dbList
	if temp != nil {
		hash = xxh3.Hash(temp)
		dbPositions = append(dbPositions, hash)

		ipList := strings.Split(string(temp), ":")
		port, err := strconv.ParseUint(ipList[1], 10, 16)

		if err != nil {
			if _, docker := os.LookupEnv("RunInDocker"); docker {
				log.Fatalln("\033[31m" + err.Error() + "\033[0m")
			}
			log.Fatalln(failedLog(err))
		}

		dbList[hash] = NodeDetails{IP: ipList[0], Port: uint16(port)}
	}

	// Send a Warning if db.txt file is empty
	if len(dbPositions) == 0 {
		log.Println("Warn: db.txt file is empty")
	}

	// Sorting the Positions
	slices.Sort(dbPositions)
}

// Function to Return the database info by checking the data
func GetDatabase(hash uint64) (ip string, port uint16) {
	dbPos := findDB(hash)
	if dbPos == 0 {
		return
	}

	ipaddr := dbList[dbPos]
	ip = ipaddr.IP
	port = ipaddr.Port
	return
}
