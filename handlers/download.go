package handlers

import (
	"github.com/gorilla/mux"
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
