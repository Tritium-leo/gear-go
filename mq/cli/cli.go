package cli

type ClientInterface interface {
	NewRabbitMQProducer()
	NewRabbitMQConsumer()
}
