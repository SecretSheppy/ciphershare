package handlers

import (
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
	"golang-encrypted-filesharing/authentication"
	"html/template"
	"log/slog"
	"os"
)

// Handlers wraps all the objects needed by each individual Handler into one object. All the individual handlers are
// then attached to this object.
type Handlers struct {
	tpl        *template.Template
	log        *slog.Logger
	collection *mongo.Collection
	store      *sessions.CookieStore
	auth       *authentication.Auth
}

func newHandlersAuth(logger *slog.Logger) *authentication.Auth {
	return authentication.New(logger, &authentication.Auth0Config{
		Auth0Domain:  os.Getenv("AUTH0_DOMAIN"),
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
	})
}

// NewHandlers creates a new Handlers object.
func NewHandlers(tpl *template.Template, logger *slog.Logger, collection *mongo.Collection, store *sessions.CookieStore) *Handlers {
	return &Handlers{
		tpl:        tpl,
		log:        logger,
		collection: collection,
		store:      store,
		auth:       newHandlersAuth(logger),
	}
}
