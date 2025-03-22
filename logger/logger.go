package logger

import (
	"errors"
	"os"

	cmlog "github.com/cometbft/cometbft/libs/log"
	ipfslog "github.com/ipfs/go-log/v2"
)

var ErrMissingValue = errors.New("(MISSING)")

type MeLogger struct {
	name    string
	Logger  *ipfslog.ZapEventLogger
	context []any
}

func (l MeLogger) Debug(msg string, keyvals ...interface{}) {
	l.Logger.Debugw(msg, append(l.context, keyvals...)...)
}
func (l MeLogger) Info(msg string, keyvals ...interface{}) {
	l.Logger.Infow(msg, append(l.context, keyvals...)...)
}
func (l MeLogger) Error(msg string, keyvals ...interface{}) {
	l.Logger.Errorw(msg, append(l.context, keyvals...)...)
}

func (l MeLogger) With(keyvals ...interface{}) cmlog.Logger {

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, ErrMissingValue)
	}
	return MeLogger{
		Logger:  l.Logger,
		context: append(l.context, keyvals...),
	}
}
func (l MeLogger) WithStacktrace(traceLevel ipfslog.LogLevel) MeLogger {
	return MeLogger{
		Logger:  ipfslog.WithStacktrace(l.Logger, traceLevel),
		context: l.context,
	}
}

func (l MeLogger) WithLevel(level string) MeLogger {
	ipfslog.SetLogLevel(l.name, level) //nolint:errcheck
	return l
}

func (l MeLogger) WithEnvLevelOr(level string) MeLogger {
	lv := os.Getenv("GOLOG_LOG_LEVEL")
	if lv != "" {
		err := ipfslog.SetLogLevel(l.name, lv)
		if err != nil {
			l.Logger.Errorf("set %s logger's level failed:%s", l.name, err)
		}
		return l
	}
	err := ipfslog.SetLogLevel(l.name, level)
	if err != nil {
		l.Logger.Errorf("set %s logger's level failed:%s", l.name, err)
	}
	return l
}
func NewLogger(name string) MeLogger {
	l := ipfslog.Logger(name)
	return MeLogger{
		name:   name,
		Logger: ipfslog.WithSkip(l, 1),
	}
}

func SetAllLoggerLevel(level string) error {
	logv, err := ipfslog.LevelFromString(level)
	if err != nil {
		return err
	}
	ipfslog.SetAllLoggers(logv)
	return nil
}
