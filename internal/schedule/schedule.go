package schedule

import (
	"fmt"
	"math/rand"
)

type Weekly struct {
	Schedule []Daily
	IsEven   bool
}

type Daily struct {
	Weekday  Weekday
	Schedule []subject
}

func NewDaily(s []subject, w Weekday) *Daily {
	return &Daily{
		Schedule: s,
		Weekday:  w,
	}
}

func (d Daily) Show() {
	for _, v := range d.Schedule {
		fmt.Println(v)
	}
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

func (w Weekday) String() string {
	return [...]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}[w-1]
}

func (w Weekday) EnumIndex() int {
	return int(w)
}

type subjectKind int

const (
	Lecture subjectKind = iota
	Practice
)

type subject struct {
	Start     string      `json:"start"`
	Name      string      `json:"name"`
	Teacher   string      `json:"teacher"`
	Classroom string      `json:"classroom"`
	Kind      subjectKind `json:"kind"`
}

func (s subjectKind) String() string {
	return [...]string{"Lecture", "Practice"}[s]
}

func (s subjectKind) EnumIndex() int {
	return int(s)
}

func NewSubject(startTime, name, teacher, classroom string, kind subjectKind) *subject {
	return &subject{
		Start:     startTime,
		Name:      name,
		Teacher:   teacher,
		Classroom: classroom,
		Kind:      kind,
	}
}

type Group struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Faculty  string  `json:"faculty"`
	Course   int     `json:"course"`
	Subjects subject `json:"subjects"`
}

func Test() {
	s1 := NewSubject("14:00", "Информатика", "Федин", "416", Lecture)
	s2 := NewSubject("15:30", "Русский", "Хз", "116", Practice)

	s2.Classroom = "0000"
	day := []subject{*s1, *s2}

	d := NewDaily(day, Monday)
	d.Show()
}

func NewGroup(name, faculty string, course int, s subject) *Group {
	return &Group{
		ID:       rand.Intn(100000),
		Name:     name,
		Faculty:  faculty,
		Course:   course,
		Subjects: s,
	}
}
