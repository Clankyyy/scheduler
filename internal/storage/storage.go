package storage

import (
	"encoding/json"
	"os"

	"github.com/Clankyyy/scheduler/internal/schedule"
)

const postfix = ".json"

type Storager interface {
	CreateSchedule(g *schedule.Group) error
	DeleteSchedule(string) error
	GetSchedule(string) (*schedule.Weekly, error)
}

type FSStorage struct {
	path string
}

func (fss FSStorage) GetSchedule(id string) (*schedule.Weekly, error) {
	f, err := os.Open(fss.path + id + postfix)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var w schedule.Weekly

	if err := json.NewDecoder(f).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (fss FSStorage) CreateSchedule(g *schedule.Group) error {
	return nil
}

func (fss FSStorage) DeleteSchedule(id string) error {
	return nil
}
func NewFSStorage(path string) *FSStorage {
	return &FSStorage{
		path: path,
	}
}

//home\clanky\projects\scheduler\data\spbgti
