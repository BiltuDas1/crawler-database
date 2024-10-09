package db

import (
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	xxh3 "github.com/zeebo/xxh3"
)

var dbList map[uint64]NodeDetails
var dbPositions []uint64

// Custom Binary Search
func findDB(value uint64) uint64 {
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

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(failedLog(err))
	}

	// Adding each line elements to the db list
	for _, value := range file {
		// When newline occurs
		if value == 10 {
			hash = xxh3.Hash(temp)
			dbPositions = append(dbPositions, hash)

			ipList := strings.Split(string(temp), ":")
			port, err := strconv.Atoi(ipList[1])

			if err != nil {
				log.Fatal(failedLog(err))
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
			log.Fatal(failedLog(err))
		}

		dbList[hash] = NodeDetails{IP: ipList[0], Port: uint16(port)}
	}

	// Sorting the Positions
	slices.Sort(dbPositions)
}

// Function to Return the database info by checking the data
func GetDatabase(hash uint64) (ip string, port uint16) {
	// log.Println(dbPositions)
	dbPos := findDB(hash)
	ipaddr := dbList[dbPos]
	ip = ipaddr.IP
	port = ipaddr.Port
	return
}
