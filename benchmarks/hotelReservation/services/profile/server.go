package profile

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	// "io/ioutil"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	pb "github.com/harlow/go-micro-services/services/profile/proto"
	// "strings"
)

// const name = "srv-profile"

// Server implements the profile service
type Server struct {
	MongoSession *mgo.Session
	// MemcClient *memcache.Client
	RedisClient *redis.Client
}

// Run starts the server
func (s *Server) Init() error {
	return nil
}

// GetProfiles returns hotel profiles for requested IDs
func (s *Server) GetProfiles(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	res := new(pb.Result)
	hotels := make([]*pb.Hotel, 0)

	session := s.MongoSession.Copy()
	defer session.Close()

	for _, i := range req.HotelIds {
		// first check redis
		val, err := s.RedisClient.Get(ctx, i).Result()
		if err == nil {
			// redis hit
			hotel_prof := new(pb.Hotel)
			json.Unmarshal([]byte(val), hotel_prof)
			hotels = append(hotels, hotel_prof)

		} else if err == redis.Nil {
			// redis miss, set up mongo connection
			c := session.DB("profile-db").C("hotels")
			// log.Println("REDIS miss: profile = ", i)
			hotel_prof := new(pb.Hotel)
			err := c.Find(bson.M{"id": i}).One(&hotel_prof)

			if err != nil {
				log.Println("Failed get hotels data: ", err)
			}

			hotels = append(hotels, hotel_prof)

			prof_json, err := json.Marshal(hotel_prof)
			if err == nil {
				// write to redis
				s.RedisClient.Set(ctx, i, prof_json, 0)
			}

		} else {
			fmt.Printf("Redis error = %s\n", err)
			panic(err)
		}
	}

	res.Hotels = hotels
	return res, nil
}
