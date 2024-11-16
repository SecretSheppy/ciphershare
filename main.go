package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"golang-encrypted-filesharing/handlers"
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

	err := srv.ListenAndServeTLS(certificate, key)
	if err != nil {
		fmt.Println(err)
	}
}
