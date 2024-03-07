package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
	rpcPort = "5001"
	// mongoURL = "mongodb://127.0.0.1:27017"
	mongoURL = "mongodb://mongo:27017"
	// mongoURL = "mongodb://admin:password@localhost:27017/logs"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo

	log.Println("hello inside logger")
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Error pinging server:", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start web server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Println(err)
	}
	go app.rpcListen()

	go app.gRPCListen()

	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) rpcListen() {
	log.Println("starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Println(err)
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			log.Println(err)
		}

		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURL)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging server:", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	return client, nil
}
