package main

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	rabbitConn, err := connect()

	if err != nil {
		log.Print(err)
	}

	defer rabbitConn.Close()
	fmt.Println("connected to rabbitMQ")

}

func connect() (*amqp.Connection, error) {
	var counts int64
	backOff := time.Second * 1
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			log.Println(err)
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
