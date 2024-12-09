package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	pb "github.com/harlow/go-micro-services/services/user/proto"
	faas "github.com/harlow/go-micro-services/worker"
	"github.com/harlow/go-micro-services/worker/types"

	"github.com/harlow/go-micro-services/services/user"
	"github.com/harlow/go-micro-services/utils"

	"gopkg.in/mgo.v2"
	// "github.com/bradfitz/gomemcache/memcache"
)

type funcHandlerFactory struct {
	mongoSession *mgo.Session
}

func (f *funcHandlerFactory) New(env types.Environment, funcName string) (types.FuncHandler, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (f *funcHandlerFactory) GrpcNew(env types.Environment, service string) (types.GrpcFuncHandler, error) {
	if service != "user.User" {
		return nil, fmt.Errorf("Unknown service: %s", service)
	}
	srv := &user.Server{MongoSession: f.mongoSession}
	err := srv.Init()
	if err != nil {
		return nil, err
	}
	return utils.NewGrpcFuncHandler(srv, pb.UserMethods)
}

func main() {
	// initializeDatabase()
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	mongo_session, err := initializeDatabase(result["UserMongoAddress"])
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer mongo_session.Close()

	faas.Serve(&funcHandlerFactory{mongoSession: mongo_session})
}
