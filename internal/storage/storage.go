package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/Clankyyy/scheduler/internal/schedule"
)

const postfix = ".json"

type Storager interface {
	CreateGroupSchedule(g *schedule.Group) error
	DeleteSchedule(string) error
	GetWeeklyBySlug(string, schedule.ScheduleRequestParam) ([]schedule.Weekly, error)
	GetDailyBySlug(string) (schedule.Daily, error)
	GetGroups() ([]GroupResponse, error)
	UpdateWeeklyBySlug([]schedule.Weekly, string) error
}

type FSStorage struct {
	path string
}

func (fss FSStorage) GetDailyBySlug(slug string) (schedule.Daily, error) {

	return schedule.Daily{}, nil
}

func (fss FSStorage) GetWeeklyBySlug(slug string, param schedule.ScheduleRequestParam) ([]schedule.Weekly, error) {
	if param == schedule.Full {
		w, err := fss.getFullWeekly(slug)
		return w, err
	} else if param == schedule.Even || param == schedule.Odd {
		f, err := os.Open(fss.path + slug + "-" + param.String() + postfix)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		var w schedule.Weekly
		if err := json.NewDecoder(f).Decode(&w); err != nil {
			return nil, err
		}
		return []schedule.Weekly{w}, nil
	}

	return nil, fmt.Errorf("got some unexpected shit")
}

func (fss FSStorage) UpdateWeeklyBySlug(s []schedule.Weekly, slug string) error {
	for i, v := range s {
		path := fss.path + slug + "-" + v.EvenString() + postfix
		f, err := os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		e := json.NewEncoder(f)
		e.SetIndent("", "   ")

		if err := e.Encode(s[i]); err != nil {
			return err
		}
	}
	return nil
}

func (fss FSStorage) CreateGroupSchedule(g *schedule.Group) error {
	for i, v := range g.Schedule {
		evenStr := v.EvenString()
		path := fss.path + fmt.Sprint(g.Course) + "-" + g.Name + "-" + evenStr + postfix
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			return err
		}
		defer f.Close()
		e := json.NewEncoder(f)
		e.SetIndent("", "   ")
		if err := e.Encode(g.Schedule[i]); err != nil {
			return err
		}
	}

	return nil
}

func (fss FSStorage) GetGroups() ([]GroupResponse, error) {
	files, err := os.ReadDir(fss.path)
	respose := make([]GroupResponse, 0, len(files))
	if err != nil {
		return respose, err
	}
	for _, file := range files {
		cleanName := strings.Split(file.Name(), ".")[0]
		details := strings.Split(cleanName, "-")
		if len(details) == 2 {
			respose = append(respose, GroupResponse{details[0], details[1]})
		}
	}
	return respose, nil
}

func (fss FSStorage) DeleteSchedule(slug string) error {
	err := os.Remove(fss.path + slug + postfix)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

func (fss FSStorage) getFullWeekly(slug string) ([]schedule.Weekly, error) {
	statusList := []string{schedule.Even.String(), schedule.Odd.String()}
	w := make([]schedule.Weekly, 2)
	for i, v := range statusList {
		f, err := os.Open(fss.path + slug + "-" + v + postfix)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if err := json.NewDecoder(f).Decode(&w[i]); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func NewFSStorage(path string) *FSStorage {
	return &FSStorage{
		path: path,
	}
}

type GroupResponse struct {
	Course string `json:"course"`
	Name   string `json:"name"`
}
