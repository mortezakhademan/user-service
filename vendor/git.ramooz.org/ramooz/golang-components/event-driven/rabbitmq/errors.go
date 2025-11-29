package rabbitmq

import componentsError "git.ramooz.org/ramooz/golang-components/error-handler"

var ErrorMessages = map[string]map[int]string{
	componentsError.LANG_FA: {
		ERROR_SERVICE_NAME:            "service name is empty",
		ERROR_URI_ADDRESS:             "uri address is invalid, please enter amqp://guest:guest@localhost:5672 for example",
		ERROR_ROUTING_KEYS_EMPTY:      "routing keys is empty",
		ERROR_CONNECTION_CLOSED:       "rabbitMQ connection closed, try to reconnect",
		ERROR_NIL_CONNECTION:          "nil rabbitmq connection",
		ERROR_EXCHANGE_ALREADY_EXISTS: "exchange already declared",
		ERROR_QUEUE_ALREADY_EXISTS:    "queue already declared",
		ERROR_EXHCNAGE_NOT_FOUND:      "exchange not declare",
	},
	componentsError.LANG_EN: {
		ERROR_SERVICE_NAME:            "service name is empty",
		ERROR_URI_ADDRESS:             "uri address is invalid, please enter amqp://guest:guest@localhost:5672 for example",
		ERROR_ROUTING_KEYS_EMPTY:      "routing keys is empty",
		ERROR_CONNECTION_CLOSED:       "rabbitMQ connection closed, try to reconnect",
		ERROR_NIL_CONNECTION:          "nil rabbitmq connection",
		ERROR_EXCHANGE_ALREADY_EXISTS: "exchange already declared",
		ERROR_QUEUE_ALREADY_EXISTS:    "queue already declared",
		ERROR_EXHCNAGE_NOT_FOUND:      "exchange not declare",
	},
}

const (
	ERROR_SERVICE_NAME            = 401110
	ERROR_URI_ADDRESS             = 401117
	ERROR_ROUTING_KEYS_EMPTY      = 401118
	ERROR_CONNECTION_CLOSED       = 401119
	ERROR_NIL_CONNECTION          = 400111
	ERROR_EXCHANGE_ALREADY_EXISTS = 400112
	ERROR_QUEUE_ALREADY_EXISTS    = 400113
	ERROR_EXHCNAGE_NOT_FOUND      = 400114
)
