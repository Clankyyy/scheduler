package storage

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"strconv"

	"github.com/Clankyyy/scheduler/internal/schedule"
)

const postfix = ".json"

type Storager interface {
	CreateGroupSchedule(g *schedule.Group) error
	DeleteSchedule(slug string) error
	GetWeeklyBySlug(slug string, query schedule.ScheduleType) ([]schedule.Weekly, error)
	GetDailyBySlug(slug string, weekday schedule.Weekday, query schedule.ScheduleType) (schedule.Daily, error)
	GetGroups() ([]GroupInfo, error)
	UpdateWeeklyBySlug(slug string, s []schedule.Weekly) error
}

type GroupInfo struct {
	Course int    `json:"course" bson:"course, omitempty"`
	Name   string `json:"name" bson:"name, omitemptys"`
}

type FSStorage struct {
	path string
}

func (fss FSStorage) GetDailyBySlug(slug string, day schedule.Weekday, dailyType schedule.ScheduleType) (schedule.Daily, error) {
	f, err := os.Open(fss.path + slug + "-" + dailyType.String() + postfix)
	if err != nil {
		log.Println("Error opening file", err.Error())
		return schedule.Daily{}, err
	}
	defer f.Close()

	var w schedule.Weekly
	if err := json.NewDecoder(f).Decode(&w); err != nil {
		log.Println("Error decoding type", err.Error())
		return schedule.Daily{}, err
	}
	d := w.Daily(day)

	return *d, nil
}

// func (fss FSStorage) GetWeeklyBySlug(slug string, param schedule.WeeklyQuery) ([]schedule.Weekly, error) {
// 	if param == schedule.WeeklyFull {
// 		w, err := fss.getFullWeekly(slug)
// 		return w, err
// 	} else if param == schedule.WeeklyEven || param == schedule.WeeklyOdd {
// 		f, err := os.Open(fss.path + slug + "-" + param.String() + postfix)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer f.Close()
// 		var w schedule.Weekly
// 		if err := json.NewDecoder(f).Decode(&w); err != nil {
// 			return nil, err
// 		}
// 		return []schedule.Weekly{w}, nil
// 	}

// 	return nil, fmt.Errorf("got some unexpected shit")
// }

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
	slug := strconv.Itoa(g.Course) + "-" + g.Name
	for i, v := range g.Schedule {
		evenStr := v.EvenString()
		path := fss.path + slug + "-" + evenStr + postfix
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

	return fss.addGroup(slug)
}

func (fss FSStorage) GetGroups() (GroupsInfo, error) {
	result := NewGroupInfo(20)
	f, err := os.Open(fss.path + "list" + postfix)
	if err != nil {
		return *result, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&result); err != nil {
		return *result, err
	}

	return *result, nil
}

func (fss FSStorage) DeleteSchedule(slug string) error {
	err := os.Remove(fss.path + slug + postfix)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}

// func (fss FSStorage) getFullWeekly(slug string) ([]schedule.Weekly, error) {
// 	statusList := []string{schedule.WeeklyEven.String(), schedule.WeeklyOdd.String()}
// 	w := make([]schedule.Weekly, 2)
// 	for i, v := range statusList {
// 		f, err := os.Open(fss.path + slug + "-" + v + postfix)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer f.Close()

// 		if err := json.NewDecoder(f).Decode(&w[i]); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return w, nil
// }

func (fss FSStorage) addGroup(slug string) error {
	f, err := os.OpenFile(fss.path+"list.json", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	groups := NewGroupInfo(20)
	if err := json.NewDecoder(f).Decode(&groups); err != nil {
		return err
	}
	if _, ok := groups.Names[slug]; ok {
		return errors.New("group already exists")
	}
	groups.Names[slug] = struct{}{}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "   ")
	return enc.Encode(&groups)
}

func NewFSStorage(path string) *FSStorage {
	return &FSStorage{
		path: path,
	}
}

type GroupsInfo struct {
	Names map[string]struct{}
}

func NewGroupInfo(len int) *GroupsInfo {
	g := &GroupsInfo{
		Names: make(map[string]struct{}, len),
	}
	return g
}
