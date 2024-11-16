package handlers

import (
	"encoding/base64"
	"fmt"
	"golang-encrypted-filesharing/cryptography"
	"golang-encrypted-filesharing/mongodb"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
		h.log.Info("Form is being parsed")
	}

	file, _, err := r.FormFile("fileUpload")
	fmt.Println(file)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File is being parsed")
	}
	key, path, err := saveFile(file, getFileName())
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File has been downloaded")
	}

	mongodb.CreateEntity(h.collection, strings.Split(r.FormValue("emails"), ","), path, key)
}

func getFileName() string {

	return base64.URLEncoding.EncodeToString(cryptography.GenerateKey())
}

func saveFile(file multipart.File, filename string) (string, string, error) {
	defer file.Close()
	folderPath := fmt.Sprintf("%s%s", RootPath, "\\files")
	tempFile, err := os.Create(folderPath + "\\" + filename)
	if err != nil {
		return "", "", err
	}
	defer tempFile.Close()
	filebytes, err := io.ReadAll(file)

	key, encryptedFileBytes := cryptography.Encrypt(filebytes)

	if err != nil {
		return "", "", err
	}
	tempFile.Write(encryptedFileBytes)

	return key, "\\files\\" + filename, nil
}
