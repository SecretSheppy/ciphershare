package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang-encrypted-filesharing/handlers"
	"html/template"
	"log"
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
	//database stuff
	//if err := godotenv.Load(); err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("security_go").Collection("go_project")
	pathToEncryptedFile := "pathname"
	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"pathToEncryptedFile", pathToEncryptedFile}}).
		Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No document was found with the name %s\n", pathToEncryptedFile)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))

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
