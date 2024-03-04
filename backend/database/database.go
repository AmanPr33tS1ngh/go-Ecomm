package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSetup() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Failed to connect to mongodb")
		return nil
	}
	fmt.Println("Connected to mongodb!")
	return mongoClient
}

var Client *mongo.Client = DBSetup()

func UserData (client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("ECommerce").Collection(collectionName)
	return collection
}

func ProductData (client *mongo.Client, collectionName string) *mongo.Collection{
	var productCollection *mongo.Collection = client.Database("ECommerce").Collection(collectionName)
	return productCollection
}
