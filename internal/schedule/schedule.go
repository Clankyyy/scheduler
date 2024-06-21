package schedule

import "math/rand"

type Group struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Faculty int    `json:"faculty"`
	Course  int    `json:"course"`
	//schedule &
}

func NewGroup(name string, faculty int, course int) *Group {
	return &Group{
		ID:      rand.Intn(100000),
		Name:    name,
		Faculty: faculty,
		Course:  course,
	}
}
