package logger

import (
	"bufio"
	"bytes"
	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"path/filepath"
)

// NewLogger initialize new logger
func NewLogger(serviceCode int32, serviceName string, opt *Options) (*LogService, error) {
	var (
		err    error
		rabbit bool
	)
	option := GetDefaultOptions()
	option.LogLevel = opt.LogLevel

	// validate newLogger
	if option, err = validateNewLogger(serviceCode, serviceName, option, opt); err != nil {
		return nil, err
	}

	// validate rabbitMQ configuration and initialized RabbitMQ
	if opt.Centralized != nil && opt.Centralized.ServerAddress != "" {
		option.Centralized = opt.Centralized
		if err = validateRabbit(option.Centralized); err != nil {
			return nil, err
		}
		if err = initializeRabbitMQ(serviceName, option.Centralized); err != nil {
			return nil, err
		}
		rabbit = true
	}

	// jw is a json writer for get json encoded data from json core
	var b bytes.Buffer
	jw := bufio.NewWriter(&b)
	core := initZapCore(serviceName, option, jw)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	if option.Sentry {
		sentryCore, err := sentryCore(option.SentryDSN, serviceName, logger)
		if err != nil {
			return nil, err
		}
		logger = sentryCore
		log.Println("the sentry core for the logger has been initiated")
	}
	return &LogService{
		zap:         logger,
		base:        logger.Sugar(),
		serviceCode: serviceCode,
		serviceName: serviceName,
		rabbitMQ:    option.Centralized,
		jsonWriter:  &jsonWriter{jw, &b},
		centralized: rabbit,

		debug:   (*zap.SugaredLogger).Debug,
		debugf:  (*zap.SugaredLogger).Debugf,
		debugw:  (*zap.SugaredLogger).Debugw,
		info:    (*zap.SugaredLogger).Info,
		infof:   (*zap.SugaredLogger).Infof,
		infow:   (*zap.SugaredLogger).Infow,
		warn:    (*zap.SugaredLogger).Warn,
		warnf:   (*zap.SugaredLogger).Warnf,
		warnw:   (*zap.SugaredLogger).Warnw,
		error:   (*zap.SugaredLogger).Error,
		errorf:  (*zap.SugaredLogger).Errorf,
		errorw:  (*zap.SugaredLogger).Errorw,
		dpanic:  (*zap.SugaredLogger).DPanic,
		dpanicf: (*zap.SugaredLogger).DPanicf,
		dpanicw: (*zap.SugaredLogger).DPanicw,
		panic:   (*zap.SugaredLogger).Panic,
		panicf:  (*zap.SugaredLogger).Panicf,
		panicw:  (*zap.SugaredLogger).Panicw,
		fatal:   (*zap.SugaredLogger).Fatal,
		fatalf:  (*zap.SugaredLogger).Fatalf,
		fatalw:  (*zap.SugaredLogger).Fatalw,
	}, nil
}

func (l *LogService) GrpcMiddleware() grpc.UnaryServerInterceptor {
	opts := []grpc_zap.Option{
		grpc_zap.WithLevels(zapLevel),
	}
	grpc_zap.ReplaceGrpcLoggerV2(l.zap)
	return grpc_zap.UnaryServerInterceptor(l.zap, opts...)
}

func zapLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		return zap.DebugLevel
	}
	return grpc_zap.DefaultCodeToLevel(code)
}

func validateNewLogger(sCode int32, sName string, oldOpts *Options, newOpts *Options) (*Options, error) {
	if newOpts.Rotation != nil {
		oldOpts.Rotation = newOpts.Rotation
	}
	if !newOpts.ConsoleWriter {
		oldOpts.ConsoleWriter = newOpts.ConsoleWriter
	}
	if len(newOpts.LogPath) != 0 {
		if ok := filepath.IsAbs(newOpts.LogPath); ok != true {
			return nil, ERROR_FILEPATH_IS_NOT_ABS
		}
		oldOpts.LogPath = newOpts.LogPath
	}
	if !newOpts.Colorable {
		oldOpts.Colorable = newOpts.Colorable
	}
	if newOpts.Development {
		oldOpts.Development = newOpts.Development
	}
	if newOpts.Sentry {
		if len(newOpts.SentryDSN) == 0 {
			return nil, ERROR_SENTRY_DSN_IS_EMPTY
		}
		oldOpts.Sentry = newOpts.Sentry
		oldOpts.SentryDSN = newOpts.SentryDSN
	}
	if sCode == 0 {
		return nil, SERVICE_CODE_ERROR
	}
	if len(sName) == 0 {
		return nil, SERVICE_NAME_ERROR
	}

	return oldOpts, nil
}

func validateRabbit(rb *RabbitMQ) error {
	if len(rb.ServerAddress) == 0 {
		return RABBIT_SERVER_EMPTY_ERROR
	}
	if len(rb.UserName) == 0 {
		return RABBIT_USER_EMPTY_ERROR
	}
	if len(rb.Password) == 0 {
		return RABBIT_PASS_EMPTY_ERROR
	}
	if len(rb.LogLevel) == 0 {
		return RABBIT_LOG_LEVEL_ERROR
	}
	return nil
}

func sentryCore(dsn, serviceName string, logger *zap.Logger) (*zap.Logger, error) {
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		ServerName:       serviceName,
	})
	if err != nil {
		return nil, err
	}
	cfg := zapsentry.Configuration{
		Level: zapcore.ErrorLevel,
		Tags: map[string]string{
			"component": "system",
		},
		EnableBreadcrumbs: true,
		BreadcrumbLevel:   zapcore.InfoLevel,
		DisableStacktrace: false,
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(sentryClient))
	if err != nil {
		return nil, err
	}
	logger = logger.With(zapsentry.NewScope())
	return zapsentry.AttachCoreToLogger(core, logger), nil
}
