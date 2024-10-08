package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	errs "github.com/Clankyyy/scheduler/internal/errors"
	"github.com/Clankyyy/scheduler/internal/schedule"
	"github.com/Clankyyy/scheduler/internal/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type APIserver struct {
	listenAddr string
	storage    storage.Storager
}

var (
	c = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_sample_metric",
		Help: "Sample metric for app",
	})

	h = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "app_sample_histogram",
		Help: "Sample histogram",
	})
	d = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_sample_devices",
		Help: "Sample counter opts devices"}, []string{"device"})
)

func UpdateMEtrics() {
	for {
		rand.New(rand.NewSource(time.Now().UnixNano()))
		h.Observe(float64(rand.Intn(100-0+1) + 0))
		d.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
		c.Inc()
		fmt.Print(".")
		time.Sleep(3 * time.Second)
	}
}

func (s *APIserver) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /groups/", makeHTTPHandleFunc(s.handleGetGroups))

	mux.HandleFunc("GET /schedule/daily/{slug}", makeHTTPHandleFunc(s.handleGetDaily))

	mux.HandleFunc("POST /schedule/weekly", makeHTTPHandleFunc(s.handleCreateWeekly))
	mux.HandleFunc("GET /schedule/weekly/full/{slug}", makeHTTPHandleFunc(s.handleGetFullWeeklyBySlug))
	mux.HandleFunc("GET /schedule/weekly/{slug}", makeHTTPHandleFunc(s.handleGetWeeklyBySlug))

	mux.HandleFunc("PUT /schedule/weekly/{slug}", makeHTTPHandleFunc(s.handleUpdateWeekly))
	mux.HandleFunc("DELETE /schedule/weekly/{slug}", makeHTTPHandleFunc(s.handleDeleteWeeklyBySlug))

	mux.HandleFunc("GET /ping", makeHTTPHandleFunc(s.handlePing))
	mux.Handle("GET /metrics", promhttp.Handler())
	log.Println("API running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIserver) handleGetDaily(w http.ResponseWriter, r *http.Request) error {
	weekday, err1 := schedule.BuildWeekday(r.URL.Query().Get("day"))
	scheduleType, err2 := schedule.BuildScheduleType(r.URL.Query().Get("type"))
	if err1 != nil || err2 != nil {
		err := errs.NewAPIError(http.StatusBadRequest, "bad parameter format")
		return WriteJSON(w, err.StatusCode, err)
	}
	slug := r.PathValue("slug")

	daily, err := s.storage.GetDailyBySlug(slug, weekday, scheduleType)

	if err != nil {
		err := errs.NewAPIError(http.StatusNotFound, fmt.Sprintf("Group %s schedule dont exist", slug))
		return WriteJSON(w, err.StatusCode, err)
	}

	return WriteJSON(w, http.StatusOK, daily)
}

func (s *APIserver) handlePing(w http.ResponseWriter, r *http.Request) error {
	dur := time.Duration(2) * time.Second
	fmt.Println("pinged")
	time.Sleep(dur)
	return WriteJSON(w, http.StatusOK, dur)
}

func (s *APIserver) handleGetGroups(w http.ResponseWriter, r *http.Request) error {
	groups, err := s.storage.GetGroups()
	groupsStr := make([]string, 0, len(groups))
	for _, v := range groups {
		slug := strconv.Itoa(v.Course) + "-" + v.Name
		groupsStr = append(groupsStr, slug)
	}
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
	}
	return WriteJSON(w, http.StatusOK, groupsStr)
}

func (s *APIserver) handleCreateWeekly(w http.ResponseWriter, r *http.Request) error {
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

func (s *APIserver) handleGetWeeklyBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("slug")
	scheduleTypeStr := r.URL.Query().Get("type")
	scheduleType, err := schedule.BuildScheduleType(scheduleTypeStr)
	if err != nil {
		err := errs.NewAPIError(http.StatusBadRequest, "Incorrect parameter format")
		return WriteJSON(w, err.StatusCode, err)
	}

	weekly, err := s.storage.GetWeeklyBySlug(slug, scheduleType)
	if err != nil {
		err := errs.NewAPIError(http.StatusInternalServerError, "Unable to get schedule")
		return WriteJSON(w, err.StatusCode, err)
	}

	return WriteJSON(w, http.StatusOK, weekly)
}

func (s *APIserver) handleGetFullWeeklyBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("slug")

	weekly, err := s.storage.GetFullWeeklyBySlug(slug)
	if err != nil {
		err := errs.NewAPIError(http.StatusInternalServerError, "Unable to get schedule")
		return WriteJSON(w, err.StatusCode, err)
	}

	return WriteJSON(w, http.StatusOK, weekly)
}

func (s *APIserver) handleUpdateWeekly(w http.ResponseWriter, r *http.Request) error {
	newSchedule := make([]schedule.Weekly, 0, 2)
	slug := r.PathValue("slug")
	if err := json.NewDecoder(r.Body).Decode(&newSchedule); err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}
	err := s.storage.UpdateWeeklyBySlug(slug, newSchedule)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}

	return WriteJSON(w, http.StatusCreated, nil)
}

func (s *APIserver) handleDeleteWeeklyBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("slug")
	return s.storage.DeleteSchedule(slug)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return nil
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	return e.Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			apiErr := errs.APIError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Service unable to process request",
				Details:    "Unexpected error",
			}
			WriteJSON(w, http.StatusBadRequest, apiErr)
		}
	}
}

func NewAPIServer(listenAddr string, storage storage.Storager) *APIserver {
	return &APIserver{
		listenAddr: listenAddr,
		storage:    storage,
	}
}
