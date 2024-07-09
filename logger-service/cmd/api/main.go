package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/iBoBoTi/microservice-arch/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start web server
	log.Println("Starting server on port: ", webPort)
	app.serve()

}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Fatal(srv.ListenAndServe())

}

func connectToMongo() (*mongo.Client, error) {
	//create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(
		options.Credential{
			Username: "admin",
			Password: "password",
		},
	)

	//connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("Error connecting: %v", err)
		return nil, err
	}

	return c, nil
}
