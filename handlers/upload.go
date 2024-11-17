package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang-encrypted-filesharing/cryptography"
	"golang-encrypted-filesharing/mongodb"
	"golang-encrypted-filesharing/validation"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	_, b, _, _         = runtime.Caller(0)
	RootPath           = filepath.Join(filepath.Dir(b), "../")
	errorCodeToMessage = map[string]string{
		"1": "You haven't given any emails!",
		"2": "The emails provided weren't valid.",
		"3": "You haven't uploaded a file!",
		"4": "There was a problem uploading your file. Please try again.",
		"5": "There was a problem storing your file. Please try again.",
	}
)

func (h *Handlers) Upload(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Error string
	}{
		Error: errorCodeToMessage[r.URL.Query().Get("error")],
	}
	err := h.tpl.ExecuteTemplate(w, "upload.gohtml", data)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("Upload page accessed")
	}
}

func (h *Handlers) UploadError(w http.ResponseWriter, r *http.Request, uploadErr string) {
	data := struct {
		Error string
	}{
		Error: errorCodeToMessage[uploadErr],
	}
	err := h.tpl.ExecuteTemplate(w, "upload.gohtml", data)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("Upload page accessed with an error displayed")
	}
}

func (h *Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(500 << 20)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("Form is being parsed")
	}
	emailForm := r.FormValue("emails")
	if emailForm == "" {
		h.log.Warn("Email is empty")
		h.UploadError(w, r, "1")
		return
	}
	emailForm = strings.ReplaceAll(emailForm, " ", "")
	emails := strings.Split(emailForm, ",")
	for _, email := range emails {
		if !validation.IsEmailValid(email) {
			h.log.Warn("Email is not valid")
			h.UploadError(w, r, "2")
			return
		}
	}

	file, fileHeader, err := r.FormFile("fileUpload")
	if file == nil {
		h.log.Warn("File is empty")
		h.UploadError(w, r, "3")
		return
	} else if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File is being parsed")
	}

	key, path, err := saveFile(file, fileHeader, getFileName())
	if err != nil {
		h.log.Error(err.Error())
		h.UploadError(w, r, "4")
		return
	} else {
		h.log.Info("File has been downloaded")
	}

	err, uuidJson := mongodb.CreateEntity(h.collection, emails, path, key)
	if err != nil {
		h.log.Error(err.Error())
		h.UploadError(w, r, "5")
		return
	} else {
		h.log.Info("Entity created")
	}
	h.log.Info("File is successfully uploaded")

	jsonPointer := make(map[string]json.RawMessage)
	err = json.Unmarshal(uuidJson, &jsonPointer)
	if err != nil {
		h.log.Error(err.Error())
	}
	uuid := string(jsonPointer["InsertedID"])
	uuid = uuid[1 : len(uuid)-1]

	data := struct {
		Id     string
		Domain string
	}{
		Id:     uuid,
		Domain: os.Getenv("DOMAIN"),
	}
	err = h.tpl.ExecuteTemplate(w, "complete.gohtml", data)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("Upload completion page accessed")
	}
}

// getFileName creates a random string
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

	_, err = tempFile.Write(encryptedFileBytes)

	return key, filepath.Clean("files/" + filename), nil
}
