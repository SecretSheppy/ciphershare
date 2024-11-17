package handlers

import (
	"encoding/base64"
	"encoding/json"
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

	file, fileHeader, err := r.FormFile("fileUpload")
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File is being parsed")
	}

	key, path, err := saveFile(file, fileHeader, getFileName())
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File has been downloaded")
	}

	uuidJson := mongodb.CreateEntity(h.collection, strings.Split(r.FormValue("emails"), ","), path, key)
	h.log.Info("File is successfully uploaded")

	jsonPointer := make(map[string]json.RawMessage)
	err = json.Unmarshal(uuidJson, &jsonPointer)
	if err != nil {
		h.log.Error(err.Error())
	}
	uuid := string(jsonPointer["InsertedID"])
	uuid = uuid[1 : len(uuid)-1]

	http.Redirect(w, r, "/complete/"+uuid, http.StatusSeeOther)
}

func getFileName() string {
	return base64.URLEncoding.EncodeToString(cryptography.GenerateKey())
}

type MetaData struct {
	FileName  string
	Extension string
}

func saveFile(file multipart.File, fileHeader *multipart.FileHeader, filename string) (string, string, error) {
	defer file.Close()
	folderPath := fmt.Sprintf("%s%s", RootPath, "\\files")
	tempFile, err := os.Create(filepath.Clean(folderPath + "\\" + filename))
	if err != nil {
		return "", "", err
	}
	defer tempFile.Close()
	fileBytes, err := io.ReadAll(file)

	// Prepend metadata before encryption
	metadata := MetaData{
		FileName:  fileHeader.Filename,
		Extension: filepath.Ext(fileHeader.Filename),
	}

	// Serialize struct to JSON string
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return "", "", err
	}

	// Convert JSON to string
	jsonString := string(jsonData)

	// Prepend
	fileBytes = []byte(jsonString + string(fileBytes))

	key, encryptedFileBytes := cryptography.Encrypt(fileBytes)

	if err != nil {
		return "", "", err
	}
	tempFile.Write(encryptedFileBytes)

	return key, filepath.Clean("files/" + filename), nil
}
