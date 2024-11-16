package handlers

import "html/template"

// Handlers wraps all the objects needed by each individual Handler into one object. All the individual handlers are
// then attached to this object.
type Handlers struct {
	tpl *template.Template
}

// NewHandlers creates a new Handlers object.
func NewHandlers(tpl *template.Template) *Handlers {
	return &Handlers{tpl}
}
