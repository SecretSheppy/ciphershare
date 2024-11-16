package handlers

import (
	"net/http"
)

func (h *Handlers) Complete(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Key string
	}{
		Key: "asmoranomardicadaistinaculdacar",
	}
	err := h.tpl.ExecuteTemplate(w, "complete.gohtml", data)
	if err != nil {
		h.log.Error(err.Error())
	} else {
		h.log.Info("Upload completion page accessed")
	}
}
