package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/assincrono/learn-pub-sub-starter/internal/gamelogic"
	"github.com/assincrono/learn-pub-sub-starter/internal/pubsub"
	"github.com/assincrono/learn-pub-sub-starter/internal/routing"
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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Something went wrong on getting username.")
		return
	}

	pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, username, routing.PauseKey, pubsub.Transient)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan

	fmt.Printf("\nReceived signal: %s\n", sig)
}
