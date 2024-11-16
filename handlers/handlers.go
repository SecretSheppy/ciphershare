package handlers

import (
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"log/slog"
)

// Handlers wraps all the objects needed by each individual Handler into one object. All the individual handlers are
// then attached to this object.
type Handlers struct {
	tpl        *template.Template
	log        *slog.Logger
	collection *mongo.Collection
}

// NewHandlers creates a new Handlers object.
func NewHandlers(tpl *template.Template, logger *slog.Logger, collection *mongo.Collection) *Handlers {
	return &Handlers{
		tpl:        tpl,
		log:        logger,
		collection: collection,
	}
}
