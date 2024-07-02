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
	GetScheduleBySlug(string) ([]schedule.Weekly, error)
	GetSchedule() ([]ScheduleResponse, error)
	UpdateScheduleBySlug([]schedule.Weekly, string) error
}

type FSStorage struct {
	path string
}

func (fss FSStorage) GetScheduleBySlug(slug string) ([]schedule.Weekly, error) {
	f, err := os.Open(fss.path + slug + postfix)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	w := make([]schedule.Weekly, 2)

	if err := json.NewDecoder(f).Decode(&w); err != nil {
		return nil, err
	}
	return w, nil
}
func (fss FSStorage) UpdateScheduleBySlug(s []schedule.Weekly, slug string) error {
	path := fss.path + slug + postfix
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "   ")

	if err := e.Encode(s); err != nil {
		return err
	}

	return nil
}

func (fss FSStorage) GetSchedule() ([]ScheduleResponse, error) {
	files, err := os.ReadDir(fss.path)
	respose := make([]ScheduleResponse, 0, len(files))
	if err != nil {
		return respose, err
	}
	for _, file := range files {
		cleanName := strings.Split(file.Name(), ".")[0]
		details := strings.Split(cleanName, "-")
		if len(details) == 2 {
			respose = append(respose, ScheduleResponse{details[0], details[1]})
		}
	}
	return respose, nil
}

func (fss FSStorage) CreateGroupSchedule(g *schedule.Group) error {

	filename := fss.path + fmt.Sprint(g.Course) + "-" + g.Name + "-" + postfix
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "   ")
	if err := e.Encode(g.Schedule); err != nil {
		return err
	}
	return nil
}

func (fss FSStorage) DeleteSchedule(slug string) error {
	err := os.Remove(fss.path + slug + postfix)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

func NewFSStorage(path string) *FSStorage {
	return &FSStorage{
		path: path,
	}
}

type ScheduleResponse struct {
	Course string `json:"course"`
	Name   string `json:"name"`
}
