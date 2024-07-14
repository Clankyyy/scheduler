package schedule

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Weekly struct {
	Schedule []Daily `json:"weekly_schedule"`
	IsEven   bool    `json:"is_even"`
}

func (w Weekly) EvenString() string {
	if w.IsEven {
		return "even"
	}
	return "odd"
}

type Daily struct {
	Schedule []subject `json:"daily_schedule"`
	Weekday  Weekday   `json:"weekday"`
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

type ScheduleRequestParam int

const (
	Even ScheduleRequestParam = iota + 1
	Odd
	Full
)

func (sr ScheduleRequestParam) String() string {
	return [...]string{"even", "odd", "full"}[sr-1]
}

func BuildWeeklyType(param string) (ScheduleRequestParam, error) {
	if param == "even" {
		return Even, nil
	} else if param == "odd" {
		return Odd, nil
	} else if param == "full" {
		return Full, nil
	}
	return Even, fmt.Errorf("bad parameter")
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
	Schedule []Weekly  `json:"subjects"`
}

func Test() {
	s1 := NewSubject("14:00", "Информатика", "Федин", "416", Lecture)
	s2 := NewSubject("15:30", "Русский", "Хз", "116", Practice)

	day1 := []subject{*s1, *s2}
	d1 := NewDaily(day1, Monday)
	d1.Show()
	w := []Weekly{}
	w = append(w, Weekly{})
	w = append(w, Weekly{
		Schedule: []Daily{*d1, *d1},
		IsEven:   false,
	})

	f, err := os.Open("data/spbgti/2-4305-even.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&w[0]); err != nil {
		fmt.Println(err)
	}
	g := NewGroup("4305", "4", 2, w)
	data, err := json.MarshalIndent(g, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(data))
	file, err := os.Create("test.json")
	if err != nil {
		panic(err)
	}

	e := json.NewEncoder(file)
	e.SetIndent("", "    ")
	if err := e.Encode(&w); err != nil {
		fmt.Print(err)
	}
}

func NewGroup(name, faculty string, course int, s []Weekly) *Group {
	return &Group{
		UUID:     uuid.New(),
		Name:     name,
		Faculty:  faculty,
		Course:   course,
		Schedule: s,
	}
}
