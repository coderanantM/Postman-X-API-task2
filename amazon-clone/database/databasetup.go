package database

import (
	"context"
	"log"
	"fmt"
	"time"
	

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
)

func DBSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err!=nil {
		log.Fatal(err)
	}

	ctx, cancel :- context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err!=nil {
		log.Fatal(err)

	}

	err = client.Ping(context.TODO(), nil)
	if err!=nil {
		log.Println("failed to connect to MongoDB :(")
		return nil
	}

	fmt.Println("Successfully connected to MongoDB")
	return client
}

	var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection  = client.Database("Amazon").Collection(collectionName)
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection{
	var productCollection *mongo.Collection = client.Database("Amazon").Collection(CollectionName(collectionName))
	return productCollection

}