package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) Upload(w http.ResponseWriter, r *http.Request) {
	err := h.tpl.ExecuteTemplate(w, "upload.gohtml", nil)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("Upload page accessed")
	}
}

func (h *Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File is being uploaded")
	}
	fmt.Println(r.FormValue("emails"))
	fmt.Println(r.FormValue("fileUpload"))

	// Encrypt the file
	// Upload file to folder
	// Upload path to database
	// Upload emails to database
}
