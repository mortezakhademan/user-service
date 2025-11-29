package logger

import (
	"go.uber.org/zap/zapcore"
)

type Level zapcore.Level

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	FatalLevel
)

// levelMap create map from levels
func (l Level) levelMap() map[Level]struct{} {
	return map[Level]struct{}{
		DebugLevel:  {},
		InfoLevel:   {},
		WarnLevel:   {},
		ErrorLevel:  {},
		DPanicLevel: {},
		PanicLevel:  {},
		FatalLevel:  {},
	}
}
