package logger

import "errors"

var (
	SERVICE_CODE_ERROR             = errors.New("set the service code to more than 0")
	SERVICE_NAME_ERROR             = errors.New("service name is empty")
	RABBIT_SERVER_EMPTY_ERROR      = errors.New("server address rabbitMQ is empty")
	RABBIT_VHOST_EMPTY_ERROR       = errors.New("vHost rabbitMQ is empty")
	RABBIT_USER_EMPTY_ERROR        = errors.New("username rabbitMQ is empty")
	RABBIT_PASS_EMPTY_ERROR        = errors.New("password rabbitMQ is empty")
	RABBIT_LOG_LEVEL_ERROR         = errors.New("log levels rabbitMQ is empty")
	ERROR_FILEPATH_IS_NOT_ABS      = errors.New("log path is not relative path")
	ERROR_FILEPATH_IS_NOT_WRITABLE = errors.New("log path don't have permission write")
	ERROR_SENTRY_DSN_IS_EMPTY      = errors.New("sentry dsn is empty")
)
