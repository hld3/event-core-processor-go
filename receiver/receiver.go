package receiver

import "github.com/streadway/amqp"

type RabbitMQReceiver struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewReceiver() (*RabbitMQReceiver, error) {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQReceiver{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (r *RabbitMQReceiver) StartReceiving() (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *RabbitMQReceiver) Close() {
	r.channel.Close()
	r.conn.Close()
}
