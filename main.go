package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"golang-encrypted-filesharing/handlers"
	"golang-encrypted-filesharing/middleware"
	"golang-encrypted-filesharing/mongodb"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	tlsPort     = ":8080"
	certificate = "server/ssl/cert.pem"
	key         = "server/ssl/key.pem"
)

func main() {
	//Connect to the database
	coll, client := mongodb.Connect()
	defer mongodb.Disconnect(client)

	//Other stuff
	if err := godotenv.Load(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	logger.Info("Starting encrypted file sharing system...")

	tpl := template.Must(template.New("").ParseGlob("./templates/*.gohtml"))

	store := sessions.NewCookieStore([]byte("example-key"))

	h := handlers.NewHandlers(tpl, logger, coll, store)
	m := middleware.New(logger)

	r := mux.NewRouter()
	r.Use(m.Logger)
	r.HandleFunc("/", h.Upload).Methods("GET")
	r.HandleFunc("/", h.UploadFile).Methods("POST")
	r.HandleFunc("/files/{key}", h.Download).Methods("GET")
	r.HandleFunc("/download/{id}", h.DownloadFile).Methods("GET")
	r.HandleFunc("/complete/{id}", h.Complete).Methods("GET")
	r.HandleFunc("/auth", h.Authenticate).Methods("POST")
	r.HandleFunc("/auth/token", h.AuthToken).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         tlsPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServeTLS(certificate, key)
	if err != nil {
		logger.Error(err.Error())
	}
}
