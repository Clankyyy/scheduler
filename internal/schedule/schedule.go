package schedule

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

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
	Faculty  int     `json:"faculty"`
	Course   int     `json:"course"`
	Subjects subject `json:"subjects"`
}

func Test() {
	s := NewSubject("14:00", "Информатика", "Федин", "416", Lecture)
	g := NewGroup("4305", 4, 2, *s)

	data, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}
	var g2 Group
	err = json.Unmarshal([]byte(data), &g2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", g)
	fmt.Printf("%+v\n", g2)
}
func NewGroup(name string, faculty int, course int, s subject) *Group {
	return &Group{
		ID:       rand.Intn(100000),
		Name:     name,
		Faculty:  faculty,
		Course:   course,
		Subjects: s,
	}
}
