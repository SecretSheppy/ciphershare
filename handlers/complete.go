package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func (h *Handlers) Complete(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Id     string
		Domain string
	}{
		Id:     mux.Vars(r)["id"],
		Domain: os.Getenv("DOMAIN"),
	}
	err := h.tpl.ExecuteTemplate(w, "complete.gohtml", data)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("Upload completion page accessed")
	}
}
