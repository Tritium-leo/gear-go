package model

import (
	"errors"
	"github.com/streadway/amqp"
)

type Producer struct {
	Nane string
	ch   amqp.Channel
	q    amqp.Queue
}

func (p *Producer) Close() error {
	return p.ch.Close()
}

func (c *Producer) PushMessage(body []byte) error {
	err := c.ch.Publish("", c.q.Name, false, false, amqp.Publishing{ContentType: "text/plane", Body: body})
	if err != nil {
		return errors.New("Push Message Failed ,Queue:" + c.q.Name + " Msg:" + string(body))
	}
	return nil
}
