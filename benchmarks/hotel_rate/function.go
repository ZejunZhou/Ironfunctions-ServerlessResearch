package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Request struct {
	HotelIds []string
	InDate   string
	OutDate  string
}

type RatePlan struct {
	// Define the fields for RatePlan
}

type RatePlans []*RatePlan

func (r RatePlans) Len() int           { return len(r) }
func (r RatePlans) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RatePlans) Less(i, j int) bool { return true } // Define your sorting logic

type Server struct {
	MemcClient *memcache.Client
}

const fifoPath = "/dev/shm/hotel_rate:0.0.1"
const memcachedAddr = "localhost:11211"
const mongoDBURL = "mongodb://localhost:27017"

func readFromFIFO(ch chan<- string) {
	fifoPathInput := fifoPath + "_input_" + os.Getenv("FUNC_ID")
	if _, err := os.Stat(fifoPathInput); os.IsNotExist(err) {
		if _, err := os.Create(fifoPathInput); err != nil {
			log.Fatalf("Failed to create FIFO: %v", err)
		}
	}

	for {
		file, err := os.Open(fifoPathInput)
		if err != nil {
			log.Fatalf("Failed to open FIFO: %v", err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			ch <- line
		}

		file.Close()

		// Empty the file after reading
		err = os.Truncate(fifoPathInput, 0)
		if err != nil {
			log.Fatalf("Failed to truncate FIFO: %v", err)
		}

		time.Sleep(1 * time.Second)
	}
}

func writeToFIFO(res string) {
	fifoPathOutput := fifoPath + "_output_" + os.Getenv("FUNC_ID")
	if _, err := os.Stat(fifoPathOutput); os.IsNotExist(err) {
		if _, err := os.Create(fifoPathOutput); err != nil {
			log.Fatalf("Failed to create FIFO: %v", err)
		}
	}

	file, err := os.OpenFile(fifoPathOutput, os.O_WRONLY|os.O_APPEND, os.ModeNamedPipe)
	if err != nil {
		log.Fatalf("Failed to open FIFO: %v", err)
	}

	_, err = file.WriteString(res + "\n")
	if err != nil {
		log.Fatalf("Failed to write to FIFO: %v", err)
	}

	file.Close()
}

func main() {
	notifyCh := make(chan string)
	go readFromFIFO(notifyCh)

	for {
		select {
		case req := <-notifyCh:
			// Parse the request
			var request Request
			err := json.Unmarshal([]byte(req), &request)
			if err != nil {
				log.Fatalf("Failed to unmarshal request: %v", err)
			}

			// convert the request to string
			reqStr, err := json.Marshal(request)
			if err != nil {
				log.Fatalf("Failed to marshal request: %v", err)
			}

			// Write to FIFO
			writeToFIFO(string(reqStr))
		}
	}

}

func processRequest(session *mgo.Session, req Request) (interface{}, error) {
	var ratePlans RatePlans
	memcClient := memcache.New(memcachedAddr)

	for _, hotelID := range req.HotelIds {
		// Check Memcached
		item, err := memcClient.Get(hotelID)
		if err == nil {
			// Memcached hit
			rateStrs := strings.Split(string(item.Value), "\n")
			for _, rateStr := range rateStrs {
				if len(rateStr) != 0 {
					rateP := new(RatePlan)
					json.Unmarshal([]byte(rateStr), rateP)
					ratePlans = append(ratePlans, rateP)
				}
			}
		} else if err == memcache.ErrCacheMiss {
			// Memcached miss, set up MongoDB connection
			c := session.DB("rate-db").C("inventory")

			memcStr := ""
			tmpRatePlans := make(RatePlans, 0)
			err := c.Find(bson.M{"hotelId": hotelID}).All(&tmpRatePlans)
			if err != nil {
				return nil, err
			}

			for _, r := range tmpRatePlans {
				ratePlans = append(ratePlans, r)
				rateJSON, err := json.Marshal(r)
				if err != nil {
					fmt.Printf("json.Marshal err = %s\n", err)
				}
				memcStr = memcStr + string(rateJSON) + "\n"
			}

			// Write to Memcached
			memcClient.Set(&memcache.Item{Key: hotelID, Value: []byte(memcStr)})
		} else {
			return nil, err
		}
	}

	sort.Sort(ratePlans)
	res := struct {
		RatePlans RatePlans
	}{RatePlans: ratePlans}

	return res, nil
}
