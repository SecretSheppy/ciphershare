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
		http.Redirect(w, r, "/?error=1", http.StatusSeeOther)
		return
	}
	emails := strings.Split(emailForm, ",")
	for _, email := range emails {
		if !validation.IsEmailValid(email) {
			h.log.Warn("Email is not valid")
			http.Redirect(w, r, "/?error=2", http.StatusSeeOther)
			return
		}
	}

	file, fileHeader, err := r.FormFile("fileUpload")
	if file == nil {
		h.log.Warn("File is empty")
		http.Redirect(w, r, "/?error=3", http.StatusSeeOther)
		return
	} else if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("File is being parsed")
	}
	key, path, err := saveFile(file, fileHeader, getFileName())
	if err != nil {
		h.log.Error(err.Error())
		// File hasn't been downloaded
		http.Redirect(w, r, "/?error=4", http.StatusSeeOther)
		return
	} else {
		h.log.Info("File has been downloaded")
	}

	err, uuidJson := mongodb.CreateEntity(h.collection, emails, path, key)
	if err != nil {
		h.log.Error(err.Error())
		http.Redirect(w, r, "/?error=5", http.StatusSeeOther)
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

	filebytes, err := io.ReadAll(file)

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
	filebytes = []byte(jsonString + string(filebytes))

	key, encryptedFileBytes := cryptography.Encrypt(filebytes)

	_, err = tempFile.Write(encryptedFileBytes)

	return key, filepath.Clean("files/" + filename), nil
}
