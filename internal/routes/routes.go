package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Clankyyy/scheduler/internal/schedule"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
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

type APIserver struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIserver {
	return &APIserver{
		listenAddr: listenAddr,
	}
}

func (s *APIserver) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /schedule/", makeHTTPHandleFunc(s.handleGetSchedule))
	mux.HandleFunc("POST /schedule/", makeHTTPHandleFunc(s.handleCreateSchedule))
	mux.HandleFunc("DELETE /schedule/", makeHTTPHandleFunc(s.handleDeleteSchedule))

	log.Println("API running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIserver) handleCreateSchedule(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIserver) handleGetSchedule(w http.ResponseWriter, r *http.Request) error {
	group := schedule.NewGroup("4305", 4, 2)
	return WriteJSON(w, http.StatusOK, group)
}

func (s *APIserver) handleDeleteSchedule(w http.ResponseWriter, r *http.Request) error {
	return nil
}
