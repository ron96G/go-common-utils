package log

import (
	log "github.com/ron96G/log15"
)

type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(ctx ...interface{}) log.Logger

	// GetHandler gets the handler associated with the logger.
	GetHandler() log.Handler

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h log.Handler)

	// Log a message at the given level with context key/value pairs
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})
}
