package model

import "github.com/streadway/amqp"

type Consumer struct {
	ch amqp.Channel
	q  amqp.Queue
}

func (c *Consumer) Close() error {
	return c.ch.Close()
}
