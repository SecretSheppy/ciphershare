package handlers

import (
	"encoding/json"
	"golang-encrypted-filesharing/mongodb"
	"golang-encrypted-filesharing/templates"
	"golang-encrypted-filesharing/utils"
	"net/http"
)

type AuthTokenPage struct {
	Key      string
	Email    string
	OTPError bool
}

func NewAuthTokenPage(key, email string, OPTError bool) *AuthTokenPage {
	return &AuthTokenPage{
		Key:      key,
		Email:    email,
		OTPError: OPTError,
	}
}

func (h *Handlers) Authenticate(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("key")
	err, rawRecord := mongodb.FindEntityViaUuid(h.collection, id)
	if err != nil {
		h.log.Error(err.Error())
	}

	record := make(map[string]json.RawMessage)
	err = json.Unmarshal(rawRecord, &record)
	if err != nil {
		h.log.Error("failed to unmarshal record", err)
	}

	var emails []string
	err = json.Unmarshal(record["list_of_emails"], &emails)
	if err != nil {
		h.log.Error("failed to unmarshal emails", err)
	}

	if utils.StringInList(r.FormValue("email"), emails) {
		h.auth.SendVerificationEmail(r.FormValue("email"))
		template := NewAuthTokenPage(id, r.FormValue("email"), false)
		err = h.tpl.ExecuteTemplate(w, "auth0_token_authentication.gohtml", template)
		if err != nil {
			h.log.Error("failed to execute template", err)
		}
		return
	}

	template := templates.NewAuth0EmailLogin(id, true)
	err = h.tpl.ExecuteTemplate(w, "auth0_email_login.gohtml", template)
	if err != nil {
		h.log.Error("failed to execute template: ", err)
	}
}

func (h *Handlers) AuthToken(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	id := r.FormValue("key")
	otp := r.FormValue("otp")

	err := h.auth.ValidateOauthToken(email, otp)
	if err != nil {
		template := NewAuthTokenPage(id, email, true)
		err = h.tpl.ExecuteTemplate(w, "auth0_token_authentication.gohtml", template)
		if err != nil {
			h.log.Error("failed to execute template", err)
		}
		return
	}

	session, _ := h.store.Get(r, "authenticated")

	session.Values["id"] = id
	session.Options.HttpOnly = true
	session.Options.SameSite = http.SameSiteStrictMode
	session.Options.MaxAge = 600 // 10 minute session timer

	err = session.Save(r, w)
	if err != nil {
		h.log.Error("failed to save session", err)
		return
	}

	http.Redirect(w, r, "/files/"+id, http.StatusFound)
	return
}
