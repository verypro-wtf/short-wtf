package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/verypro-wtf/short-wtf/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Log() *zerolog.Event
	Fatal() *zerolog.Event
	Err(err error) *zerolog.Event
	Panic() *zerolog.Event
	Error() *zerolog.Event
	Warn() *zerolog.Event
	Info() *zerolog.Event
	Trace() *zerolog.Event
	Debug() *zerolog.Event
	With() zerolog.Context
	SetLogLevel(level string)
	Printf(format string, v ...any)
	Print(v ...any)
}

type logger struct {
	logger zerolog.Logger
}

var sharedLogger *logger
var globalConfig config.LoggerConfig
var globalEnv string

func (logger *logger) SetLogLevel(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)
}

func Log() *zerolog.Event {
	return getDefaultLogger().logger.Log()
}

func Fatal() *zerolog.Event {
	return getDefaultLogger().logger.Fatal()
}

func Error() *zerolog.Event {
	return getDefaultLogger().logger.Error()
}

func Err(err error) *zerolog.Event {
	return getDefaultLogger().logger.Err(err)
}

func Warn() *zerolog.Event {
	return getDefaultLogger().logger.Warn()
}

func Info() *zerolog.Event {
	return getDefaultLogger().logger.Info()
}

func Debug() *zerolog.Event {
	return getDefaultLogger().logger.Debug()
}

func Panic() *zerolog.Event {
	return getDefaultLogger().logger.Panic()
}

func Trace() *zerolog.Event {
	return getDefaultLogger().logger.Trace()
}

func Print(v ...interface{}) {
	getDefaultLogger().logger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	getDefaultLogger().logger.Printf(format, v...)
}

func With() zerolog.Context {
	return getDefaultLogger().logger.With().Caller()
}

func getDefaultLogger() *logger {
	if sharedLogger == nil {
		panic("Logger: Package not initialized. Please call Init() first!")
	}
	return sharedLogger
}

func Init(config config.LoggerConfig, env string) {
	globalConfig = config
	globalEnv = env

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	lvl, err := zerolog.ParseLevel(config.Log_Level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)

	sharedLogger = &logger{
		logger: createLogger("main", config.Log_File_Path),
	}
}

func New(serviceName string) zerolog.Logger {
	if globalConfig.Log_File_Path == "" {
		return createLogger(serviceName, "")
	}

	dir := filepath.Dir(globalConfig.Log_File_Path)
	ext := filepath.Ext(globalConfig.Log_File_Path)
	base := filepath.Base(globalConfig.Log_File_Path)
	extRemovedBase := base[:len(base)-len(ext)]

	logPath := filepath.Join(dir, serviceName + "." + extRemovedBase + ext)


	return createLogger(serviceName, logPath)
}

func createLogger(serviceName, logPath string) zerolog.Logger {
	var writers []io.Writer

	if globalEnv == "dev" {
		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})
	} else {
		writers = append(writers, os.Stderr)
	}

	if logPath != "" {
		writers = append(writers, &lumberjack.Logger{
			Filename:   logPath,
			Compress:   globalConfig.Compress_Logs,
			MaxBackups: globalConfig.Max_Backups,
			MaxSize:    globalConfig.Max_Size,
			MaxAge:     globalConfig.Max_Age,
		})
	}

	return zerolog.New(io.MultiWriter(writers...)).
		With().
		Timestamp().
		Stack().
		Str("service", serviceName).
		Logger()
}
