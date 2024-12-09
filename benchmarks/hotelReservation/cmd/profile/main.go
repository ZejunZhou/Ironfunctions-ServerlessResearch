package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
	"gopkg.in/mgo.v2"

	"github.com/harlow/go-micro-services/services/profile"
	pb "github.com/harlow/go-micro-services/services/profile/proto"
	"github.com/harlow/go-micro-services/utils"
	faas "github.com/harlow/go-micro-services/worker"
	"github.com/harlow/go-micro-services/worker/types"
)

type funcHandlerFactory struct {
	mongoSession *mgo.Session
	// memcClient   *memcache.Client
	redisClient *redis.Client
}

func (f *funcHandlerFactory) New(env types.Environment, funcName string) (types.FuncHandler, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f *funcHandlerFactory) GrpcNew(env types.Environment, service string) (types.GrpcFuncHandler, error) {
	if service != "profile.Profile" {
		return nil, fmt.Errorf("unknown service: %s", service)
	}
	srv := &profile.Server{MongoSession: f.mongoSession, RedisClient: f.redisClient}
	err := srv.Init()
	if err != nil {
		return nil, err
	}
	return utils.NewGrpcFuncHandler(srv, pb.ProfileMethods)
}

func main() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	mongo_session := initializeDatabase(result["ProfileMongoAddress"])
	defer mongo_session.Close()

	// fmt.Printf("profile memc addr port = %s\n", result["ProfileMemcAddress"])
	// memc_client := memcache.New(result["ProfileMemcAddress"])
	// memc_client.Timeout = 100 * time.Millisecond
	// memc_client.MaxIdleConns = 64

	fmt.Printf("profile redis addr port = %s\n", result["ProfileRedisAddress"])
	redis_client := redis.NewClient(&redis.Options{
		Addr:     result["ProfileRedisAddress"],
		Password: "123", // no password set
		DB:       0,     // use default DB
	})
	pong, err := redis_client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to redis", err)
	} else {
		fmt.Println(pong)
	}

	faas.Serve(&funcHandlerFactory{mongoSession: mongo_session, redisClient: redis_client})
}
