package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"golang-encrypted-filesharing/handlers"
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
	coll, client := mongodb.Connect()
	defer func() {
		mongodb.Disconnect(client)
	}()
	//example of finding entity and getting entity
	mongodb.FindEntityViaEmail(coll, "owalmsley1@sheffield.ac.uk")
	emails := []string{"owalmsley1@sheffield.ac.uk", "Rturner3@sheffield.ac.uk"}
	print(string(mongodb.CreateEntity(coll, emails, "this magic path", "This file key")))

	//Other stuff
	if err := godotenv.Load(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	logger.Info("Starting encrypted file sharing system...")

	tpl := template.Must(template.New("").ParseGlob("./templates/*.gohtml"))

	h := handlers.NewHandlers(tpl, logger)

	r := mux.NewRouter()
	r.HandleFunc("/", h.Upload).Methods("GET")
	r.HandleFunc("/download/{key}", h.Download).Methods("GET")
	r.HandleFunc("/complete", h.Complete).Methods("GET")
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
