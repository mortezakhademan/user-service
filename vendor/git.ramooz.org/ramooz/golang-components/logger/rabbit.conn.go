package logger

import (
	"git.ramooz.org/ramooz/golang-components/event-driven/rabbitmq"
)

var event *rabbitmq.Connection

func initializeRabbitMQ(serviceName string, rb *RabbitMQ) error {
	var err error
	event, err = rabbitmq.NewConnection(serviceName, &rabbitmq.Options{
		UriAddress:      rabbitmq.CreateURIAddress(rb.UserName, rb.Password, rb.ServerAddress, rb.VHost),
		DurableExchange: true,
	}, nil)
	if err != nil {
		return err
	}
	if err := event.ExchangeDeclare("centralizedLogging", rabbitmq.TOPIC); err != nil {
		return err
	}
	if err := event.DeclarePublisherQueue("logger", "centralizedLogging", "log"); err != nil {
		return err
	}
	return nil
}
