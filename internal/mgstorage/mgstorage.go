package mgstorage

import (
	"context"
	"errors"
	"log"
	"strconv"
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

func (mgs *MGStorage) GetDailyBySlug(slug string, weekday schedule.Weekday, query schedule.ScheduleType) (schedule.Daily, error) {
	course, name, err := ParseSlug(slug)
	if err != nil {
		return schedule.Daily{}, nil
	}

	col := mgs.client.Database("scheduler").Collection("groups")
	matchCN := bson.D{{"$match", bson.D{{"name", name}, {"course", course}}}}
	unwindSchedule := bson.D{{"$unwind", "$schedule"}}
	matchEven := bson.D{{"$match", bson.D{{"schedule.is_even", query.Boolean()}}}}
	unwindWeekly := bson.D{{"$unwind", "$schedule.weekly_schedule"}}
	matchDay := bson.D{{"$match", bson.D{{"schedule.weekly_schedule.weekday", weekday}}}}
	replaceRoot := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$schedule.weekly_schedule"}}}}

	cursor, err := col.Aggregate(context.Background(), mongo.Pipeline{matchCN, unwindSchedule, matchEven, unwindWeekly, matchDay, replaceRoot})
	if err != nil {
		log.Println(err)
		return schedule.Daily{}, err
	}

	var d []schedule.Daily
	if err = cursor.All(context.Background(), &d); err != nil {
		return schedule.Daily{}, err
	}
	if len(d) < 1 {
		return schedule.Daily{}, errors.New("no objects found")
	}

	return d[0], nil
}

func (mgs *MGStorage) GetWeeklyBySlug(slug string, query schedule.ScheduleType) (schedule.Weekly, error) {
	course, name, err := ParseSlug(slug)
	if err != nil {
		return schedule.Weekly{}, nil
	}

	col := mgs.client.Database("scheduler").Collection("groups")
	matchCN := bson.D{{"$match", bson.D{{"name", name}, {"course", course}}}}
	unwindSchedule := bson.D{{"$unwind", "$schedule"}}
	matchEven := bson.D{{"$match", bson.D{{"schedule.is_even", query.Boolean()}}}}
	replaceRoot := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$schedule"}}}}

	cursor, err := col.Aggregate(context.Background(), mongo.Pipeline{matchCN, unwindSchedule, matchEven, replaceRoot})
	if err != nil {
		log.Println(err)
		return schedule.Weekly{}, err
	}

	var w []schedule.Weekly
	if err = cursor.All(context.Background(), &w); err != nil {
		return schedule.Weekly{}, err
	}
	if len(w) < 1 {
		return schedule.Weekly{}, errors.New("no objects found")
	}

	return w[0], nil
}

func (mgs *MGStorage) UpdateWeeklyBySlug(slug string, s []schedule.Weekly) error {
	course, name, err := ParseSlug(slug)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "name", Value: name}, {Key: "course", Value: course}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "schedule", Value: s}}}}
	col := mgs.client.Database("scheduler").Collection("groups")

	_, err = col.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (mgs *MGStorage) GetFullWeeklyBySlug(slug string) ([]schedule.Weekly, error) {
	course, name, err := ParseSlug(slug)
	if err != nil {
		return nil, err
	}
	//opts := options.Find().SetProjection(bson.D{{Key: "schedule", Value: 1}})
	filter := bson.D{{Key: "name", Value: name}, {Key: "course", Value: course}}
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
	course, name, err := ParseSlug(slug)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "course", Value: course}, {Key: "name", Value: name}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))
	defer cancel()

	col := mgs.client.Database("scheduler").Collection("groups")
	res, err := col.DeleteOne(ctx, filter)
	log.Println("Filter is: ", filter, "deleted count:", res.DeletedCount)
	return err
}

func (mgs *MGStorage) CreateGroupSchedule(g *schedule.Group) error {
	collection := mgs.client.Database("scheduler").Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(500*time.Millisecond))
	defer cancel()
	_, err := collection.InsertOne(ctx, *g)
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

func ParseSlug(slug string) (course int, name string, err error) {
	splited := strings.Split(slug, "-")
	if len(splited) != 2 {
		return 0, "", errors.New("bad slug format")
	}
	courseStr, name := splited[0], splited[1]

	course, err = strconv.Atoi(courseStr)
	if err != nil {
		return 0, "", err
	}
	return course, name, nil
}

// db.groups.find(
// 	{ name: "4307",
// 		schedule: {
// 			$elemMatch: {
// 				is_even: true
// 			}
// 		},

// 	}
// )
// db.groups.findOne({"schedule.weekly_schedule": {$elemMatch: {weekday: 1}}}, {"schedule.weekly_schedule.$": 1})

/* db.groups.aggregate([{
	 $match: {course: 2, name: "4307"}
 },
 {
	 "$unwind": "$schedule"
 },
 {
	 $match: {"schedule.is_even": true}
 },
{
	 $unwind: "$schedule.weekly_schedule"
 },
 {
	 $match: {"schedule.weekly_schedule.weekday": 1}
 },
 {
	 $limit: 1
 },
 {
	 $replaceRoot: {newRoot: "$schedule.weekly_schedule"}
 },
 ])*/

/* db.groups.aggregate([{
	 $match: {course: 2, name: "4307"}
 },
 {
	 "$unwind": "$schedule"
 },
 {
	 $match: {"schedule.is_even": true}
 },
 {
	 $limit: 1
 },
 {
	 $replaceRoot: {newRoot: "$schedule"}
 },
 ])*/
