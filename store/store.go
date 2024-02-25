package store

import (
	"context"
	"log"
	"main/config"
	"main/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Store interface {
	GetInfo(request model.GetInfoRequest) (model.DBResponse, error)
}

type Mongo struct {
	DB     *mongo.Database
	Client *mongo.Client
}

var (
	timeoutTimeBox = time.Duration(10)
)

func NewStore(c config.Mongo) *Mongo {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutTimeBox*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(c.URI)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	db := client.Database(c.Database)

	return &Mongo{DB: db, Client: client}
}

func (m *Mongo) GetInfo(request model.GetInfoRequest) (model.DBResponse, error) {
	var results []bson.M

	var dbResponse model.DBResponse

	ctx, cancel := context.WithTimeout(context.Background(), timeoutTimeBox*time.Second)
	defer cancel()

	collection := m.DB.Collection("records")

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"createdAt", bson.D{
				{"$gte", request.StartDate},
				{"$lte", request.EndDate},
			}},
		}}},
		{{"$unwind", "$counts"}},
		{{"$group", bson.D{
			{"_id", "$key"},
			{"totalCount", bson.D{{"$sum", "$counts"}}},
			{"createdAt", bson.D{{"$first", "$createdAt"}}},
		}}},
		{{"$match", bson.D{
			{"totalCount", bson.D{
				{"$gte", request.MinCount},
				{"$lte", request.MaxCount},
			}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return dbResponse, err
	}

	if err := cursor.All(ctx, &results); err != nil {
		return dbResponse, err
	}

	for _, result := range results {
		var record model.Record

		if key, ok := result["key"].(string); ok {
			record.Key = key
		}

		if createdAt, ok := result["createdAt"].(primitive.DateTime); ok {
			record.CreatedAt = createdAt.Time()
		}

		if totalCount, ok := result["totalCount"].(int32); ok {
			record.TotalCount = int(totalCount)
		}

		dbResponse.Records = append(dbResponse.Records, record)
	}

	return dbResponse, nil
}
