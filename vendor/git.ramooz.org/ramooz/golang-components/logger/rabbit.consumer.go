package logger

import (
	"fmt"
	"git.ramooz.org/ramooz/golang-components/event-driven/rabbitmq"
)

// SendLogToCentralizedService send message to centralized logging service
func (r *RabbitMQ) SendLogToCentralizedService(log *jsonWriter, sCode int32, sName string) {
	if log == nil {
		return
	}
	err := event.Publish("centralizedLogging", "log", logData{Data: createMessage(log, sCode, sName)}, rabbitmq.PublishingOptions{})
	if err != nil {
		fmt.Println("logger can't be send log to centralized service, got error ", err)
	}
}

func createMessage(log *jsonWriter, sCode int32, sName string) message {
	log.bufWriter.Flush()
	newMsg := message{ServiceCode: sCode, ServiceName: sName, Log: log.byteBuffer.String()}
	// clear buffer data of buffer writer
	log.byteBuffer.Truncate(log.bufWriter.Buffered())
	return newMsg
}
