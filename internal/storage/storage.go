package storage

import "github.com/Clankyyy/scheduler/internal/schedule"

type Storager interface {
	CreateGroup(*schedule.Group) error
	DeleteGroup(string) error
	GetSchedule(string) error
}

type FSStorage struct {
	path string
}

func NewFSStorage(path string) *FSStorage {
	return &FSStorage{
		path: path,
	}
}

//home\clanky\projects\scheduler\data\spbgti
