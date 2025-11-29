package logger

// Debug uses fmt.Sprint to construct and log a message.
func (l *LogService) Debug(args ...interface{}) {
	l.debug(l.base, args...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, DebugLevel, l.rabbitMQ)
	}
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *LogService) Debugf(message string, args ...interface{}) {
	l.debugf(l.base, message, args...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, DebugLevel, l.rabbitMQ)
	}
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func (l *LogService) Debugw(message string, keysAndValues ...interface{}) {
	l.debugw(l.base, message, keysAndValues...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, DebugLevel, l.rabbitMQ)
	}
}

// Info uses fmt.Sprint to construct and log a message.
func (l *LogService) Info(args ...interface{}) {
	l.info(l.base, args...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, InfoLevel, l.rabbitMQ)
	}
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *LogService) Infof(message string, args ...interface{}) {
	l.infof(l.base, message, args...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, InfoLevel, l.rabbitMQ)
	}
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *LogService) Infow(message string, keysAndValues ...interface{}) {
	l.infow(l.base, message, keysAndValues...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, InfoLevel, l.rabbitMQ)
	}
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *LogService) Warn(args ...interface{}) {
	l.warn(l.base, args...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, WarnLevel, l.rabbitMQ)
	}
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *LogService) Warnf(message string, args ...interface{}) {
	l.warnf(l.base, message, args...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, WarnLevel, l.rabbitMQ)
	}
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *LogService) Warnw(message string, keysAndValues ...interface{}) {
	l.warnw(l.base, message, keysAndValues...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, WarnLevel, l.rabbitMQ)
	}
}

// Error uses fmt.Sprint to construct and log a message.
func (l *LogService) Error(args ...interface{}) {
	l.error(l.base, args...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, ErrorLevel, l.rabbitMQ)
	}
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *LogService) Errorf(message string, args ...interface{}) {
	l.errorf(l.base, message, args...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, ErrorLevel, l.rabbitMQ)
	}
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *LogService) Errorw(message string, keysAndValues ...interface{}) {
	l.errorw(l.base, message, keysAndValues...)
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, ErrorLevel, l.rabbitMQ)
	}
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *LogService) DPanic(args ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, DPanicLevel, l.rabbitMQ)
	}
	l.dpanic(l.base, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *LogService) DPanicf(message string, args ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, DPanicLevel, l.rabbitMQ)
	}
	l.dpanicf(l.base, message, args...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (l *LogService) DPanicw(message string, keysAndValues ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, DPanicLevel, l.rabbitMQ)
	}
	l.dpanicw(l.base, message, keysAndValues...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *LogService) Panic(args ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, PanicLevel, l.rabbitMQ)
	}
	l.panic(l.base, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *LogService) Panicf(message string, args ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, PanicLevel, l.rabbitMQ)
	}
	l.panicf(l.base, message, args...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (l *LogService) Panicw(message string, keysAndValues ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, PanicLevel, l.rabbitMQ)
	}
	l.panicw(l.base, message, keysAndValues...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *LogService) Fatal(args ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, FatalLevel, l.rabbitMQ)
	}
	l.fatal(l.base, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *LogService) Fatalf(message string, args ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, FatalLevel, l.rabbitMQ)
	}
	l.fatalf(l.base, message, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (l *LogService) Fatalw(message string, keysAndValues ...interface{}) {
	if l.centralized {
		centralizedLog(l.jsonWriter, l.serviceCode, l.serviceName, FatalLevel, l.rabbitMQ)
	}
	l.fatalw(l.base, message, keysAndValues...)
}
