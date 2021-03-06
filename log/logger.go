package log

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	log "github.com/ron96G/log15"
)

type contextKey string

var (
	Root          log.Logger
	logCtxKey     contextKey = "logger"
	loggerNameKey            = "logger"
	formats                  = map[string]log.Format{
		"json":   JsonFormat(),
		"logfmt": LogfmtFormat(),
	}
	defaultLoglevel  = log.LvlInfo
	defaultLogformat = formats["logfmt"]

	TimeFormat = time.RFC3339
	LevelKey   = "level"
	TimeKey    = "time"
	MsgKey     = "message"
)

func init() {
	Reset()
}

func Configure(loglevel, format string, output io.Writer, params ...interface{}) {
	level, err := log.LvlFromString(strings.ToLower(loglevel))
	if err != nil {
		Root.Crit("unknown loglevel", "loglevel", loglevel)
		level = defaultLoglevel
	}

	var frmt log.Format
	var ok bool
	if frmt, ok = formats[format]; !ok {
		Root.Crit("unable to find log format", "format", format)
		frmt = defaultLogformat
	}

	Root.SetHandler(log.MultiHandler(
		log.CallerFileHandler(log.DiscardHandler()),
		log.LvlFilterHandler(level, log.StreamHandler(output, frmt)),
	))

	Root = Root.New(params...)
}

func Reset() {
	log.TimeFormat = TimeFormat
	log.TimeKey = TimeKey
	log.LvlKey = LevelKey
	log.MsgKey = MsgKey
	Root = log.Root()
}

func New(logger string, ctx ...interface{}) Logger {
	params := append([]interface{}{loggerNameKey, logger}, ctx...)
	return Root.New(params...)
}

func ToContext(ctx context.Context, logger Logger, params ...interface{}) interface{} {
	l := logger.New(params...)
	return context.WithValue(ctx, logCtxKey, l)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(logCtxKey).(Logger)
	if !ok {
		return Root.New()
	}
	return logger
}

func Trace(msg string, v ...interface{}) {
	Root.Debug(msg, v...)
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

func Debug(msg string, v ...interface{}) {
	Root.Debug(msg, v...)
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

func Info(msg string, v ...interface{}) {
	Root.Info(msg, v...)
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
