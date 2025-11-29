/*
Logger components for logging service logs into file, console and centralized logging service.

Example logger usage:
	package main

	import (
		"git.ramooz.org/ramooz/golang-components/logger"
		"math/rand"
		"time"
	)

	var logObj *logger.LogService

	func main() {
		logObj = logger.NewLogger(3, "test", &logger.Options{
			Rotation:      &logger.RotationConfig{MaxAge: 1, FileSize: 2, Compress: false},
			LogLevel:      logger.DebugLevel,
			ConsoleWriter: true,
			Centralized:   &logger.RabbitMQ{ServerAddress: "127.0.0.1", UserName: "user", Password: "password", LogLevel: []logger.Level{logger.ErrorLevel, logger.WarnLevel}},
			Colorable:     true,
			TimeFormat:    logger.ISO8601_TIME,
			Development:   false,
		})

		for {
			testWarn()
			time.Sleep(10 * time.Second)
			testError()
		}
	}

	func testWarn() {
		logObj.Warn("this is a warn")
		logObj.Warnf("warn message %d", rand.Int())
	}

	func testError() {
		logObj.Errorf("message %v", 1)
		logObj.Errorw("errorw message for test", "ip", "127.0.0.1")
	}
*/
package logger
