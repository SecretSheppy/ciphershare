package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang-encrypted-filesharing/cryptography"
	"golang-encrypted-filesharing/mongodb"
	"io"
	"log"
	"net/http"
	"os"
)

func (h *Handlers) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err, _ := mongodb.FindEntityViaUuid(h.collection, vars["key"])
	if err != nil {
		h.log.Warn("Invalid download page ID")
		h.log.Error(err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	err = h.tpl.ExecuteTemplate(w, "download.gohtml", vars)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("ID " + vars["key"] + " download page accessed")
	}
}

func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err, jsonData := mongodb.FindEntityViaUuid(h.collection, id)
	if err != nil {
		h.log.Warn("Invalid download ID attempted post")
	}

	jsonPointer := make(map[string]json.RawMessage)
	err = json.Unmarshal(jsonData, &jsonPointer)
	if err != nil {
		h.log.Error(err.Error())
	}

	key := jsonPointer["encrypted_file_key"]
	encryptedPath := jsonPointer["path_to_encrypted_file"]
	// Open the file
	fmt.Println(os.Getwd())
	file, err := os.Open(string(encryptedPath)[1 : len(string(encryptedPath))-1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the entire content of the file into a byte slice
	encrypted, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	plaintext := cryptography.Decrypt(string(key)[1:len(string(key))-1], encrypted)

	metadata, content := splitPlainText(string(plaintext))

	// Set the headers to indicate a file download
	w.Header().Set("Content-Disposition", "attachment; filename="+metadata.FileName) // Set the filename

	// TODO: make dynamic
	if metadata.Extension == ".txt" {
		w.Header().Set("Content-Type", "text/plain")
	} else if metadata.Extension == ".pdf" {
		w.Header().Set("Content-Type", "application/pdf")
	}

	// Write the content to the response body, which will be downloaded as a file
	w.Write([]byte(content))

}

func splitPlainText(plaintext string) (MetaData, string) {
	// Extract JSON dynamically
	var metadata MetaData
	var jsonLength int

	// Find where JSON ends
	for i := 1; i <= len(plaintext); i++ {
		partial := plaintext[:i]
		err := json.Unmarshal([]byte(partial), &metadata)
		if err == nil {
			jsonLength = i
			break
		}
	}

	if jsonLength == 0 {
		fmt.Println("No valid JSON found in the message")
		return MetaData{}, ""
	}

	// Extract JSON and original message
	return metadata, plaintext[jsonLength:]
}
