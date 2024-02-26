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

// Store defines methods for interacting with the data store.
type Store interface {
	GetInfo(request model.GetInfoRequest) (model.DBResponse, error)
}

// Mongo implements the Store interface for MongoDB.
type Mongo struct {
	DB     *mongo.Database
	Client *mongo.Client
}

// timeoutTimeBox sets the timeout duration for MongoDB operations.
const timeoutTimeBox = time.Duration(10)

// NewStore initializes a new MongoDB connection and returns a Store instance.
func NewStore(c config.Mongo) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutTimeBox*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(c.URI)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Printf("error occurs while connecting mongo db: %v", err)
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("error occurs while connecting mongo db: %v", err)
		return nil, err
	}

	db := client.Database(c.Database)

	return &Mongo{DB: db, Client: client}, nil
}

// GetInfo retrieves information from the MongoDB store based on the provided request.
func (m *Mongo) GetInfo(request model.GetInfoRequest) (model.DBResponse, error) {
	var results []bson.M

	var dbResponse model.DBResponse

	ctx, cancel := context.WithTimeout(context.Background(), timeoutTimeBox*time.Second)
	defer cancel()

	collection := m.DB.Collection("records")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, request.StartDate)
	endDate, _ := time.Parse(layout, request.EndDate)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "createdAt", Value: bson.D{
				{Key: "$gte", Value: primitive.NewDateTimeFromTime(startDate)},
				{Key: "$lte", Value: primitive.NewDateTimeFromTime(endDate)},
			}},
		}}},
		{{Key: "$unwind", Value: "$counts"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$key"},
			{Key: "totalCount", Value: bson.D{{Key: "$sum", Value: "$counts"}}},
			{Key: "counts", Value: bson.D{{Key: "$push", Value: "$counts"}}},
			{Key: "createdAt", Value: bson.D{{Key: "$first", Value: "$createdAt"}}},
		}}},
		{{Key: "$match", Value: bson.D{
			{Key: "totalCount", Value: bson.D{
				{Key: "$gte", Value: request.MinCount},
				{Key: "$lte", Value: request.MaxCount},
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

		if key, ok := result["_id"].(string); ok {
			record.Key = key
		}

		if createdAt, ok := result["createdAt"].(primitive.DateTime); ok {
			record.CreatedAt = createdAt.Time().UTC()
		}

		if totalCount, ok := result["totalCount"].(int64); ok { // veya int32, tipinize bağlı olarak
			record.TotalCount = int(totalCount)
		}

		dbResponse.Records = append(dbResponse.Records, record)
	}

	return dbResponse, nil
}
