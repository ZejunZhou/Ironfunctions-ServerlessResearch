package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

const MongoDBURL = "mongodb://pc99.cloudlab.umass.edu:27017"

func initializeUserDatabase(client *mongo.Client) bool {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("user-db").Collection("user")

	// Try to connect to the collection and print some data
	var result User
	err = collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
	return true
}

func main() {
	clientOptions := options.Client().ApplyURI(MongoDBURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	initializeUserDatabase(client)
}
