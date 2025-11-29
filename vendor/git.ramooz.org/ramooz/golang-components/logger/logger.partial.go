package logger

import (
	"bufio"
	"fmt"
	colorable2 "github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

// initZapCore initialize zap cores
func initZapCore(serviceName string, opt *Options, jw *bufio.Writer) zapcore.Core {
	var zapCores []zapcore.Core

	if opt.ConsoleWriter {
		zapCores = append(zapCores, consoleCoreByLevel(opt.LogLevel, opt.Colorable, opt.Development)...)
	}
	if opt.Centralized != nil {
		zapCores = append(zapCores, jsonCoreByLevel(opt.Centralized.LogLevel, jw, opt.TimeFormat, opt.Development)...)
	}
	if !opt.Development {
		zapCores = append(zapCores, fileCoreByLevel(serviceName, opt.LogPath, opt.LogLevel, opt.Rotation, opt.Development)...)
	}

	return zapcore.NewTee(zapCores...)
}

// jsonCoreByLevel create cores base of levels
func jsonCoreByLevel(levels []Level, jw *bufio.Writer, timeFormat TimeFormat, development bool) []zapcore.Core {
	var zapCores []zapcore.Core
	for _, level := range levels {
		zapCores = append(zapCores, newCoreJson(zapcore.Level(level), jw, timeFormat, development))
	}
	return zapCores
}

// consoleCoreByLevel create cores base on input level to up
func consoleCoreByLevel(level Level, colorable bool, development bool) []zapcore.Core {
	var zapCores []zapcore.Core
	for l := range level.levelMap() {
		if l >= level {
			zapCores = append(zapCores, newCoreConsole(zapcore.Level(l), colorable, development))
		}
	}
	return zapCores
}

// fileCoreByLevel create cores base on input level to up
func fileCoreByLevel(serviceName, logPath string, level Level, rotation *RotationConfig, development bool) []zapcore.Core {
	var zapCores []zapcore.Core
	for l := range level.levelMap() {
		if l >= level {
			zapCores = append(zapCores, newCoreFile(serviceName, logPath, zapcore.Level(l), rotation, development))
		}
	}
	return zapCores
}

// newCoreJson add json writer for zap core
func newCoreJson(logLevel zapcore.Level, jw *bufio.Writer, timeFormat TimeFormat, development bool) zapcore.Core {
	Syncer := getJsonWriter(jw)
	return zapcore.NewCore(setEncoderJson(development, timeFormat), zapcore.Lock(Syncer), zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return logLevel == level
	}))
}

// newCoreConsole add console writer for zap core
func newCoreConsole(logLevel zapcore.Level, colorable, development bool) zapcore.Core {
	return zapcore.NewCore(setEncoderConsole(colorable, development), zapcore.AddSync(colorable2.NewColorableStdout()), zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return logLevel == level
	}))
}

// newCoreFile write logs into file
func newCoreFile(service, logPath string, logLevel zapcore.Level, rotation *RotationConfig, development bool) zapcore.Core {
	Syncer := getLogWriter(service, logPath, logLevel, rotation)
	return zapcore.NewCore(setEncoderFile(development), zapcore.Lock(Syncer), zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return logLevel == level
	}))
}

// setEncoderJson set encoder json structure for zap
func setEncoderJson(development bool, timeFormat TimeFormat) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{}
	if development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		encoderConfig.StacktraceKey = "stacktrace"
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.StacktraceKey = "stacktrace"
	}
	switch timeFormat {
	case ISO8601_TIME:
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	case RFC3339_TIME:
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	case RFC3339NANO_TIME:
		encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	default:
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	encoderConfig.FunctionKey = "function"
	return zapcore.NewJSONEncoder(encoderConfig)
}

// setEncoderConsole set encoder console writer for zap
func setEncoderConsole(colorable, development bool) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{}
	if development {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		if colorable {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		if colorable {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// setEncoderFile set encoder file writer for zap
func setEncoderFile(development bool) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{}
	if development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		encoderConfig.StacktraceKey = "stacktrace"
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	encoderConfig.FunctionKey = "function"
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriter create log file
func getLogWriter(service, logPath string, logType zapcore.Level, rotation *RotationConfig) zapcore.WriteSyncer {
	logTypeStr := logType.String()
	r := &lumberjack.Logger{}
	path := filepath.Join(logPath, fmt.Sprintf("%v/%v_%v.log", logTypeStr, service, logTypeStr))
	if rotation != nil {
		r.Filename = path
		r.MaxAge = rotation.MaxAge
		r.MaxSize = rotation.FileSize
		r.Compress = rotation.Compress
		return zapcore.AddSync(r)
	}
	return zapcore.AddSync(createBasicLogFile(path))
}

// getJsonWriter for marshal to json
func getJsonWriter(jw *bufio.Writer) zapcore.WriteSyncer {
	return zapcore.AddSync(jw)
}

func createBasicLogFile(logPath string) *os.File {
	file, err := os.Create(logPath)
	if err != nil {
		panic(err)
	}
	return file
}

// centralizedLog send log to centralized logging service
func centralizedLog(jw *jsonWriter, sCode int32, sName string, level Level, r *RabbitMQ) {
	for _, centLevel := range r.LogLevel {
		if level == centLevel {
			r.SendLogToCentralizedService(jw, sCode, sName)
		}
	}
}
