package cli

import (
	"errors"
	"github.com/Tritium-leo/gear-go/cfg"
	"github.com/Tritium-leo/gear-go/mq/model"
	"github.com/streadway/amqp"
	"sync"
)

type RabbitMQClient struct {
	rw        sync.RWMutex
	Client    *amqp.Connection
	Producers map[string]*model.Producer
	Consumers map[string]*model.Consumer
}
type RabbitMQConfig struct {
	Username string
	password string
	Address  string
	Port     int
}

func NewRabbitMQClient(config *cfg.Config, parentPath string) (cli *RabbitMQClient, err error) {
	var conn *amqp.Connection
	conn, err = amqp.Dial("amqp://")
	if err != nil {
		return nil, err
	}
	cli = &RabbitMQClient{
		Client:    conn,
		Producers: map[string]*model.Producer{},
		Consumers: map[string]*model.Consumer{},
	}
	return cli, err
}

func (c *RabbitMQClient) Close() {
	for _, p := range c.Producers {
		p.Close()
	}
	for _, cons := range c.Consumers {
		cons.Close()
	}
	c.Client.Close()
}

func (c *RabbitMQClient) NewRabbitMQProducer(resourceName string) (ch *amqp.Channel, q amqp.Queue, err error) {
	ch, err = c.Client.Channel()
	q, err = ch.QueueDeclare(resourceName, false, false, false, false, nil)

	if err != nil {
		return nil, nil, errors.New("Create RabbitMq Producer Failed :" + resourceName)
	}

	return

}

type ProcessFunc func(msg []byte, headers map[string]interface{}) (err error)

func (c *RabbitMQClient) NewRabbitMQConsumer(cli *amqp.Connection, resourceName string, processfunc ProcessFunc) {

}
