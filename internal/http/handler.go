package http

import (
	"html/template"
	"net/http"
)

type Handler struct {
	template *template.Template
}

func NewHandler() *Handler {
	return &Handler{
		template: template.Must(template.ParseFS(htmlFiles, "templates/*")),
	}
}

func (h *Handler) Poll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Trigger", "getupdate")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello"))
}

func (h *Handler) GetUpdate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Got updated"))
}

func (h *Handler) GetRedirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Redirect", "/bye")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Bye(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bye"))
}
