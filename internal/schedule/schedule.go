package schedule

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
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
	fmt.Print(d.Weekday)
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
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Faculty  string    `json:"faculty"`
	Course   int       `json:"course"`
	Schedule Weekly    `json:"subjects"`
}

func Test() {
	s1 := NewSubject("14:00", "Информатика", "Федин", "416", Lecture)
	s2 := NewSubject("15:30", "Русский", "Хз", "116", Practice)

	day1 := []subject{*s1, *s2}
	day2 := []subject{*s2, *s1}
	d1 := NewDaily(day1, Monday)
	d2 := NewDaily(day2, Thursday)
	d1.Show()
	_ = d2
	w := Weekly{}

	f, err := os.Open("test.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&w); err != nil {
		fmt.Println(err)
	}
	fmt.Println(w)
	// file, err := os.Create("test.json")
	// if err != nil {
	// 	panic(err)
	// }

	// e := json.NewEncoder(file)
	// e.SetIndent("", "    ")
	// if err := e.Encode(&w); err != nil {
	// 	fmt.Print(err)
	// }
}

func NewGroup(name, faculty string, course int, s Weekly) *Group {
	return &Group{
		UUID:     uuid.New(),
		Name:     name,
		Faculty:  faculty,
		Course:   course,
		Schedule: s,
	}
}
