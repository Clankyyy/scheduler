package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Clankyyy/scheduler/internal/schedule"
)

const postfix = ".json"

type Storager interface {
	CreateGroupSchedule(g *schedule.Group) error
	DeleteSchedule(string) error
	GetSchedule(string) ([]schedule.Weekly, error)
}

type FSStorage struct {
	path string
}

func (fss FSStorage) GetSchedule(slug string) ([]schedule.Weekly, error) {
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

func (fss FSStorage) CreateGroupSchedule(g *schedule.Group) error {

	filename := fss.path + fmt.Sprint(g.Course) + "-" + g.Name + "-" + postfix
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
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

//home\clanky\projects\scheduler\data\spbgti
