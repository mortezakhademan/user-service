package logger

import (
	"fmt"
	"os"
)

type TimeFormat uint

const (
	ISO8601_TIME = iota // https://en.wikipedia.org/wiki/ISO_8601
	RFC3339_TIME        // https://en.wikipedia.org/wiki/ISO_8601#RFCs
	RFC3339NANO_TIME
)

var (
	defaultLogPath = fmt.Sprintf("%v/logs", getCWD()) // default local path for store log
)

// Options logger options for first initialization
type Options struct {
	Rotation      *RotationConfig // get rotation config
	LogLevel      Level           // get log level for store log level x to upper
	LogPath       string          // set log path store
	Centralized   *RabbitMQ       // send logs to centralized logging service using message broker rabbitMQ
	ConsoleWriter bool            // write logs in console
	TimeFormat    TimeFormat      // set Time format as ISO8601, RFC3339, RFC3339NANO for json encoder
	Colorable     bool            // set color mode levels for console writer
	Sentry        bool            // enable sentry for get log in sentry panel
	SentryDSN     string          // dsn for connect to sentry hub
	Development   bool            // logging in development mode
}

// RotationConfig create rotate flow for log files per day or log size create new log file
type RotationConfig struct {
	MaxAge   int  // per x day for create new log file
	FileSize int  // per x size create new log file
	Compress bool // compress log file with gzip
}

type RabbitMQ struct {
	ServerAddress string `json:"server_address"`
	VHost         string `json:"v_host"`
	UserName      string `json:"user_name"`
	Password      string `json:"password"`
	// LogLevel set some levels for send to centralized logging service,
	// for example []logger.Level{logger.ErrorLevel, logger.WarnLevel}
	LogLevel []Level
}

func GetDefaultOptions() *Options {
	return &Options{
		Rotation:      &RotationConfig{MaxAge: 10, FileSize: 200, Compress: false},
		LogLevel:      InfoLevel,
		LogPath:       defaultLogPath,
		ConsoleWriter: true,
		Centralized:   &RabbitMQ{},
		Colorable:     true,
		TimeFormat:    ISO8601_TIME,
		Sentry:        false,
		SentryDSN:     "",
		Development:   false,
	}
}

// getCWD get current work directory
func getCWD() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}
