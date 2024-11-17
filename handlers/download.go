package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang-encrypted-filesharing/cryptography"
	"golang-encrypted-filesharing/mongodb"
	"io"
	"mime"
	"net/http"
	"os"
)

func (h *Handlers) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := h.tpl.ExecuteTemplate(w, "download.gohtml", vars)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("ID " + vars["key"] + " download page accessed")
	}
}

// Handles actually downloading the file onto the device
func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	jsonData := mongodb.FindEntityViaUuid(h.collection, id)

	jsonPointer := make(map[string]json.RawMessage)
	err := json.Unmarshal(jsonData, &jsonPointer)
	if err != nil {
		h.log.Error(err.Error())
	}

	key := jsonPointer["encrypted_file_key"]
	encryptedPath := jsonPointer["path_to_encrypted_file"]

	// Open the file
	file, err := os.Open(trimFirstLast(string(encryptedPath)))
	if err != nil {
		h.log.Error(err.Error())
	}
	defer file.Close()

	// Read the entire content of the file into a byte slice
	encrypted, err := io.ReadAll(file)
	if err != nil {
		h.log.Error(err.Error())
	}
	plaintext := cryptography.Decrypt(trimFirstLast(string(key)), encrypted)

	metadata, content := splitPlainText(plaintext)

	// Set the headers to indicate a file download
	w.Header().Set("Content-Disposition", "attachment; filename="+metadata.FileName) // Set the filename

	// Use the mime.TypeByExtension function to get the MIME type
	contentType := mime.TypeByExtension(metadata.Extension)

	// If the MIME type is not found, set it to "application/octet-stream"
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", contentType)

	// Write the content to the response body, which will be downloaded as a file
	w.Write([]byte(content))

	h.log.Info("ID " + id + " downloaded successfully")
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

func trimFirstLast(s string) string {
	return s[1 : len(string(s))-1]
}
