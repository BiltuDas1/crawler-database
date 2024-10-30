package db

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/zeebo/xxh3"
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
	for low <= high {
		mid := low + (high-low)/2
		if dbPositions[mid] < value {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	if low < len(dbPositions) {
		return dbPositions[low]
	}

	// Returning first database hash cause any database can't hold the data
	return dbPositions[0]
}

// InitializeDBList Function to load all the ip addresses and ports of all of the databases from the file
func InitializeDBList(filename string, runningInDocker bool) {
	dbList = make(map[uint64]NodeDetails)
	ipv4PortRegex := `^(([0-9][0-9]?|[0-1][0-9][0-9]|[2][0-4][0-9]|[2][5][0-5])\.){3}([0-9][0-9]?|[0-1][0-9][0-9]|[2][0-4][0-9]|[2][5][0-5]):([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
	ipv4Port := regexp.MustCompile(ipv4PortRegex)

	file, err := os.Open(filename)
	if err != nil {
		if runningInDocker {
			log.Fatalln("\033[31mNo db.txt file found, Please make sure it is attached with the container as a secret file.\033[0m")
		}
		log.Fatalln(failedLog("No db.txt file found, Please create one in the same path where the main.go file located. For More information: https://github.com/BiltuDas1/crawler-db/wiki"))
	}

	// Close the file when function ends
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Warn: Unable to close the file correctly")
		}
	}(file)

	scanner := bufio.NewScanner(file)
	// Adding each line elements to the db list
	for scanner.Scan() {
		line := scanner.Text()

		// When it's valid ip address + port not found, then skip the line
		if !ipv4Port.MatchString(line) {
			continue
		}

		hash := xxh3.HashString(line)
		dbPositions = append(dbPositions, hash)

		ipList := strings.Split(line, ":")
		port, err := strconv.Atoi(ipList[1])

		if err != nil {
			if runningInDocker {
				log.Fatalln("\033[31m" + err.Error() + "\033[0m")
			}
			log.Fatalln(failedLog(err))
		}

		dbList[hash] = NodeDetails{IP: ipList[0], Port: uint16(port)}
	}

	// If the scanner failed to scan the file
	if err := scanner.Err(); err != nil {
		if runningInDocker {
			log.Fatalln("\033[31m" + err.Error() + "\033[0m")
		}
		log.Fatalln(failedLog(err))
	}

	// Send a Warning if db.txt file is empty
	if len(dbPositions) == 0 {
		log.Println("Warn: db.txt file is empty")
	}

	// Sorting the Positions
	slices.Sort(dbPositions)
}

// GetDatabase Function to Return the database info by checking the data
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
