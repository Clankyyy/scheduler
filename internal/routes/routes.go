package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Clankyyy/scheduler/internal/storage"
)

type APIserver struct {
	listenAddr string
	storage    storage.Storager
}

func NewAPIServer(listenAddr string, storage storage.Storager) *APIserver {
	return &APIserver{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *APIserver) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /schedule/{id}", makeHTTPHandleFunc(s.handleGetSchedule))
	mux.HandleFunc("POST /schedule/", makeHTTPHandleFunc(s.handleCreateSchedule))
	mux.HandleFunc("DELETE /schedule/", makeHTTPHandleFunc(s.handleDeleteSchedule))

	log.Println("API running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIserver) handleCreateSchedule(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIserver) handleGetSchedule(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	weekly, err := s.storage.GetSchedule(id)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	return WriteJSON(w, http.StatusOK, weekly)
}

func (s *APIserver) handleDeleteSchedule(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	return e.Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
