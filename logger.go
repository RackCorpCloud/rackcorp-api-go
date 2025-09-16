package api

import (
	"log/slog"
	"reflect"
)

var DefaultLogger = slog.Default()

func withError(l *slog.Logger, err error) *slog.Logger {
	if err == nil {
		return l
	}
	return l.With(
		slog.String("exception.message", err.Error()),
		slog.String("exception.type", reflect.TypeOf(err).String()),
	)
}
