package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)

	if err != nil {
		fmt.Println("Something went wrong in marshalling val.")
		return err
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: data})
}

type SimpleQueueType string

const (
    // Durable queues survive broker restarts
    Durable SimpleQueueType = "durable"

    // Transient queues disappear when the broker restarts
    Transient SimpleQueueType = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel();
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to create channel: %w", err)
	}

	durable := queueType == Durable
	autoDelete := queueType == Transient
	exclusive := queueType == Transient

	queue, err := channel.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil);
	if err != nil {
		fmt.Println("Something went wrong on creating a queue.");
		return nil, amqp.Queue{}, fmt.Errorf("failed to create queue: %w", err)
	}

	err = channel.QueueBind(queue.Name, key, exchange, false, nil);
	if err != nil {
        channel.Close()
        return nil, amqp.Queue{}, fmt.Errorf("failed to bind queue: %w", err)
    }
	
	return channel, queue, nil
}