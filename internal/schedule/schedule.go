package schedule

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Weekly struct {
	Schedule []Daily `json:"weekly_schedule" bson:"weekly_schedule"`
	IsEven   bool    `json:"is_even" bson:"is_even"`
}

func (w Weekly) Daily(day Weekday) *Daily {
	for i := range w.Schedule {
		if w.Schedule[i].Weekday == day {
			return &w.Schedule[i]
		}
	}
	return &Daily{}
}

func (w Weekly) EvenString() string {
	if w.IsEven {
		return "even"
	}
	return "odd"
}

type Daily struct {
	Schedule []Subject `json:"daily_schedule" bson:"daily_schedule"`
	Weekday  Weekday   `json:"weekday" bson:"weekday"`
}

func NewDaily(s []Subject, w Weekday) *Daily {
	return &Daily{
		Schedule: s,
		Weekday:  w,
	}
}

func (d Daily) Show() {
	fmt.Print(d.Weekday)
	for _, v := range d.Schedule {
		fmt.Println(v)
	}
}

type ScheduleType int

const (
	ScheduleEven ScheduleType = iota + 1
	ScheduleOdd
)

func (sq ScheduleType) String() string {
	return [...]string{"even", "odd"}[sq-1]
}

func (sq ScheduleType) Boolean() bool {
	return [...]bool{true, false}[sq-1]
}

func BuildScheduleType(param string) (ScheduleType, error) {
	if param == "even" {
		return ScheduleEven, nil
	} else if param == "odd" {
		return ScheduleOdd, nil
	}
	return ScheduleEven, fmt.Errorf("bad value")
}

type Weekday int

const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (w Weekday) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

func (w Weekday) String() string {
	return [...]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}[w-1]
}

func (w Weekday) EnumIndex() int {
	return int(w)
}

func BuildWeekday(day string) (Weekday, error) {
	switch day {
	case "monday":
		return Monday, nil
	case "tuesday":
		return Tuesday, nil
	case "wednesday":
		return Wednesday, nil
	case "thursday":
		return Thursday, nil
	case "friday":
		return Friday, nil
	case "saturday":
		return Saturday, nil
	case "sunday":
		return Sunday, nil
	default:
		return Monday, fmt.Errorf("bad string")
	}
}

type SubjectKind int

const (
	Lecture SubjectKind = iota
	Practice
)

func (s SubjectKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s SubjectKind) String() string {
	return [...]string{"lecture", "practice"}[s]
}

func (s SubjectKind) EnumIndex() int {
	return int(s)
}

type Subject struct {
	Start     string      `json:"start"`
	Name      string      `json:"name"`
	Teacher   string      `json:"teacher"`
	Classroom string      `json:"classroom"`
	Kind      SubjectKind `json:"kind"`
}

func NewSubject(startTime, name, teacher, classroom string, kind SubjectKind) *Subject {
	return &Subject{
		Start:     startTime,
		Name:      name,
		Teacher:   teacher,
		Classroom: classroom,
		Kind:      kind,
	}
}

type Group struct {
	Name     string   `json:"name"`
	Faculty  string   `json:"faculty"`
	Course   int      `json:"course"`
	Schedule []Weekly `json:"subjects"`
}

func (g Group) Slug() string {
	return strconv.Itoa(g.Course) + "-" + g.Name
}

func NewGroup(name, faculty string, course int, s []Weekly) *Group {
	return &Group{
		Name:     name,
		Faculty:  faculty,
		Course:   course,
		Schedule: s,
	}
}
