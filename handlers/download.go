package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handlers) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := h.tpl.ExecuteTemplate(w, "download.gohtml", vars)
	if err != nil {
		fmt.Println(err)
	}
}
