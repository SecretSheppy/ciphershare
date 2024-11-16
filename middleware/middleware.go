package middleware

import "log/slog"

type Middleware struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Middleware {
	return &Middleware{log}
}
