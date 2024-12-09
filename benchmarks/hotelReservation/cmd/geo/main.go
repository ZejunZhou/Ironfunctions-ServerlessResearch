package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/mgo.v2"

	"github.com/harlow/go-micro-services/services/geo"
	pb "github.com/harlow/go-micro-services/services/geo/proto"
	"github.com/harlow/go-micro-services/utils"
	faas "github.com/harlow/go-micro-services/worker"
	"github.com/harlow/go-micro-services/worker/types"
)

type funcHandlerFactory struct {
	mongoSession *mgo.Session
}

func (f *funcHandlerFactory) New(env types.Environment, funcName string) (types.FuncHandler, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (f *funcHandlerFactory) GrpcNew(env types.Environment, service string) (types.GrpcFuncHandler, error) {
	if service != "geo.Geo" {
		return nil, fmt.Errorf("Unknown service: %s", service)
	}
	srv := &geo.Server{MongoSession: f.mongoSession}
	err := srv.Init()
	if err != nil {
		return nil, err
	}
	return utils.NewGrpcFuncHandler(srv, pb.GeoMethods)
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

	mongo_session := initializeDatabase(result["GeoMongoAddress"])
	defer mongo_session.Close()

	faas.Serve(&funcHandlerFactory{mongoSession: mongo_session})
}
