package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"golang-encrypted-filesharing/handlers"
	"golang-encrypted-filesharing/mongodb"
	"html/template"
	"net/http"
	"time"
)

const (
	tlsPort     = ":8080"
	certificate = "server/ssl/cert.pem"
	key         = "server/ssl/key.pem"
)

func main() {
	coll := mongodb.Connect()
	//example of finding entity and getting entity
	mongodb.FindEntityViaEmail(coll, "Email1@gmail.com")
	emails := []string{"owalmsley1@sheffield.ac.uk", "Rturner3@sheffield.ac.uk"}
	print(string(mongodb.CreateEntity(coll, emails, "this magic path", "This file key")))

	//Other stuff
	tpl := template.Must(template.New("").ParseGlob("./templates/*.gohtml"))

	h := handlers.NewHandlers(tpl)

	r := mux.NewRouter()
	r.HandleFunc("/", h.Upload).Methods("GET")
	r.HandleFunc("/download/{key}", h.Download).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         tlsPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	nerr := srv.ListenAndServeTLS(certificate, key)
	if nerr != nil {
		fmt.Println(err)
	}
}
