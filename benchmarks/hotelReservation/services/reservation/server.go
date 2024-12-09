package reservation

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	pb "github.com/harlow/go-micro-services/services/reservation/proto"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"strconv"
	"time"
)

// const name = "srv-reservation"

// Server implements the user service
type Server struct {
	MongoSession *mgo.Session
	RedisClient  *redis.Client
}

// Run starts the server
func (s *Server) Init() error {
	return nil
}

// MakeReservation makes a reservation based on given information
func (s *Server) MakeReservation(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	res := new(pb.Result)
	res.HotelId = make([]string, 0)

	session := s.MongoSession.Copy()
	defer session.Close()

	c := session.DB("reservation-db").C("reservation")
	c1 := session.DB("reservation-db").C("number")

	inDate, _ := time.Parse(
		time.RFC3339,
		req.InDate+"T12:00:00+00:00")

	outDate, _ := time.Parse(
		time.RFC3339,
		req.OutDate+"T12:00:00+00:00")
	hotelId := req.HotelId[0]

	indate := inDate.String()[0:10]

	redis_date_num_map := make(map[string]int)

	for inDate.Before(outDate) {
		count := 0
		inDate = inDate.AddDate(0, 0, 1)
		outdate := inDate.String()[0:10]

		redis_key := hotelId + "_" + inDate.String()[0:10] + "_" + outdate
		val, err := s.RedisClient.Get(ctx, redis_key).Result()
		if err == nil {
			count, _ = strconv.Atoi(val)
			redis_date_num_map[redis_key] = count + int(req.RoomNumber)
		} else if err == redis.Nil {
			reserve := make([]reservation, 0)
			err := c.Find(&bson.M{"hotelId": hotelId, "inDate": indate, "outDate": outdate}).All(&reserve)
			if err != nil {
				panic(err)
			}

			for _, r := range reserve {
				count += r.Number
			}

			redis_date_num_map[redis_key] = count + int(req.RoomNumber)
		} else {
			fmt.Printf("Redis error = %s\n", err)
			panic(err)
		}

		redis_cap_key := hotelId + "_cap"
		val, err = s.RedisClient.Get(ctx, redis_cap_key).Result()
		hotel_cap := 0
		if err == nil {
			hotel_cap, _ = strconv.Atoi(val)
		} else if err == redis.Nil {
			var num number
			err = c1.Find(&bson.M{"hotelId": hotelId}).One(&num)
			if err != nil {
				panic(err)
			}
			hotel_cap = int(num.Number)
			s.RedisClient.Set(ctx, redis_cap_key, strconv.Itoa(hotel_cap), 0)
		} else {
			fmt.Printf("Redis error = %s\n", err)
			panic(err)
		}

		if count+int(req.RoomNumber) > hotel_cap {
			return res, nil
		}
		indate = outdate
	}

	for key, val := range redis_date_num_map {
		s.RedisClient.Set(ctx, key, strconv.Itoa(val), 0)
	}

	inDate, _ = time.Parse(
		time.RFC3339,
		req.InDate+"T12:00:00+00:00")

	indate = inDate.String()[0:10]

	for inDate.Before(outDate) {
		inDate = inDate.AddDate(0, 0, 1)
		outdate := inDate.String()[0:10]
		err := c.Insert(&reservation{
			HotelId:      hotelId,
			CustomerName: req.CustomerName,
			InDate:       indate,
			OutDate:      outdate,
			Number:       int(req.RoomNumber)})
		if err != nil {
			panic(err)
		}
		indate = outdate
	}

	res.HotelId = append(res.HotelId, hotelId)

	return res, nil
}

// CheckAvailability checks if given information is available
func (s *Server) CheckAvailability(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	res := new(pb.Result)
	res.HotelId = make([]string, 0)

	session := s.MongoSession.Copy()
	defer session.Close()

	c := session.DB("reservation-db").C("reservation")
	c1 := session.DB("reservation-db").C("number")

	for _, hotelId := range req.HotelId {
		inDate, _ := time.Parse(
			time.RFC3339,
			req.InDate+"T12:00:00+00:00")

		outDate, _ := time.Parse(
			time.RFC3339,
			req.OutDate+"T12:00:00+00:00")

		indate := inDate.String()[0:10]

		for inDate.Before(outDate) {
			count := 0
			inDate = inDate.AddDate(0, 0, 1)
			outdate := inDate.String()[0:10]

			redis_key := hotelId + "_" + inDate.String()[0:10] + "_" + outdate
			val, err := s.RedisClient.Get(ctx, redis_key).Result()

			if err == nil {
				count, _ = strconv.Atoi(val)
			} else if err == redis.Nil {
				// log.Println("REDIS miss: reservation = ", redis_key)

				reserve := make([]reservation, 0)
				err := c.Find(&bson.M{"hotelId": hotelId, "inDate": indate, "outDate": outdate}).All(&reserve)
				if err != nil {
					panic(err)
				}
				for _, r := range reserve {
					count += r.Number
				}
				s.RedisClient.Set(ctx, redis_key, strconv.Itoa(count), 0)
			} else {
				fmt.Printf("Redis error = %s\n", err)
				panic(err)
			}

			redis_cap_key := hotelId + "_cap"
			val, err = s.RedisClient.Get(ctx, redis_cap_key).Result()
			hotel_cap := 0

			if err == nil {
				hotel_cap, _ = strconv.Atoi(val)
			} else if err == redis.Nil {
				var num number
				// log.Println("REDIS miss: reservation hotelId = ", hotelId)
				err = c1.Find(&bson.M{"hotelId": hotelId}).One(&num)
				if err != nil {
					panic(err)
				}
				hotel_cap = int(num.Number)
				s.RedisClient.Set(ctx, redis_cap_key, strconv.Itoa(hotel_cap), 0)
			} else {
				fmt.Printf("Redis error = %s\n", err)
				panic(err)
			}

			if count+int(req.RoomNumber) > hotel_cap {
				break
			}
			indate = outdate

			if inDate.Equal(outDate) {
				res.HotelId = append(res.HotelId, hotelId)
			}
		}
	}

	return res, nil
}

type reservation struct {
	HotelId      string `bson:"hotelId"`
	CustomerName string `bson:"customerName"`
	InDate       string `bson:"inDate"`
	OutDate      string `bson:"outDate"`
	Number       int    `bson:"number"`
}

type number struct {
	HotelId string `bson:"hotelId"`
	Number  int    `bson:"numberOfRoom"`
}
