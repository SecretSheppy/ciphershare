package handlers

import (
	"html/template"
	"log/slog"
)

// Handlers wraps all the objects needed by each individual Handler into one object. All the individual handlers are
// then attached to this object.
type Handlers struct {
	tpl *template.Template
	log *slog.Logger
}

// NewHandlers creates a new Handlers object.
func NewHandlers(tpl *template.Template, logger *slog.Logger) *Handlers {
	return &Handlers{
		tpl: tpl,
		log: logger,
	}
}
