package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) Download(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
