package mgstorage

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Clankyyy/scheduler/internal/schedule"
	"github.com/Clankyyy/scheduler/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// type Storager interface {
// 	GetWeeklyBySlug(string, schedule.WeeklyQuery) ([]schedule.Weekly, error)
// 	GetDailyBySlug(string, schedule.Weekday, schedule.DailyQuery) (schedule.Daily, error)
// 	GetGroups() (GroupsInfo, error)
// 	UpdateWeeklyBySlug([]schedule.Weekly, string) error
// }

type MGStorage struct {
	client *mongo.Client
}

func (mgs *MGStorage) GetDailyBySlug(slug string, weekday schedule.Weekday, query schedule.DailyQuery) (schedule.Daily, error) {
	return schedule.Daily{}, nil
}
func (mgs *MGStorage) UpdateWeeklyBySlug(slug string, s []schedule.Weekly) error {
	return nil
}

func (mgs *MGStorage) GetWeeklyBySlug(slug string, query schedule.WeeklyQuery) ([]schedule.Weekly, error) {
	_, _, err := ParseSlug(slug)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"name", "4307"}, {"course", 2}, {"schedule", bson.M{"$elemMatch": bson.M{"is_even": true}}}}
	// filter1 := fmt.Sprintf(`{"schedule": {"$elemMatch": {"is_even": true}}}`, name)
	// var filter bson.D
	// json.Unmarshal([]byte(filter1), &filter)

	col := mgs.client.Database("scheduler").Collection("groups")

	var g schedule.Group

	err = col.FindOne(context.Background(), filter).Decode(&g)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return g.Schedule, nil
}

func (mgs *MGStorage) GetGroups() ([]storage.GroupInfo, error) {
	col := mgs.client.Database("scheduler").Collection("groups")

	opts := options.Find().SetProjection(bson.D{{Key: "course", Value: 1}, {Key: "name", Value: 1}})
	filter := bson.D{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
	defer cancel()
	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var results []storage.GroupInfo
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (mgs *MGStorage) DeleteSchedule(slug string) error {
	col := mgs.client.Database("scheduler").Collection("groups")

	course, name, err := ParseSlug(slug)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "course", Value: course}, {Key: "name", Value: name}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))
	defer cancel()
	res, err := col.DeleteOne(ctx, filter)
	log.Println("Filter is: ", filter, "deleted count:", res.DeletedCount)
	return err
}

func (mgs *MGStorage) CreateGroupSchedule(g *schedule.Group) error {
	collection := mgs.client.Database("scheduler").Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))
	defer cancel()
	_, err := collection.InsertOne(ctx, *g)
	//log.Println("Inserted id:", res.InsertedID)
	return err
}

func NewMGStorage(uri string) *MGStorage {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*1))
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Cant connect to database with error: %s", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	return &MGStorage{
		client: client,
	}
}

func ParseSlug(slug string) (course, name string, err error) {
	splited := strings.Split(slug, "-")
	if len(splited) != 2 {
		return "", "", errors.New("bad slug format")
	}
	course, name = splited[0], splited[1]

	return course, name, nil
}

func MongoPipeline(str string) mongo.Pipeline {
	var pipeline = []bson.D{}
	str = strings.TrimSpace(str)
	if strings.Index(str, "[") != 0 {
		var doc bson.D
		bson.UnmarshalExtJSON([]byte(str), false, &doc)
		pipeline = append(pipeline, doc)
	} else {
		bson.UnmarshalExtJSON([]byte(str), false, &pipeline)
	}
	return pipeline
}

// db.groups.find(
// 	{ name: 4307,
// 		schedule: {
// 			$elemMatch: {
// 				is_even: true
// 			}
// 		}
// 	}
// )
