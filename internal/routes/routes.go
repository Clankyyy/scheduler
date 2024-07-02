package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Clankyyy/scheduler/internal/schedule"
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

	mux.HandleFunc("GET /schedule/{slug}", makeHTTPHandleFunc(s.handleGetScheduleBySlug))
	mux.HandleFunc("GET /schedule/", makeHTTPHandleFunc(s.handleGetSchedule))
	mux.HandleFunc("DELETE /schedule/{slug}", makeHTTPHandleFunc(s.handleDeleteScheduleBySlug))
	mux.HandleFunc("POST /schedule/", makeHTTPHandleFunc(s.handleCreateSchedule))

	log.Println("API running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIserver) handleCreateSchedule(w http.ResponseWriter, r *http.Request) error {
	g := &schedule.Group{}
	if err := json.NewDecoder(r.Body).Decode(g); err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}
	err := s.storage.CreateGroupSchedule(g)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}

	return WriteJSON(w, http.StatusCreated, *g)
}

func (s *APIserver) handleGetSchedule(w http.ResponseWriter, r *http.Request) error {
	groups, err := s.storage.GetSchedule()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
	}
	return WriteJSON(w, http.StatusOK, groups)
}

func (s *APIserver) handleGetScheduleBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("slug")
	weekly, err := s.storage.GetScheduleBySlug(slug)
	if err != nil {
		return WriteJSON(w, http.StatusNotFound, err.Error())
	}

	return WriteJSON(w, http.StatusOK, weekly)
}

func (s *APIserver) handleDeleteScheduleBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("slug")
	return s.storage.DeleteSchedule(slug)
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
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
