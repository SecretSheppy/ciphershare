package handlers

import (
	"fmt"
	"golang-encrypted-filesharing/cryptography"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	RootPath   = filepath.Join(filepath.Dir(b), "../")
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
	err := r.ParseMultipartForm(500 << 20)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("Form is being parse")
	}

	file, handler, err := r.FormFile("fileUpload")
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File is being parsed")
	}
	err = saveFile(file, handler, getFileName())
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File has been downloaded")
	}

	// Encrypt the file
	key, encryptedFile := cryptography.Encrypt(r.FormValue("fileUpload"))

	// Upload file to folder

	// Upload path to database

	// Upload emails to database

}

func getFileName() string {
	return "name"
}

func saveFile(file multipart.File, handler *multipart.FileHeader, filename string) error {
	defer file.Close()
	folderPath := fmt.Sprintf("%s%s", RootPath, "/files")
	tempFile, err := os.CreateTemp(folderPath, filename)
	if err != nil {
		return err
	}
	defer tempFile.Close()
	filebytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	tempFile.Write(filebytes)
	return nil
}
