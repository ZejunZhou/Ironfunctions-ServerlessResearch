package rate

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	// "io/ioutil"
	// "log"
	// "os"
	"context"
	"sort"

	"github.com/go-redis/redis/v8"
	pb "github.com/harlow/go-micro-services/services/rate/proto"

	"strings"
)

// const name = "srv-rate"

// Server implements the rate service
type Server struct {
	MongoSession *mgo.Session
	// MemcClient *memcache.Client
	RedisClient *redis.Client
}

// Run starts the server
func (s *Server) Init() error {
	return nil
}

// GetRates gets rates for hotels for specific date range.
func (s *Server) GetRates(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	res := new(pb.Result)
	// session, err := mgo.Dial("mongodb-rate")
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()

	ratePlans := make(RatePlans, 0)

	session := s.MongoSession.Copy()
	defer session.Close()

	for _, hotelID := range req.HotelIds {
		// first check redis
		item, err := s.RedisClient.Get(ctx, hotelID).Result()
		if err == nil {
			// redis hit
			rate_strs := strings.Split(item, "\n")

			for _, rate_str := range rate_strs {
				if len(rate_str) != 0 {
					rate_p := new(pb.RatePlan)
					json.Unmarshal([]byte(rate_str), rate_p)
					ratePlans = append(ratePlans, rate_p)
				}
			}
		} else if err == redis.Nil {
			// redis miss, set up mongo connection
			c := session.DB("rate-db").C("inventory")
			// log.Println("REDIS miss: hotelID = ", hotelID)
			memc_str := ""

			tmpRatePlans := make(RatePlans, 0)
			err := c.Find(&bson.M{"hotelId": hotelID}).All(&tmpRatePlans)
			if err != nil {
				panic(err)
			} else {
				for _, r := range tmpRatePlans {
					ratePlans = append(ratePlans, r)
					rate_json, err := json.Marshal(r)
					if err != nil {
						fmt.Printf("json.Marshal err = %s\n", err)
					}
					memc_str = memc_str + string(rate_json) + "\n"
				}
			}

			// write to redis
			s.RedisClient.Set(ctx, hotelID, memc_str, 0)

		} else {
			fmt.Printf("Redis error = %s\n", err)
			panic(err)
		}
	}

	sort.Sort(ratePlans)
	res.RatePlans = ratePlans

	return res, nil
}

type RatePlans []*pb.RatePlan

func (r RatePlans) Len() int {
	return len(r)
}

func (r RatePlans) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RatePlans) Less(i, j int) bool {
	return r[i].RoomType.TotalRate > r[j].RoomType.TotalRate
}
