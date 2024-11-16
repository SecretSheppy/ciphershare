package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang-encrypted-filesharing/cryptography"
	"net/http"
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

func (h *Handlers) DownloadFile(w http.ResponseWriter, r *http.Request) {
	var encrypted []byte = []byte("lemons")
	var key = "lemons"

	plaintext := cryptography.Decrypt(key, encrypted)

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
