package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Println("Something went wrong on connection.")
		return
	}

	defer connection.Close()

	fmt.Println("Connection went successfully!")

	channel, err := connection.Channel();

	if err != nil {
		fmt.Println("Something went wrong on creating a channel.");
		return;
	}

	pubsub.publishJSON()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan

	fmt.Printf("\nReceived signal: %s\n", sig)
}
