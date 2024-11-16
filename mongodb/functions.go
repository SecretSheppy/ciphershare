package mongodb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// Connect return the database connection
func Connect() (*mongo.Collection, *mongo.Client) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	coll := client.Database("security_go").Collection("go_project")
	return coll, client
}
func Disconnect(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func FindEntityViaEmail(collection *mongo.Collection, emailToFind string) []byte {
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.D{{"list_of_emails", emailToFind}}).Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No document was found with the name %s\n", emailToFind)
		return nil
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
	return jsonData
}
func FindEntityViaUuid(collection *mongo.Collection, uuidToFind string) (error, []byte) {
	var result bson.M
	objectId, idErr := primitive.ObjectIDFromHex(uuidToFind)
	if idErr != nil {
		return idErr, nil
	}
	err := collection.FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No document was found with the name %s\n", uuidToFind)
		return err, nil
	}
	if err != nil {
		return err, nil
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return err, nil
	}
	fmt.Println(string(jsonData))
	return nil, jsonData
}
func CreateEntity(collection *mongo.Collection, listOfEmails []string, pathToEncryptedFile string, fileKey string) []byte {
	doc := bson.D{
		{"list_of_emails", listOfEmails},
		{"path_to_encrypted_file", pathToEncryptedFile},
		{"encrypted_file_key", fileKey},
	}
	result, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted a single document: ", result.InsertedID)
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	return jsonData
}
