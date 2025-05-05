package mq

import (
	"github.com/rabbitmq/amqp091-go"
)

type MQManager struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewMQManager(url, queueName string) (*MQManager, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, //auto-delete
		false, //exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &MQManager{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (m *MQManager) Publish(message string) error {
	return m.channel.Publish(
		"",           // exchange
		m.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (m *MQManager) Close() {
	m.channel.Close()
	m.conn.Close()
}
