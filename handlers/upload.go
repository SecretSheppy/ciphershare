package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) Upload(w http.ResponseWriter, r *http.Request) {
	err := h.tpl.ExecuteTemplate(w, "upload.gohtml", nil)
	if err != nil {
		fmt.Println(err)
	}
}
