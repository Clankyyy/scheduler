package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Clankyyy/scheduler/internal/schedule"
)

const postfix = ".json"

type Storager interface {
	CreateGroupSchedule(g *schedule.Group) error
	DeleteSchedule(string) error
	GetSchedule(string) (*schedule.Weekly, error)
}

type FSStorage struct {
	path string
}

func (fss FSStorage) GetSchedule(slug string) (*schedule.Weekly, error) {
	f, err := os.Open(fss.path + slug + postfix)
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

func (fss FSStorage) CreateGroupSchedule(g *schedule.Group) error {
	evenStatus := ""
	if g.Schedule.IsEven {
		evenStatus = "even"
	} else {
		evenStatus = "odd"
	}
	filename := fss.path + fmt.Sprint(g.Course) + "-" + g.Name + "-" + evenStatus + postfix
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	e := json.NewEncoder(f)
	e.SetIndent("", "   ")
	if err := e.Encode(g); err != nil {
		return err
	}
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
