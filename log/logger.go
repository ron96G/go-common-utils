package log

import (
	"context"
	"fmt"
	"io"
	"strings"

	log "github.com/inconshreveable/log15"
)

var (
	Root log.Logger
)

func init() {
	Root = log.Root()
}

func Configure(loglevel string, output io.Writer) {
	level, err := log.LvlFromString(strings.ToLower(loglevel))
	if err != nil {
		Root.Error("unknown loglevel", "loglevel", loglevel)
		level = log.LvlInfo
	}

	Root.SetHandler(log.MultiHandler(
		log.CallerFileHandler(log.DiscardHandler()),
		log.LvlFilterHandler(level, log.StreamHandler(output, log.LogfmtFormat())),
	))
}

func Reset() {
	Root = log.Root()
}

func New(logger string, ctx ...interface{}) Logger {
	params := append([]interface{}{"logger", logger}, ctx...)
	return Root.New(params...)
}

func ToContext(ctx context.Context, logger Logger, params ...interface{}) context.Context {
	l := logger.New(params...)
	return context.WithValue(ctx, "logger", l)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value("logger").(Logger)
	if !ok {
		return Root.New()
	}
	return logger
}

func Tracef(format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}

	Root.Debug(message)
}

func Debugf(format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}

	Root.Debug(message)
}

func Infof(format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}

	Root.Info(message)
}

func Warn(msg string, v ...interface{}) {
	Root.Warn(msg, v...)
}

func Warnf(format string, v ...interface{}) {
	var message string
	if len(v) > 0 {
		message = fmt.Sprintf(format, v...)
	} else {
		message = format
	}

	Root.Warn(message)
}

func Error(msg string, args ...interface{}) {
	Root.Error(msg, args...)
}

func Errorf(skip int, format string, v ...interface{}) {
	Root.Error(fmt.Sprintf(format, v...))
}
