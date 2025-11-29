package logger

import (
	"bufio"
	"bytes"
	"go.uber.org/zap"
)

// logFunc method driven
type logFunc func(logger *zap.SugaredLogger, args ...interface{})

// logfFunc method driven
type logfFunc func(logger *zap.SugaredLogger, format string, args ...interface{})

// logwFunc method driven
type logwFunc func(logger *zap.SugaredLogger, msg string, keysAndValues ...interface{})

type LogService struct {
	zap         *zap.Logger
	base        *zap.SugaredLogger
	serviceCode int32
	serviceName string
	rabbitMQ    *RabbitMQ
	jsonWriter  *jsonWriter
	centralized bool

	debug  logFunc
	debugf logfFunc
	debugw logwFunc

	info  logFunc
	infof logfFunc
	infow logwFunc

	warn  logFunc
	warnf logfFunc
	warnw logwFunc

	error  logFunc
	errorf logfFunc
	errorw logwFunc

	dpanic  logFunc
	dpanicf logfFunc
	dpanicw logwFunc

	panic  logFunc
	panicf logfFunc
	panicw logwFunc

	fatal  logFunc
	fatalf logfFunc
	fatalw logwFunc
}

type jsonWriter struct {
	bufWriter  *bufio.Writer
	byteBuffer *bytes.Buffer
}

type logData struct {
	Data interface{} `bson:"data" json:"data"`
}

type message struct {
	ServiceCode int32  `bson:"service_code,omitempty" json:"service_code,omitempty"`
	ServiceName string `bson:"service_name,omitempty" json:"service_name,omitempty"`
	Log         string `bson:"log,omitempty" json:"log,omitempty"`
}
