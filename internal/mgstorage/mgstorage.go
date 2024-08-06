package mgstorage

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// type Storager interface {
// 	CreateGroupSchedule(g *schedule.Group) error
// 	DeleteSchedule(string) error
// 	GetWeeklyBySlug(string, schedule.WeeklyQuery) ([]schedule.Weekly, error)
// 	GetDailyBySlug(string, schedule.Weekday, schedule.DailyQuery) (schedule.Daily, error)
// 	GetGroups() (GroupsInfo, error)
// 	UpdateWeeklyBySlug([]schedule.Weekly, string) error
// }

type MGStorage struct {
	client *mongo.Client
}

func NewMGStorage(uri string) *MGStorage {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Cant connect to database with error: %s", err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	return &MGStorage{
		client: client,
	}
}
