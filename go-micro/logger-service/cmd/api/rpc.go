package main

import (
	"context"
	"log"
	"logger/data"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload *RPCPayload, resp *string) error {
	log.Println("inside LogInfo")
	collection := client.Database("logs").Collection("logs")

	res, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error writing into mongo", err)
	}

	log.Println(res.InsertedID)

	*resp = "processed payload via RPC" + payload.Name

	return nil
}
