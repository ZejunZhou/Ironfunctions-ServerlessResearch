package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hotel_recommendations/types"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const MongoDBURL = "mongodb://pc21.cloudlab.umass.edu"
const maxRecommendationResults = 5

type HotelRecommendation struct {
	Require string  `json:"require"` // dis, rate, price
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	// HLat    float64
	// HLon    float64
}

type Hotel struct {
	ID     bson.ObjectId `bson:"_id"`
	HId    string        `bson:"hotelId"`
	HLat   float64       `bson:"lat"`
	HLon   float64       `bson:"lon"`
	HRate  float64       `bson:"rate"`
	HPrice float64       `bson:"price"`
}

// loadRecommendations loads hotel recommendations from mongodb.
func loadRecommendations() map[string]Hotel {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MongoDBURL))
	if err != nil {
		log.Println("Failed to create mongodb client: ", err)
		return nil
	}

	collection := client.Database("recommendation-db").Collection("recommendation")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Failed to get hotels data: ", err)
		return nil
	}
	defer cursor.Close(ctx)

	var hotels []Hotel
	if err = cursor.All(ctx, &hotels); err != nil {
		log.Println("Failed to decode hotels data: ", err)
		return nil
	}

	profiles := make(map[string]Hotel)
	for _, hotel := range hotels {
		profiles[hotel.HId] = hotel
	}

	return profiles
}

// should connect to mongodb
func main() {
	for {
		res := http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			StatusCode: 200,
			Status:     "OK",
		}

		r := bufio.NewReader(os.Stdin)
		req, err := http.ReadRequest(r)
		var buf bytes.Buffer
		if err != nil {
			res.StatusCode = 500
			res.Status = http.StatusText(res.StatusCode)
			fmt.Fprintln(&buf, err)
		} else {
			l, _ := strconv.Atoi(req.Header.Get("Content-Length"))
			p := make([]byte, l)
			r.Read(p)
			req := &HotelRecommendation{}
			if err := json.Unmarshal(p, req); err != nil {
				res.StatusCode = 400
				res.Status = http.StatusText(res.StatusCode)
				fmt.Fprintln(&buf, err)
			}

			HotelIds := []string{}
			if req.Require == "" {
				fmt.Println(&buf, "Require field is required")
			} else if req.Require == "dis" {
				hotels := loadRecommendations()
				p1 := &types.GeoPoint{
					Pid:  "",
					Plat: req.Lat,
					Plon: req.Lon,
				}
				min := math.MaxFloat64
				for _, hotel := range hotels {
					tmp := float64(types.Distance(p1, &types.GeoPoint{
						Pid:  "",
						Plat: hotel.HLat,
						Plon: hotel.HLon,
					})) / 1000
					if tmp < min {
						min = tmp
					}
				}
				for _, hotel := range hotels {
					tmp := float64(types.Distance(p1, &types.GeoPoint{
						Pid:  "",
						Plat: hotel.HLat,
						Plon: hotel.HLon,
					})) / 1000
					if tmp == min && len(HotelIds) < maxRecommendationResults {
						HotelIds = append(HotelIds, hotel.HId)
					}
				}
			} else if req.Require == "rate" {
				hotels := loadRecommendations()
				max := 0.0
				for _, hotel := range hotels {
					if hotel.HRate > max {
						max = hotel.HRate
					}
				}
				for _, hotel := range hotels {
					if hotel.HRate == max && len(HotelIds) < maxRecommendationResults {
						HotelIds = append(HotelIds, hotel.HId)
					}
				}
			} else if req.Require == "price" {
				min := math.MaxFloat64
				hotels := loadRecommendations()
				for _, hotel := range hotels {
					if hotel.HPrice < min {
						min = hotel.HPrice
					}
				}
				for _, hotel := range hotels {
					if hotel.HPrice == min && len(HotelIds) < maxRecommendationResults {
						HotelIds = append(HotelIds, hotel.HId)
					}
				}
			}
			fmt.Fprintf(&buf, "Hotel Recommendation: %v\n", HotelIds)
			// for k, vs := range req.Header {
			// 	fmt.Fprintf(&buf, "ENV: %s %#v\n", k, vs)
			// }
		}
		res.Body = ioutil.NopCloser(&buf)
		res.ContentLength = int64(buf.Len())
		res.Write(os.Stdout)
	}
}
